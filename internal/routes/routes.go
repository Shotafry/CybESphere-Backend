// internal/routes/routes.go
package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/config"
	"cybesphere-backend/internal/handlers"
	"cybesphere-backend/internal/helpers"
	"cybesphere-backend/internal/mappers"
	"cybesphere-backend/internal/middleware"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/permissions"
	"cybesphere-backend/internal/repositories"
	"cybesphere-backend/internal/services"
	"cybesphere-backend/pkg/auth"
	"cybesphere-backend/pkg/database"
)

// Application estructura que contiene todas las dependencias
type Application struct {
	Config       *config.Config
	Repositories *repositories.RepositoryManager
	Services     *ServiceContainer
	Handlers     *HandlerContainer
	Mapper       *mappers.UnifiedMapper
}

// ServiceContainer contiene todos los servicios
type ServiceContainer struct {
	Auth          services.AuthService
	Authorization services.AuthorizationService
	Events        services.EventService
	Organizations services.OrganizationService
	Users         services.UserService
}

// HandlerContainer contiene todos los handlers
type HandlerContainer struct {
	Auth          *handlers.AuthHandler
	Events        *handlers.EventHandler
	Organizations *handlers.OrganizationHandler
	Users         *handlers.UserHandler
	Capabilities  *handlers.UserCapabilitiesHandler
}

// InitializeApplication inicializa toda la aplicación con sus dependencias
func InitializeApplication(cfg *config.Config, jwtManager *auth.JWTManager) *Application {
	// 1. Crear repositories
	repoManager := repositories.NewRepositoryManager()

	// 2. Crear mapper unificado
	mapper := mappers.NewUnifiedMapper()

	// 3. Crear authorization service
	authorizationService := services.NewAuthorizationService()

	// 4. Crear auth service
	authService := services.NewAuthServiceImpl(
		repoManager.Users,
		repoManager.RefreshTokens,
		jwtManager,
		mapper,
	)

	// 5. Crear service manager
	serviceManager := services.NewServiceManager(
		repoManager,
		mapper,
		authorizationService,
	)

	// 6. Container de servicios
	serviceContainer := &ServiceContainer{
		Auth:          authService,
		Authorization: authorizationService,
		Events:        serviceManager.Events,
		Organizations: serviceManager.Organizations,
		Users:         serviceManager.Users,
	}

	// 7. Crear handlers
	handlerContainer := &HandlerContainer{
		Auth: handlers.NewAuthHandler(
			authService,
			serviceManager.Users,
			mapper,
			jwtManager,
		),
		Events: handlers.NewEventHandler(
			serviceManager.Events,
			mapper,
		),
		Organizations: handlers.NewOrganizationHandler(
			serviceManager.Organizations,
			mapper,
		),
		Users: handlers.NewUserHandler(
			serviceManager.Users,
			mapper,
		),
		Capabilities: handlers.NewUserCapabilitiesHandler(
			authorizationService,
			serviceManager.Users,
			mapper,
		),
	}

	return &Application{
		Config:       cfg,
		Repositories: repoManager,
		Services:     serviceContainer,
		Handlers:     handlerContainer,
		Mapper:       mapper,
	}
}

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(r *gin.Engine, cfg *config.Config, authMiddleware *middleware.AuthMiddleware, app *Application) {
	// Health check endpoint
	r.GET("/health", healthCheck)

	// Rutas de documentación
	if cfg.Server.EnableDocs && cfg.Monitoring.Environment != "production" {
		r.GET("/docs", documentationEndpoint(cfg))
	}

	// API v1
	v1 := r.Group("/api/" + cfg.Server.Version)

	// Configurar rutas públicas de auth
	setupAuthRoutes(v1, app.Handlers.Auth)

	// Configurar rutas públicas con auth opcional
	setupPublicRoutes(v1, cfg, authMiddleware, app)

	// Configurar rutas protegidas
	setupProtectedRoutes(v1, authMiddleware, app)

	// Configurar rutas de administración
	setupAdminRoutes(v1, authMiddleware, app)

	// Configurar rutas de organizador
	setupOrganizerRoutes(v1, authMiddleware)
}

// setupAuthRoutes configura las rutas de autenticación (públicas)
func setupAuthRoutes(v1 *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}
}

// setupPublicRoutes configura las rutas públicas con autenticación opcional
func setupPublicRoutes(v1 *gin.RouterGroup, cfg *config.Config, authMiddleware *middleware.AuthMiddleware, app *Application) {
	public := v1.Group("/public")
	public.Use(authMiddleware.OptionalAuthFlow())
	public.Use(middleware.QueryOptions())
	{
		// Ping endpoint con información opcional de usuario
		public.GET("/ping", pingEndpoint(cfg))

		// Eventos públicos
		public.GET("/events", app.Handlers.Events.GetAll)
		public.GET("/events/:id", app.Handlers.Events.GetPublicEvent)
		public.GET("/events/featured", app.Handlers.Events.GetFeaturedEvents)
		public.GET("/events/upcoming", app.Handlers.Events.GetUpcomingEvents)

		// Organizaciones públicas
		public.GET("/organizations", app.Handlers.Organizations.GetAll)
		public.GET("/organizations/:id", app.Handlers.Organizations.GetByID)
		public.GET("/organizations/active", app.Handlers.Organizations.GetActiveOrganizations)

		// Estadísticas públicas
		public.GET("/stats", publicStatsEndpoint)

		// Debug endpoint para desarrollo
		if cfg.Monitoring.Environment != "production" {
			public.GET("/debug/db", debugDBEndpoint)
		}
	}
}

// setupProtectedRoutes configura las rutas protegidas que requieren autenticación
func setupProtectedRoutes(v1 *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, app *Application) {
	// Grupo protegido con autenticación completa
	protected := v1.Group("")
	protected.Use(authMiddleware.AuthFlow())
	protected.Use(middleware.EnhancedUserContext())
	protected.Use(middleware.QueryOptions())
	{
		// Auth endpoints
		authGroup := protected.Group("/auth")
		{
			authGroup.GET("/me", app.Handlers.Auth.Me)
			authGroup.POST("/logout-all", app.Handlers.Auth.LogoutAll)
		}

		// User capabilities
		userGroup := protected.Group("/user")
		{
			userGroup.GET("/capabilities", app.Handlers.Capabilities.GetUserCapabilities)
			userGroup.POST("/check-access", app.Handlers.Capabilities.CheckResourceAccess)
			userGroup.GET("/available-actions", app.Handlers.Capabilities.GetAvailableActions)
			userGroup.GET("/sessions", app.Handlers.Capabilities.GetUserSessions)
			userGroup.DELETE("/sessions/:sessionId", app.Handlers.Capabilities.RevokeSession)
			userGroup.GET("/roles", app.Handlers.Capabilities.GetRoleInfo)
			userGroup.GET("/profile", app.Handlers.Users.GetUserProfile)
		}

		// Events - CRUD con BaseHandler
		eventsGroup := protected.Group("/events")
		{
			eventsGroup.GET("", app.Handlers.Events.GetAll)
			eventsGroup.GET("/:id", app.Handlers.Events.GetByID)

			// Crear evento requiere permisos de escritura
			eventsGroup.POST("",
				authMiddleware.RequirePermissionEnhanced(permissions.WriteEvent),
				app.Handlers.Events.Create)

			// Actualizar y eliminar requieren ownership
			eventsGroup.PUT("/:id",
				authMiddleware.GuardEvent(permissions.WriteEvent),
				app.Handlers.Events.Update)

			eventsGroup.DELETE("/:id",
				authMiddleware.GuardEvent(permissions.DeleteEvent),
				app.Handlers.Events.Delete)

			// Acciones específicas
			eventsGroup.POST("/:id/publish",
				authMiddleware.GuardEvent(permissions.PublishEvent),
				app.Handlers.Events.PublishEvent)

			eventsGroup.POST("/:id/cancel",
				authMiddleware.GuardEvent(permissions.WriteEvent),
				app.Handlers.Events.CancelEvent)

			// Favoritos (cualquier usuario autenticado)
			eventsGroup.POST("/:id/favorite", app.Handlers.Events.AddToFavorites)
			eventsGroup.DELETE("/:id/favorite", app.Handlers.Events.RemoveFromFavorites)

			// Eventos por organización
			eventsGroup.GET("/organization/:orgId", app.Handlers.Events.GetEventsByOrganization)
		}

		// Organizations - CRUD con BaseHandler
		orgsGroup := protected.Group("/organizations")
		{
			orgsGroup.GET("", app.Handlers.Organizations.GetAll)
			orgsGroup.GET("/:id", app.Handlers.Organizations.GetByID)

			// Crear organización requiere usuario verificado
			orgsGroup.POST("",
				authMiddleware.RequirePermissionEnhanced(permissions.WriteOrganization),
				app.Handlers.Organizations.Create)

			// Actualizar requiere ser miembro o admin
			orgsGroup.PUT("/:id",
				authMiddleware.GuardOrganization(permissions.WriteOrganization),
				app.Handlers.Organizations.Update)

			// Solo admin puede eliminar
			orgsGroup.DELETE("/:id",
				authMiddleware.ForAdminOnly(),
				app.Handlers.Organizations.Delete)

			// Verificación solo admin
			orgsGroup.POST("/:id/verify",
				authMiddleware.ForAdminOnly(),
				app.Handlers.Organizations.VerifyOrganization)

			// Miembros (solo miembros de la org o admin)
			orgsGroup.GET("/:id/members",
				authMiddleware.GuardOrganization(permissions.ReadOrganization),
				app.Handlers.Organizations.GetMembers)
		}

		// Users management
		usersGroup := protected.Group("/users")
		{
			// Lista de usuarios solo admin
			usersGroup.GET("",
				authMiddleware.ForAdminOnly(),
				app.Handlers.Users.GetAll)

			// Ver perfil (propio o si es admin)
			usersGroup.GET("/:id",
				authMiddleware.GuardUser(permissions.ReadProfile),
				app.Handlers.Users.GetByID)

			// Actualizar perfil (propio o si es admin)
			usersGroup.PUT("/:id",
				authMiddleware.GuardUser(permissions.WriteProfile),
				app.Handlers.Users.Update)

			// Acciones de admin
			usersGroup.PUT("/:id/role",
				authMiddleware.ForAdminOnly(),
				app.Handlers.Users.UpdateRole)

			usersGroup.POST("/:id/activate",
				authMiddleware.ForAdminOnly(),
				app.Handlers.Users.ActivateUser)

			usersGroup.POST("/:id/deactivate",
				authMiddleware.ForAdminOnly(),
				app.Handlers.Users.DeactivateUser)

			// Sesiones del usuario
			usersGroup.GET("/:id/sessions",
				authMiddleware.GuardUser(permissions.ReadProfile),
				app.Handlers.Users.GetUserSessions)
		}
	}
}

// setupAdminRoutes rutas de administración
func setupAdminRoutes(v1 *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware, app *Application) {
	admin := v1.Group("/admin")
	admin.Use(authMiddleware.ForAdminOnly())
	admin.Use(middleware.QueryOptions())
	{
		// Dashboard
		admin.GET("/dashboard", adminDashboard)

		// Estadísticas del sistema
		admin.GET("/system/stats", systemStatsEndpoint)

		// Logs de auditoría
		admin.GET("/audit-logs", auditLogsEndpoint)

		// Configuración del sistema
		admin.GET("/system/config", systemConfigEndpoint)

		// Gestión masiva de usuarios
		admin.GET("/users/export", app.Handlers.Users.GetAll)

		// Gestión masiva de organizaciones
		admin.POST("/organizations/bulk-verify", bulkVerifyOrganizations)

		// Gestión masiva de eventos
		admin.POST("/events/bulk-moderate", bulkModerateEvents)
	}
}

// setupOrganizerRoutes rutas de organizador
func setupOrganizerRoutes(v1 *gin.RouterGroup, authMiddleware *middleware.AuthMiddleware) {
	organizer := v1.Group("/organizer")
	organizer.Use(authMiddleware.RequireOrganizerOrAdmin())
	organizer.Use(middleware.EnhancedUserContext())
	organizer.Use(middleware.QueryOptions())
	{
		// Dashboard del organizador
		organizer.GET("/dashboard", organizerDashboard)

		// Eventos de la organización
		organizer.GET("/events", organizerEventsEndpoint)

		// Estadísticas de la organización
		organizer.GET("/stats", organizerStatsEndpoint)

		// Miembros de la organización
		organizer.GET("/members", organizerMembersEndpoint)
	}
}

// ============================================
// Endpoints individuales
// ============================================

func healthCheck(c *gin.Context) {
	dbStats, _ := database.GetStats()

	helpers.FormatSuccessResponse(c, gin.H{
		"status":    "ok",
		"message":   "CybESphere Backend is running",
		"version":   "0.1.0",
		"timestamp": time.Now().UTC(),
		"database": gin.H{
			"connected": database.IsConnected(),
			"stats":     dbStats,
		},
	}, "Health check successful")
}

func documentationEndpoint(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		helpers.FormatSuccessResponse(c, gin.H{
			"version": cfg.Server.Version,
			"endpoints": gin.H{
				"auth": gin.H{
					"POST /api/v1/auth/register":   "Registro de usuario",
					"POST /api/v1/auth/login":      "Inicio de sesión",
					"POST /api/v1/auth/refresh":    "Renovar tokens",
					"POST /api/v1/auth/logout":     "Cerrar sesión",
					"POST /api/v1/auth/logout-all": "Cerrar todas las sesiones",
					"GET  /api/v1/auth/me":         "Información del usuario actual",
				},
				"public": gin.H{
					"GET /api/v1/public/ping":                 "Ping test",
					"GET /health":                             "Health check",
					"GET /api/v1/public/events":               "Lista de eventos públicos",
					"GET /api/v1/public/events/:id":           "Detalle de evento público",
					"GET /api/v1/public/events/featured":      "Eventos destacados",
					"GET /api/v1/public/events/upcoming":      "Próximos eventos",
					"GET /api/v1/public/organizations":        "Lista de organizaciones públicas",
					"GET /api/v1/public/organizations/:id":    "Detalle de organización",
					"GET /api/v1/public/organizations/active": "Organizaciones activas",
					"GET /api/v1/public/stats":                "Estadísticas públicas",
				},
				"protected": gin.H{
					"GET /api/v1/user/capabilities":         "Capacidades del usuario",
					"GET /api/v1/user/profile":              "Perfil del usuario actual",
					"GET /api/v1/user/sessions":             "Sesiones activas",
					"GET /api/v1/user/roles":                "Información de roles",
					"GET /api/v1/events":                    "Lista de eventos",
					"POST /api/v1/events":                   "Crear evento",
					"PUT /api/v1/events/:id":                "Actualizar evento",
					"DELETE /api/v1/events/:id":             "Eliminar evento",
					"POST /api/v1/events/:id/publish":       "Publicar evento",
					"POST /api/v1/events/:id/cancel":        "Cancelar evento",
					"GET /api/v1/organizations":             "Lista de organizaciones",
					"POST /api/v1/organizations":            "Crear organización",
					"PUT /api/v1/organizations/:id":         "Actualizar organización",
					"GET /api/v1/organizations/:id/members": "Miembros de organización",
				},
				"admin": gin.H{
					"GET /api/v1/admin/dashboard":           "Dashboard de administrador",
					"GET /api/v1/admin/system/stats":        "Estadísticas del sistema",
					"GET /api/v1/admin/audit-logs":          "Logs de auditoría",
					"GET /api/v1/admin/system/config":       "Configuración del sistema",
					"POST /api/v1/organizations/:id/verify": "Verificar organización",
					"PUT /api/v1/users/:id/role":            "Cambiar rol de usuario",
				},
				"organizer": gin.H{
					"GET /api/v1/organizer/dashboard": "Dashboard de organizador",
					"GET /api/v1/organizer/events":    "Eventos de la organización",
					"GET /api/v1/organizer/stats":     "Estadísticas de la organización",
					"GET /api/v1/organizer/members":   "Miembros de la organización",
				},
			},
		}, "API Documentation")
	}
}

func pingEndpoint(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx := handlers.GetUserContext(c)

		response := gin.H{
			"message":       "pong",
			"environment":   cfg.Monitoring.Environment,
			"timestamp":     time.Now().UTC(),
			"authenticated": false,
		}

		if userCtx != nil {
			response["authenticated"] = true
			response["user"] = gin.H{
				"id":    userCtx.ID,
				"email": userCtx.Email,
				"role":  string(userCtx.Role),
			}
		}

		helpers.FormatSuccessResponse(c, response, "Ping successful")
	}
}

func debugDBEndpoint(c *gin.Context) {
	db := database.GetDB()
	if db == nil {
		helpers.FormatErrorResponse(c, http.StatusInternalServerError, "database_error", "Database connection is nil")
		return
	}

	var stats struct {
		Organizations int64 `json:"organizations"`
		Events        int64 `json:"events"`
		Users         int64 `json:"users"`
	}

	db.Model(&models.Organization{}).Count(&stats.Organizations)
	db.Model(&models.Event{}).Count(&stats.Events)
	db.Model(&models.User{}).Count(&stats.Users)

	helpers.FormatSuccessResponse(c, gin.H{"database_stats": stats}, "Database stats retrieved")
}

func publicStatsEndpoint(c *gin.Context) {
	db := database.GetDB()

	var stats struct {
		PublicEvents          int64 `json:"public_events"`
		ActiveOrganizations   int64 `json:"active_organizations"`
		VerifiedOrganizations int64 `json:"verified_organizations"`
		UpcomingEvents        int64 `json:"upcoming_events"`
	}

	db.Model(&models.Event{}).Where("status = ? AND is_public = ?", "published", true).Count(&stats.PublicEvents)
	db.Model(&models.Organization{}).Where("status = ?", "active").Count(&stats.ActiveOrganizations)
	db.Model(&models.Organization{}).Where("is_verified = true").Count(&stats.VerifiedOrganizations)
	db.Model(&models.Event{}).Where("status = ? AND start_date > ?", "published", time.Now()).Count(&stats.UpcomingEvents)

	helpers.FormatSuccessResponse(c, gin.H{
		"data":         stats,
		"generated_at": time.Now().UTC(),
	}, "Public statistics retrieved")
}

func adminDashboard(c *gin.Context) {
	db := database.GetDB()

	var stats struct {
		TotalUsers           int64 `json:"total_users"`
		ActiveUsers          int64 `json:"active_users"`
		VerifiedUsers        int64 `json:"verified_users"`
		TotalOrganizations   int64 `json:"total_organizations"`
		ActiveOrganizations  int64 `json:"active_organizations"`
		PendingOrganizations int64 `json:"pending_organizations"`
		TotalEvents          int64 `json:"total_events"`
		PublishedEvents      int64 `json:"published_events"`
		UpcomingEvents       int64 `json:"upcoming_events"`
	}

	// Users stats
	db.Model(&models.User{}).Count(&stats.TotalUsers)
	db.Model(&models.User{}).Where("is_active = true").Count(&stats.ActiveUsers)
	db.Model(&models.User{}).Where("is_verified = true").Count(&stats.VerifiedUsers)

	// Organizations stats
	db.Model(&models.Organization{}).Count(&stats.TotalOrganizations)
	db.Model(&models.Organization{}).Where("status = ?", "active").Count(&stats.ActiveOrganizations)
	db.Model(&models.Organization{}).Where("status = ?", "pending").Count(&stats.PendingOrganizations)

	// Events stats
	db.Model(&models.Event{}).Count(&stats.TotalEvents)
	db.Model(&models.Event{}).Where("status = ?", "published").Count(&stats.PublishedEvents)
	db.Model(&models.Event{}).Where("status = ? AND start_date > ?", "published", time.Now()).Count(&stats.UpcomingEvents)

	helpers.FormatSuccessResponse(c, gin.H{
		"data":         stats,
		"generated_at": time.Now().UTC(),
	}, "Admin dashboard data retrieved")
}

func organizerDashboard(c *gin.Context) {
	userCtx := handlers.GetUserContext(c)
	if userCtx == nil || userCtx.OrganizationID == nil {
		helpers.FormatErrorResponse(c, http.StatusForbidden, "no_organization",
			"Debes pertenecer a una organización para acceder a esta sección")
		return
	}

	db := database.GetDB()
	orgID := *userCtx.OrganizationID

	var organization models.Organization
	if err := db.First(&organization, "id = ?", orgID).Error; err != nil {
		helpers.FormatErrorResponse(c, http.StatusNotFound, "organization_not_found",
			"Organización no encontrada")
		return
	}

	var stats struct {
		TotalEvents     int64 `json:"total_events"`
		PublishedEvents int64 `json:"published_events"`
		DraftEvents     int64 `json:"draft_events"`
		UpcomingEvents  int64 `json:"upcoming_events"`
		TotalAttendees  int64 `json:"total_attendees"`
		MembersCount    int64 `json:"members_count"`
	}

	// Events stats
	db.Model(&models.Event{}).Where("organization_id = ?", orgID).Count(&stats.TotalEvents)
	db.Model(&models.Event{}).Where("organization_id = ? AND status = ?", orgID, "published").Count(&stats.PublishedEvents)
	db.Model(&models.Event{}).Where("organization_id = ? AND status = ?", orgID, "draft").Count(&stats.DraftEvents)
	db.Model(&models.Event{}).Where("organization_id = ? AND status = ? AND start_date > ?",
		orgID, "published", time.Now()).Count(&stats.UpcomingEvents)

	// Members count
	db.Model(&models.User{}).Where("organization_id = ?", orgID).Count(&stats.MembersCount)

	// Total attendees across all events
	db.Model(&models.Event{}).Where("organization_id = ?", orgID).
		Select("COALESCE(SUM(current_attendees), 0)").Scan(&stats.TotalAttendees)

	helpers.FormatSuccessResponse(c, gin.H{
		"organization": organization,
		"stats":        stats,
		"generated_at": time.Now().UTC(),
	}, "Organizer dashboard data retrieved")
}

func organizerEventsEndpoint(c *gin.Context) {
	userCtx := handlers.GetUserContext(c)
	if userCtx == nil || userCtx.OrganizationID == nil {
		helpers.FormatErrorResponse(c, http.StatusForbidden, "no_organization",
			"No perteneces a ninguna organización")
		return
	}

	// Usar el filtro de organización
	c.Set("organization_filter", *userCtx.OrganizationID)
	helpers.FormatSuccessResponse(c, gin.H{
		"message": "Use /api/v1/events with organization filter",
	}, "Filter set for organization events")
}

func organizerStatsEndpoint(c *gin.Context) {
	userCtx := handlers.GetUserContext(c)
	if userCtx == nil || userCtx.OrganizationID == nil {
		helpers.FormatErrorResponse(c, http.StatusForbidden, "no_organization",
			"No perteneces a ninguna organización")
		return
	}

	db := database.GetDB()
	orgID := *userCtx.OrganizationID

	// Estadísticas detalladas de la organización
	var stats struct {
		EventsByType      map[string]int `json:"events_by_type"`
		EventsByStatus    map[string]int `json:"events_by_status"`
		AttendeesPerMonth []struct {
			Month string `json:"month"`
			Count int    `json:"count"`
		} `json:"attendees_per_month"`
		TopEvents []struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			Attendees int    `json:"attendees"`
		} `json:"top_events"`
	}

	// Eventos por tipo
	var eventTypes []struct {
		Type  string
		Count int
	}
	db.Model(&models.Event{}).
		Select("event_type as type, COUNT(*) as count").
		Where("organization_id = ?", orgID).
		Group("event_type").
		Scan(&eventTypes)

	stats.EventsByType = make(map[string]int)
	for _, et := range eventTypes {
		stats.EventsByType[et.Type] = et.Count
	}

	// Eventos por estado
	var eventStatus []struct {
		Status string
		Count  int
	}
	db.Model(&models.Event{}).
		Select("status, COUNT(*) as count").
		Where("organization_id = ?", orgID).
		Group("status").
		Scan(&eventStatus)

	stats.EventsByStatus = make(map[string]int)
	for _, es := range eventStatus {
		stats.EventsByStatus[es.Status] = es.Count
	}

	helpers.FormatSuccessResponse(c, gin.H{
		"data":         stats,
		"generated_at": time.Now().UTC(),
	}, "Organization statistics retrieved")
}

func organizerMembersEndpoint(c *gin.Context) {
	userCtx := handlers.GetUserContext(c)
	if userCtx == nil || userCtx.OrganizationID == nil {
		helpers.FormatErrorResponse(c, http.StatusForbidden, "no_organization",
			"No perteneces a ninguna organización")
		return
	}

	helpers.FormatSuccessResponse(c, gin.H{
		"message":         "Use /api/v1/organizations/:id/members",
		"organization_id": *userCtx.OrganizationID,
	}, "Use the organizations endpoint for members")
}

func systemStatsEndpoint(c *gin.Context) {
	db := database.GetDB()

	var stats struct {
		DatabaseSize   string `json:"database_size"`
		TotalRecords   int64  `json:"total_records"`
		ActiveSessions int64  `json:"active_sessions"`
		SystemHealth   string `json:"system_health"`
	}

	// Contar registros totales
	var userCount, eventCount, orgCount int64
	db.Model(&models.User{}).Count(&userCount)
	db.Model(&models.Event{}).Count(&eventCount)
	db.Model(&models.Organization{}).Count(&orgCount)
	stats.TotalRecords = userCount + eventCount + orgCount

	// Contar sesiones activas
	db.Model(&models.RefreshToken{}).Where("is_revoked = false AND expires_at > ?", time.Now()).Count(&stats.ActiveSessions)

	// Estado del sistema
	stats.SystemHealth = "healthy"
	stats.DatabaseSize = "N/A"

	helpers.FormatSuccessResponse(c, gin.H{
		"data":         stats,
		"generated_at": time.Now().UTC(),
	}, "System statistics retrieved")
}

func auditLogsEndpoint(c *gin.Context) {
	db := database.GetDB()

	// Usar helpers para extraer parámetros de paginación
	page, limit := helpers.ExtractPaginationParams(c)
	offset := helpers.CalculateOffset(page, limit)

	var logs []models.AuditLog
	var total int64

	query := db.Model(&models.AuditLog{})

	// Filtrar por usuario si se especifica
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Filtrar por acción si se especifica
	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}

	// Contar total
	query.Count(&total)

	// Obtener logs con paginación
	query.Order("timestamp DESC").Offset(offset).Limit(limit).Find(&logs)

	// Usar helper para formatear respuesta con paginación
	meta := helpers.BuildPaginationMeta(page, limit, total)
	helpers.FormatPaginationResponse(c, logs, meta, "Audit logs retrieved")
}

func systemConfigEndpoint(c *gin.Context) {
	config := gin.H{
		"version":     "0.1.0",
		"environment": "development",
		"features": gin.H{
			"registration_enabled": true,
			"email_verification":   false,
			"social_login":         false,
			"two_factor_auth":      false,
			"api_rate_limit":       true,
			"file_upload":          false,
		},
		"limits": gin.H{
			"max_events_per_org":      100,
			"max_attendees_per_event": 1000,
			"max_file_size_mb":        10,
			"rate_limit_per_minute":   60,
		},
		"api": gin.H{
			"version":      "v1",
			"base_url":     "/api/v1",
			"docs_enabled": true,
		},
	}

	helpers.FormatSuccessResponse(c, gin.H{
		"data":         config,
		"generated_at": time.Now().UTC(),
	}, "System configuration retrieved")
}

// bulkVerifyOrganizations verifica organizaciones en lote (admin)
func bulkVerifyOrganizations(c *gin.Context) {
	var req struct {
		OrganizationIDs []string `json:"organization_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.FormatValidationErrorResponse(c, err.Error())
		return
	}

	db := database.GetDB()
	userCtx := handlers.GetUserContext(c)

	var updated int64
	result := db.Model(&models.Organization{}).
		Where("id IN ? AND is_verified = false", req.OrganizationIDs).
		Updates(map[string]interface{}{
			"is_verified": true,
			"verified_at": time.Now(),
			"verified_by": userCtx.ID,
			"status":      models.OrgStatusActive,
		})

	updated = result.RowsAffected

	helpers.FormatSuccessResponse(c, gin.H{
		"requested": len(req.OrganizationIDs),
		"verified":  updated,
		"timestamp": time.Now().UTC(),
	}, "Organizations verified successfully")
}

// bulkModerateEvents modera eventos en lote (admin)
func bulkModerateEvents(c *gin.Context) {
	var req struct {
		EventIDs []string `json:"event_ids" binding:"required"`
		Action   string   `json:"action" binding:"required,oneof=publish unpublish cancel"`
		Reason   string   `json:"reason,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.FormatValidationErrorResponse(c, err.Error())
		return
	}

	db := database.GetDB()
	userCtx := handlers.GetUserContext(c)

	var updated int64
	var updates map[string]interface{}

	switch req.Action {
	case "publish":
		updates = map[string]interface{}{
			"status":       models.EventStatusPublished,
			"published_at": time.Now(),
			"published_by": userCtx.ID,
		}
	case "unpublish":
		updates = map[string]interface{}{
			"status": models.EventStatusDraft,
		}
	case "cancel":
		updates = map[string]interface{}{
			"status":          models.EventStatusCanceled,
			"canceled_at":     time.Now(),
			"canceled_by":     userCtx.ID,
			"canceled_reason": req.Reason,
		}
	}

	result := db.Model(&models.Event{}).
		Where("id IN ?", req.EventIDs).
		Updates(updates)

	updated = result.RowsAffected

	// Log de auditoría para acciones masivas cuando esté implementado
	// go logBulkAction(userCtx.ID, "bulk_moderate_events", req.Action, len(req.EventIDs))

	helpers.FormatSuccessResponse(c, gin.H{
		"action":    req.Action,
		"requested": len(req.EventIDs),
		"updated":   updated,
		"timestamp": time.Now().UTC(),
	}, "Events moderated successfully")
}

// ============================================
// Funciones Helper Exportadas (Solo las necesarias)
// ============================================

// GetUserContext helper exportado para obtener user context
func GetUserContext(c *gin.Context) *common.UserContext {
	if userCtx, exists := c.Get("user_context"); exists {
		if ctx, ok := userCtx.(*common.UserContext); ok {
			return ctx
		}
	}
	return nil
}

// GetQueryOptions helper exportado para obtener query options
func GetQueryOptions(c *gin.Context) *common.QueryOptions {
	if opts, exists := c.Get("query_options"); exists {
		if queryOpts, ok := opts.(*common.QueryOptions); ok {
			return queryOpts
		}
	}
	return &common.QueryOptions{
		Page:     1,
		Limit:    20,
		OrderBy:  "created_at",
		OrderDir: "desc",
		Filters:  make(map[string]interface{}),
	}
}
