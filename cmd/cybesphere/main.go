// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/config"
	"cybesphere-backend/internal/middleware"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/routes"
	"cybesphere-backend/pkg/auth"
	"cybesphere-backend/pkg/database"
	"cybesphere-backend/pkg/logger"
)

func main() {
	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Inicializar logger
	if err := logger.Init(&cfg.Logging); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Starting CybESphere Backend")
	logger.Infof("Environment: %s", cfg.Monitoring.Environment)
	logger.Infof("Version: %s", cfg.Server.Version)

	// 3. Conectar a la base de datos
	if err := database.Connect(&cfg.Database); err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	logger.Info("Database connected successfully")

	// 4. Ejecutar auto-migraciones
	if err := runDatabaseMigrations(); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// 5. Seed data en desarrollo
	if cfg.Monitoring.Environment == "development" {
		if err := seedDevelopmentData(); err != nil {
			logger.Warnf("Failed to seed development data: %v", err)
		}
	}

	// 6. Inicializar JWT Manager
	jwtManager, err := auth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
		cfg.JWT.Issuer,
	)
	if err != nil {
		logger.Fatalf("Failed to initialize JWT manager: %v", err)
	}

	// 7. Crear AuthMiddleware desde JWT manager
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	logger.Info("JWT manager and auth middleware initialized successfully")

	// 8. Configurar Gin
	gin.SetMode(cfg.Server.Mode)
	r := setupRouter(cfg)

	// 9. Aplicar middleware global
	applyGlobalMiddleware(r, cfg)

	// 10. Inicializar aplicación con todas las dependencias
	app := routes.InitializeApplication(cfg, jwtManager)
	logger.Info("Application dependencies initialized")

	// 11. Configurar todas las rutas
	routes.SetupRoutes(r, cfg, authMiddleware, app)
	logger.Info(" Routes configured successfully")

	// 12. Iniciar servidor con graceful shutdown
	startServerWithGracefulShutdown(r, cfg)
}

// setupRouter configura el router de Gin con configuración básica
func setupRouter(cfg *config.Config) *gin.Engine {
	// Crear router con configuración según environment
	var r *gin.Engine

	if cfg.Monitoring.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	} else {
		r = gin.Default()
	}

	// Trust proxies para obtener IP real del cliente
	if err := r.SetTrustedProxies(cfg.Security.TrustedProxies); err != nil {
		logger.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Configurar tamaño máximo de multipart
	r.MaxMultipartMemory = 8 << 20 // 8 MB

	return r
}

// applyGlobalMiddleware aplica todos los middlewares globales
func applyGlobalMiddleware(r *gin.Engine, cfg *config.Config) {
	// Middleware básico (orden importante)
	r.Use(middleware.RequestID())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RecoveryWithLogger())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.ErrorHandler())

	// CORS para desarrollo
	if cfg.Monitoring.Environment == "development" || cfg.Security.CORSEnabled {
		r.Use(middleware.CORSMiddleware(cfg))
		logger.Info("CORS enabled")
	}

	// Rate limiting
	if cfg.RateLimit.Enabled {
		r.Use(middleware.RateLimitMiddleware(cfg))
		logger.Infof("Rate limiting enabled: %d requests/minute", cfg.RateLimit.RequestsPerMinute)
	}

	// Timeout global
	r.Use(middleware.Timeout(5 * time.Minute))

	// Validación de Content-Type para APIs
	r.Use(middleware.ContentTypeValidation("application/json", "multipart/form-data"))

	// Versionado de API
	r.Use(middleware.APIVersioning())

	// Validación de errores
	r.Use(middleware.ValidationError())

	// Context básico para logging
	r.Use(middleware.BasicUserContext())
}

// runDatabaseMigrations ejecuta las migraciones de base de datos
func runDatabaseMigrations() error {
	logger.Info("Running database migrations...")

	if err := models.AutoMigrate(database.GetDB()); err != nil {
		return err
	}

	logger.Info("Creating database indexes...")
	if err := models.CreateIndexes(database.GetDB()); err != nil {
		logger.Warnf("Some indexes could not be created: %v", err)
		// No fallar si algunos índices no se pueden crear
	}

	logger.Info("Database setup completed successfully")
	return nil
}

// seedDevelopmentData crea datos de prueba para desarrollo
func seedDevelopmentData() error {
	logger.Info("Seeding development data...")

	if err := models.SeedData(database.GetDB()); err != nil {
		return err
	}

	logger.Info("Development data seeded successfully")
	return nil
}

// startServerWithGracefulShutdown inicia el servidor con graceful shutdown
func startServerWithGracefulShutdown(r *gin.Engine, cfg *config.Config) {
	address := cfg.Server.GetAddress()

	// Crear servidor HTTP
	srv := &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Goroutine para iniciar el servidor
	go func() {
		printStartupBanner(cfg, address)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Canal para señales del sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Intentar shutdown graceful
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	// Cerrar conexiones de base de datos
	if err := database.Close(); err != nil {
		logger.Errorf("Error cerrando la base de datos: %v", err)
	}

	logger.Info("Server exited successfully")
}

// printStartupBanner imprime información de inicio
func printStartupBanner(cfg *config.Config, address string) {
	logger.Info("")
	logger.Info("╔════════════════════════════════════════════════════════╗")
	logger.Info("║            CybESphere Backend Started                  ║")
	logger.Info("╚════════════════════════════════════════════════════════╝")
	logger.Info("")
	logger.Infof(" Server running on %s", address)
	logger.Info("")
	logger.Info(" Available endpoints:")
	logger.Info("├── Health & Status")
	logger.Info("│   ├── GET /health - Health check")
	logger.Info("│   ├── GET /api/v1/ping - Ping test")
	logger.Info("│   └── GET /docs - API Documentation")
	logger.Info("│")
	logger.Info("├── Authentication")
	logger.Info("│   ├── POST /api/v1/auth/register - Register new user")
	logger.Info("│   ├── POST /api/v1/auth/login - User login")
	logger.Info("│   ├── POST /api/v1/auth/refresh - Refresh tokens")
	logger.Info("│   ├── POST /api/v1/auth/logout - Logout current session")
	logger.Info("│   ├── POST /api/v1/auth/logout-all - Logout all sessions")
	logger.Info("│   └── GET  /api/v1/auth/me - Current user info")
	logger.Info("│")
	logger.Info("├── Public Endpoints")
	logger.Info("│   ├── GET /api/v1/public/events - Public events")
	logger.Info("│   ├── GET /api/v1/public/events/:id - Event details")
	logger.Info("│   ├── GET /api/v1/public/organizations - Public organizations")
	logger.Info("│   └── GET /api/v1/public/stats - Public statistics")
	logger.Info("│")
	logger.Info("├── Protected Endpoints")
	logger.Info("│   ├── GET /api/v1/user/capabilities - User capabilities")
	logger.Info("│   ├── GET /api/v1/events - List events")
	logger.Info("│   ├── POST /api/v1/events - Create event")
	logger.Info("│   ├── GET /api/v1/organizations - List organizations")
	logger.Info("│   └── POST /api/v1/organizations - Create organization")
	logger.Info("│")
	logger.Info("├── Admin Endpoints")
	logger.Info("│   ├── GET /api/v1/admin/dashboard - Admin dashboard")
	logger.Info("│   ├── GET /api/v1/admin/system/stats - System statistics")
	logger.Info("│   └── GET /api/v1/admin/audit-logs - Audit logs")
	logger.Info("│")
	logger.Info("└── Organizer Endpoints")
	logger.Info("    ├── GET /api/v1/organizer/dashboard - Organizer dashboard")
	logger.Info("    └── GET /api/v1/organizer/events - Organization events")
	logger.Info("")
	logger.Info("Configuration:")
	logger.Infof("  ├── Environment: %s", cfg.Monitoring.Environment)
	logger.Infof("  ├── Log level: %s", cfg.Logging.Level)
	logger.Infof("  ├── JWT expiration: %v", cfg.JWT.AccessTokenDuration)
	logger.Infof("  ├── Database: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	logger.Infof("  ├── Rate limiting: %v", cfg.RateLimit.Enabled)

	if cfg.Monitoring.Environment == "development" {
		logger.Infof("  ├── CORS origins: %v", cfg.Security.CORSAllowedOrigins)
		logger.Infof("  └── Docs enabled: %v", cfg.Server.EnableDocs)
	} else {
		logger.Info("  └── Production mode: Security headers enabled")
	}

	logger.Info("")
	logger.Info("Ready to serve requests!")
	logger.Info("")
}
