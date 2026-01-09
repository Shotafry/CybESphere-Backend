package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config estructura principal de configuración
type Config struct {
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	JWT        JWTConfig        `json:"jwt"`
	Security   SecurityConfig   `json:"security"`
	Logging    LoggingConfig    `json:"logging"`
	Monitoring MonitoringConfig `json:"monitoring"`
	Email      EmailConfig      `json:"email"`
	Upload     UploadConfig     `json:"upload"`
	Geo        GeoConfig        `json:"geo"`
	RateLimit  RateLimitConfig  `json:"rate_limit"`
}

// ServerConfig configuración del servidor
type ServerConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Mode       string `json:"mode"`
	Version    string `json:"version"`
	EnableDocs bool   `json:"enable_docs"`
}

// DatabaseConfig configuración de base de datos
type DatabaseConfig struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	Name            string        `json:"name"`
	User            string        `json:"user"`
	Password        string        `json:"-"` // No exponer en JSON
	SSLMode         string        `json:"ssl_mode"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	LogLevel        string        `json:"log_level"`
}

// JWTConfig configuración JWT
type JWTConfig struct {
	Secret               string        `json:"-"` // No exponer en JSON
	AccessTokenDuration  time.Duration `json:"access_token_duration"`
	RefreshTokenDuration time.Duration `json:"refresh_token_duration"`
	Issuer               string        `json:"issuer"`
}

// SecurityConfig configuración de seguridad
type SecurityConfig struct {
	BcryptCost         int      `json:"bcrypt_cost"`
	CORSAllowedOrigins []string `json:"cors_allowed_origins"`
	CORSAllowedMethods []string `json:"cors_allowed_methods"`
	CORSAllowedHeaders []string `json:"cors_allowed_headers"`
	CORSEnabled        bool     `json:"cors_enabled"`
	TrustedProxies     []string `json:"trusted_proxies"`
}

// LoggingConfig configuración de logging
type LoggingConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	Output     string `json:"output"`
	FilePath   string `json:"file_path"`
	AuditLevel string `json:"audit_level"`
}

// MonitoringConfig configuración de monitoreo
type MonitoringConfig struct {
	Environment        string `json:"environment"`
	DebugMode          bool   `json:"debug_mode"`
	HealthCheckEnabled bool   `json:"health_check_enabled"`
	MetricsEnabled     bool   `json:"metrics_enabled"`
	ProfilingEnabled   bool   `json:"profiling_enabled"`
}

// EmailConfig configuración de email
type EmailConfig struct {
	SMTPEnabled  bool   `json:"smtp_enabled"`
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"-"` // No exponer en JSON
	FromEmail    string `json:"from_email"`
	FromName     string `json:"from_name"`
}

// UploadConfig configuración de archivos
type UploadConfig struct {
	MaxSizeMB         int      `json:"max_size_mb"`
	AllowedExtensions []string `json:"allowed_extensions"`
	Path              string   `json:"path"`
}

// GeoConfig configuración geoespacial
type GeoConfig struct {
	DefaultRadiusKM int `json:"default_radius_km"`
	MaxRadiusKM     int `json:"max_radius_km"`
}

// RateLimitConfig configuración de rate limiting
type RateLimitConfig struct {
	Enabled           bool `json:"enabled"`
	RequestsPerMinute int  `json:"requests_per_minute"`
	Burst             int  `json:"burst"`
}

// Load carga la configuración desde variables de entorno
func Load() (*Config, error) {
	// Cargar .env si existe
	if err := godotenv.Load(); err != nil {
		// No es error crítico si no existe .env
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	config := &Config{
		Server: ServerConfig{
			Host:       getEnvString("SERVER_HOST", "localhost"),
			Port:       getEnvInt("SERVER_PORT", 8080),
			Mode:       getEnvString("SERVER_MODE", "release"),
			Version:    getEnvString("API_VERSION", "v1"),
			EnableDocs: getEnvBool("ENABLE_DOCS", true),
		},
		Database: DatabaseConfig{
			Host:            getEnvString("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			Name:            getEnvString("DB_NAME", "cybesphere_dev"),
			User:            getEnvString("DB_USER", "cybesphere_user"),
			Password:        getEnvString("DB_PASSWORD", ""),
			SSLMode:         getEnvString("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", "300s"),
			LogLevel:        getEnvString("DB_LOG_LEVEL", "silent"),
		},
		JWT: JWTConfig{
			Secret:               getEnvString("JWT_SECRET", ""),
			AccessTokenDuration:  getEnvDuration("JWT_EXPIRATION", "15m"),
			RefreshTokenDuration: getEnvDuration("JWT_REFRESH_EXPIRATION", "168h"),
			Issuer:               getEnvString("JWT_ISSUER", "cybesphere-api"),
		},
		Security: SecurityConfig{
			BcryptCost:         getEnvInt("BCRYPT_COST", 12),
			CORSAllowedOrigins: getEnvStringSlice("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173"),
			CORSAllowedMethods: getEnvStringSlice("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			CORSAllowedHeaders: getEnvStringSlice("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization,X-Requested-With"),
			CORSEnabled:        getEnvBool("CORS_ENABLED", true),         // <-- NUEVO
			TrustedProxies:     getEnvStringSlice("TRUSTED_PROXIES", ""), // <-- NUEVO
		},
		Logging: LoggingConfig{
			Level:      getEnvString("LOG_LEVEL", "warn"),
			Format:     getEnvString("LOG_FORMAT", "json"),
			Output:     getEnvString("LOG_OUTPUT", "stdout"),
			FilePath:   getEnvString("LOG_FILE_PATH", "logs/app.log"),
			AuditLevel: getEnvString("AUDIT_LOG_LEVEL", "info"),
		},
		Monitoring: MonitoringConfig{
			Environment:        getEnvString("ENVIRONMENT", "development"),
			DebugMode:          getEnvBool("DEBUG_MODE", false),
			HealthCheckEnabled: getEnvBool("HEALTH_CHECK_ENABLED", true),
			MetricsEnabled:     getEnvBool("METRICS_ENABLED", true),
			ProfilingEnabled:   getEnvBool("PROFILING_ENABLED", false),
		},
		Email: EmailConfig{
			SMTPEnabled:  getEnvBool("SMTP_ENABLED", false),
			SMTPHost:     getEnvString("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvInt("SMTP_PORT", 587),
			SMTPUsername: getEnvString("SMTP_USERNAME", ""),
			SMTPPassword: getEnvString("SMTP_PASSWORD", ""),
			FromEmail:    getEnvString("SMTP_FROM_EMAIL", "noreply@cybesphere.com"),
			FromName:     getEnvString("SMTP_FROM_NAME", "CybESphere"),
		},
		Upload: UploadConfig{
			MaxSizeMB:         getEnvInt("UPLOAD_MAX_SIZE_MB", 10),
			AllowedExtensions: getEnvStringSlice("UPLOAD_ALLOWED_EXTENSIONS", "jpg,jpeg,png,gif,pdf,doc,docx"),
			Path:              getEnvString("UPLOAD_PATH", "uploads/"),
		},
		Geo: GeoConfig{
			DefaultRadiusKM: getEnvInt("GEO_DEFAULT_RADIUS_KM", 50),
			MaxRadiusKM:     getEnvInt("GEO_MAX_RADIUS_KM", 200),
		},
		RateLimit: RateLimitConfig{
			Enabled:           getEnvBool("RATE_LIMIT_ENABLED", true),
			RequestsPerMinute: getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 100),
			Burst:             getEnvInt("RATE_LIMIT_BURST", 20),
		},
	}

	// Validaciones
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate valida la configuración
func (c *Config) Validate() error {
	// Validar JWT Secret
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	// Validar configuración de base de datos
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	// Validar puerto del servidor
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("SERVER_PORT must be between 1 and 65535")
	}

	// Validar modo de servidor
	validModes := map[string]bool{"debug": true, "release": true, "test": true}
	if !validModes[c.Server.Mode] {
		return fmt.Errorf("SERVER_MODE must be one of: debug, release, test")
	}

	return nil
}

// GetAddress retorna la dirección completa del servidor
func (s *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// GetDSN retorna el Data Source Name para PostgreSQL
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		d.Host, d.User, d.Password, d.Name, d.Port, d.SSLMode,
	)
}

// Funciones auxiliares para obtener variables de entorno

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue string) time.Duration {
	value := getEnvString(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	// Si no se puede parsear, usar valor por defecto
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return 15 * time.Minute // fallback final
}

func getEnvStringSlice(key string, defaultValue string) []string {
	value := getEnvString(key, defaultValue)
	if value == "" {
		return []string{}
	}

	// Dividir por comas y limpiar espacios
	var result []string
	for _, item := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
