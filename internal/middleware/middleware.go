// internal/middleware/middleware.go
package middleware

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/config"
	"cybesphere-backend/pkg/logger"
)

// CORSMiddleware configura CORS usando la configuración del .env
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Usar orígenes permitidos de la configuración
		isAllowed := false
		for _, allowedOrigin := range cfg.Security.CORSAllowedOrigins {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed || origin == "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", strings.Join(cfg.Security.CORSAllowedHeaders, ", "))
		c.Header("Access-Control-Allow-Methods", strings.Join(cfg.Security.CORSAllowedMethods, ", "))
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SecurityHeaders middleware para headers de seguridad
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevenir sniffing de MIME types
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevenir clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Protección XSS
		c.Header("X-XSS-Protection", "1; mode=block")

		// Política de referrer
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy básica
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self' https:; media-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'")

		// HSTS (solo en HTTPS)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// Prevenir información del servidor
		c.Header("Server", "CybESphere-API")

		c.Next()
	}
}

// RequestID middleware para agregar ID único a cada request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		c.Next()
	}
}

// RateLimitMiddleware usando configuración del .env
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	if !cfg.RateLimit.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Almacén en memoria simple (en producción usar Redis)
	requestCounts := make(map[string][]time.Time)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		// Limpiar requests antiguos
		if requests, exists := requestCounts[clientIP]; exists {
			var validRequests []time.Time
			cutoff := now.Add(-time.Minute)

			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}
			requestCounts[clientIP] = validRequests
		}

		// Verificar límite usando configuración
		currentRequests := len(requestCounts[clientIP])
		if currentRequests >= cfg.RateLimit.RequestsPerMinute {
			c.JSON(429, gin.H{
				"error":       "rate_limit_exceeded",
				"message":     "Demasiadas requests. Intenta nuevamente en unos minutos.",
				"retry_after": 60,
			})
			c.Abort()
			return
		}

		// Agregar request actual
		requestCounts[clientIP] = append(requestCounts[clientIP], now)

		c.Next()
	}
}

// RequestLogger middleware personalizado para logging de requests
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Formato personalizado para logs HTTP
		return fmt.Sprintf("[%s] %d %s %s %s %v %s \"%s\"\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Method,
			param.Path,
			param.ClientIP,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// ErrorHandler middleware para manejo centralizado de errores
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Si hay errores en el contexto, manejarlos
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Log del error
			logger.Errorf("Request error: %v - Path: %s - IP: %s",
				err.Error(), c.Request.URL.Path, c.ClientIP())

			// Si no se ha enviado respuesta aún
			if c.Writer.Status() == 200 {
				c.JSON(500, gin.H{
					"error":      "internal_server_error",
					"message":    "Error interno del servidor",
					"request_id": c.GetString("request_id"),
				})
			}
		}
	}
}

// RecoveryWithLogger middleware personalizado de recovery
func RecoveryWithLogger() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, recovered interface{}) {
		// Log detallado del panic
		logger.Errorf("Panic recovered: %v - Path: %s - IP: %s - User-Agent: %s",
			recovered, c.Request.URL.Path, c.ClientIP(), c.GetHeader("User-Agent"))

		requestID := c.GetString("request_id")

		c.JSON(500, gin.H{
			"error":      "internal_server_error",
			"message":    "Error interno del servidor",
			"request_id": requestID,
		})
	})
}

// Timeout middleware para requests con timeout
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if timeout <= 0 {
			c.Next()
			return
		}

		// Crear contexto con timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Establecer contexto con timeout
		c.Request = c.Request.WithContext(ctx)

		// Canal para detectar si la request terminó
		done := make(chan struct{})

		go func() {
			defer close(done)
			c.Next()
		}()

		select {
		case <-done:
			// Request completada normalmente
			return
		case <-ctx.Done():
			// Timeout alcanzado
			if ctx.Err() != nil {
				logger.Warnf("Request timeout: %s - IP: %s", c.Request.URL.Path, c.ClientIP())
				c.JSON(408, gin.H{
					"error":   "request_timeout",
					"message": "La request tardó demasiado en completarse",
				})
				c.Abort()
			}
		}
	}
}

// BasicUserContext middleware para enriquecer el contexto con información básica del usuario para logging
func BasicUserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Si hay información de usuario, enriquecer el contexto para logs
		if userID, exists := c.Get("user_id"); exists {
			// Crear contexto básico de usuario para logs
			logContext := map[string]interface{}{
				"user_id": userID,
			}

			if userEmail, exists := c.Get("user_email"); exists {
				logContext["email"] = userEmail
			}

			if userRole, exists := c.Get("user_role"); exists {
				logContext["role"] = userRole
			}

			if orgID, exists := c.Get("organization_id"); exists {
				logContext["organization_id"] = orgID
			}

			c.Set("log_context", logContext)
		}

		c.Next()
	}
}

// ValidationError middleware para formatear errores de validación
func ValidationError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Procesar errores de binding/validación
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Si es un error de validación de Gin
			if err.Type == gin.ErrorTypeBind {
				c.JSON(400, gin.H{
					"error":   "validation_failed",
					"message": "Datos de entrada inválidos",
					"details": err.Error(),
				})
				return
			}
		}
	}
}

// APIVersioning middleware para versionado de API
func APIVersioning() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer versión del path o header
		version := extractAPIVersion(c)
		c.Set("api_version", version)

		// Validar versión soportada
		if !isVersionSupported(version) {
			c.JSON(400, gin.H{
				"error":              "unsupported_api_version",
				"message":            "Versión de API no soportada",
				"supported_versions": []string{"v1"},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ContentTypeValidation middleware para validar Content-Type
func ContentTypeValidation(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Solo validar para métodos que envían datos
		if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "PATCH" {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")

		// Si no hay tipos permitidos especificados, permitir cualquiera
		if len(allowedTypes) == 0 {
			c.Next()
			return
		}

		// Verificar si el Content-Type está permitido
		for _, allowedType := range allowedTypes {
			if strings.Contains(contentType, allowedType) {
				c.Next()
				return
			}
		}

		c.JSON(400, gin.H{
			"error":         "invalid_content_type",
			"message":       "Content-Type no soportado",
			"allowed_types": allowedTypes,
		})
		c.Abort()
	}
}

// Helper functions

func generateRequestID() string {
	// Generar ID único usando timestamp + random
	timestamp := time.Now().UnixNano()
	return strconv.FormatInt(timestamp, 36) + generateRandomString(6)
}

func generateRandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
	}

	return string(result)
}

func extractAPIVersion(c *gin.Context) string {
	// Intentar extraer versión del path
	path := c.Request.URL.Path
	pathParts := strings.Split(path, "/")

	for _, part := range pathParts {
		if strings.HasPrefix(part, "v") && len(part) > 1 {
			return part
		}
	}

	// Si no se encuentra en el path, usar header
	if version := c.GetHeader("API-Version"); version != "" {
		return version
	}

	// Versión por defecto
	return "v1"
}

func isVersionSupported(version string) bool {
	supportedVersions := []string{"v1"}

	for _, supported := range supportedVersions {
		if version == supported {
			return true
		}
	}

	return false
}
