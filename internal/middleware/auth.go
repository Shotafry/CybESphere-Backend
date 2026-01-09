// internal/middleware/auth.go - VERSIÓN MODERNIZADA
package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/permissions"
	"cybesphere-backend/pkg/auth"
	"cybesphere-backend/pkg/database"
	"cybesphere-backend/pkg/logger"
)

// AuthMiddleware middleware de autenticación MODERNIZADO
type AuthMiddleware struct {
	jwtManager        *auth.JWTManager
	permissionChecker *permissions.PermissionChecker
	db                *gorm.DB
}

// NewAuthMiddleware crea una nueva instancia del middleware
func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager:        jwtManager,
		permissionChecker: permissions.NewPermissionChecker(),
		db:                database.GetDB(),
	}
}

// =============================================================================
// MÉTODOS EXISTENTES (mantenidos por compatibilidad)
// =============================================================================

// RequireAuth middleware que requiere autenticación válida
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.ErrorResponse(c, common.NewBusinessError("missing_token", "Token de autorización requerido"))
			c.Abort()
			return
		}

		tokenString, err := m.jwtManager.ExtractTokenFromHeader(authHeader)
		if err != nil {
			common.ErrorResponse(c, common.NewBusinessError("invalid_token_format", "Formato de token inválido"))
			c.Abort()
			return
		}

		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			var errorCode, errorMessage string
			switch err {
			case auth.ErrExpiredToken:
				errorCode, errorMessage = "token_expired", "Token expirado"
			case auth.ErrInvalidToken:
				errorCode, errorMessage = "invalid_token", "Token inválido"
			case auth.ErrInvalidTokenType:
				errorCode, errorMessage = "invalid_token_type", "Tipo de token inválido"
			case auth.ErrInvalidClaims:
				errorCode, errorMessage = "invalid_claims", "Claims del token inválidos"
			default:
				errorCode, errorMessage = "token_validation_failed", "Error al validar token"
			}

			logger.Warnf("Token validation failed: %v", err)
			common.ErrorResponse(c, common.NewBusinessError(errorCode, errorMessage))
			c.Abort()
			return
		}

		// Verificar que el usuario siga activo
		var user models.User
		if err := m.db.First(&user, "id = ?", claims.UserID).Error; err != nil {
			logger.Warnf("User not found for valid token: %s", claims.UserID)
			common.ErrorResponse(c, common.NewBusinessError("user_not_found", "Usuario no encontrado"))
			c.Abort()
			return
		}

		if !user.IsActive {
			common.ErrorResponse(c, common.NewBusinessError("account_disabled", "Cuenta deshabilitada"))
			c.Abort()
			return
		}

		// Almacenar información del usuario en el contexto
		m.setUserContext(c, claims, &user)

		c.Next()
	}
}

// OptionalAuth middleware que extrae información del usuario si está presente
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString, err := m.jwtManager.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Verificar que el usuario siga activo
		var user models.User
		if err := m.db.First(&user, "id = ?", claims.UserID).Error; err != nil || !user.IsActive {
			c.Next()
			return
		}

		// Token válido, establecer información del usuario
		m.setUserContext(c, claims, &user)
		c.Set("authenticated", true)

		c.Next()
	}
}

// =============================================================================
// NUEVOS MÉTODOS MODERNIZADOS
// =============================================================================

// InjectUserContext inyecta UserContext unificado
func (m *AuthMiddleware) InjectUserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := m.extractUserContext(c)
		if userContext != nil {
			c.Set("user_context", userContext)
		}
		c.Next()
	}
}

// ParseQueryOptions parsea y valida QueryOptions automáticamente
func (m *AuthMiddleware) ParseQueryOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		opts := &common.QueryOptions{}

		// Bind query parameters
		if err := c.ShouldBindQuery(opts); err != nil {
			common.ErrorResponse(c, common.NewValidationError("query_params", "Parámetros de consulta inválidos"))
			c.Abort()
			return
		}

		// Validar y normalizar
		if err := opts.Validate(); err != nil {
			common.ErrorResponse(c, err)
			c.Abort()
			return
		}

		// Inyectar UserContext
		opts.UserContext = m.extractUserContext(c)

		c.Set("query_options", opts)
		c.Next()
	}
}

// GuardResource protege un recurso específico con lógica de ownership
func (m *AuthMiddleware) GuardResource(resourceType, resourceIDParam string, permission permissions.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := m.extractUserContext(c)
		if userContext == nil {
			common.ErrorResponse(c, common.ErrUnauthorized)
			c.Abort()
			return
		}

		resourceID := c.Param(resourceIDParam)

		// Verificar acceso al recurso
		if !m.permissionChecker.CanAccessResource(
			userContext.ID,
			userContext.Role,
			permission,
			resourceType,
			resourceID,
		) {
			reason := m.permissionChecker.GetAccessDenialReason(
				userContext.Role,
				permission,
				resourceType,
				resourceID,
			)

			common.ErrorResponse(c, common.NewBusinessError("access_denied", reason))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermissionEnhanced versión moderna de RequirePermission
func (m *AuthMiddleware) RequirePermissionEnhanced(permission permissions.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := m.extractUserContext(c)
		if userContext == nil {
			common.ErrorResponse(c, common.ErrUnauthorized)
			c.Abort()
			return
		}

		if !userContext.HasPermission(permission.Resource, permission.Action) {
			message := fmt.Sprintf("No tienes permisos para %s en %s", permission.Action, permission.Resource)
			common.ErrorResponse(c, common.NewBusinessError("insufficient_permissions", message))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthFlow cadena completa de autenticación moderna
func (m *AuthMiddleware) AuthFlow() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Autenticación
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// 2. Inyectar UserContext
		m.InjectUserContext()(c)

		// 3. Parsear QueryOptions
		m.ParseQueryOptions()(c)
		if c.IsAborted() {
			return
		}

		c.Next()
	}
}

// OptionalAuthFlow para endpoints que no requieren autenticación obligatoria
func (m *AuthMiddleware) OptionalAuthFlow() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Autenticación opcional
		m.OptionalAuth()(c)

		// 2. Si hay usuario, inyectar contexto
		if _, exists := c.Get("user_id"); exists {
			m.InjectUserContext()(c)
		}

		// 3. Parsear QueryOptions
		m.ParseQueryOptions()(c)
		if c.IsAborted() {
			return
		}

		c.Next()
	}
}

// =============================================================================
// GUARD METHODS ESPECÍFICOS PARA RECURSOS
// =============================================================================

// GuardEvent protege recursos de eventos
func (m *AuthMiddleware) GuardEvent(permission permissions.Permission) gin.HandlerFunc {
	return m.GuardResource("event", "eventId", permission)
}

// GuardOrganization protege recursos de organizaciones
func (m *AuthMiddleware) GuardOrganization(permission permissions.Permission) gin.HandlerFunc {
	return m.GuardResource("organization", "orgId", permission)
}

// GuardUser protege recursos de usuarios
func (m *AuthMiddleware) GuardUser(permission permissions.Permission) gin.HandlerFunc {
	return m.GuardResource("user", "userId", permission)
}

// RequireEventOwnership verifica ownership específico de eventos
func (m *AuthMiddleware) RequireEventOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := m.extractUserContext(c)
		if userContext == nil {
			common.ErrorResponse(c, common.ErrUnauthorized)
			c.Abort()
			return
		}

		// Admin puede acceder a todo
		if userContext.IsAdmin() {
			c.Next()
			return
		}

		eventID := c.Param("eventId")
		if eventID == "" {
			// Para creación de eventos, verificar organización
			m.checkEventCreationPermissions(c, userContext)
			return
		}

		// Para eventos existentes, verificar ownership
		m.checkEventOwnership(c, userContext, eventID)
	}
}

// RequireOrganizationOwnership verifica ownership específico de organizaciones
func (m *AuthMiddleware) RequireOrganizationOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := m.extractUserContext(c)
		if userContext == nil {
			common.ErrorResponse(c, common.ErrUnauthorized)
			c.Abort()
			return
		}

		// Admin puede acceder a todo
		if userContext.IsAdmin() {
			c.Next()
			return
		}

		orgID := c.Param("orgId")
		if orgID == "" {
			// Para creación, verificar que esté verificado
			if !userContext.IsVerified {
				common.ErrorResponse(c, common.NewBusinessError("verification_required",
					"Debes tener una cuenta verificada para crear organizaciones"))
				c.Abort()
				return
			}
			c.Next()
			return
		}

		if !userContext.CanManageOrganization(orgID) {
			common.ErrorResponse(c, common.NewBusinessError("organization_access_denied",
				"Solo puedes gestionar tu propia organización"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// =============================================================================
// MIDDLEWARE FACTORY METHODS
// =============================================================================

// ForPublicEndpoint middleware para endpoints públicos
func (m *AuthMiddleware) ForPublicEndpoint() gin.HandlerFunc {
	return m.OptionalAuthFlow()
}

// ForAuthenticatedEndpoint middleware para endpoints autenticados
func (m *AuthMiddleware) ForAuthenticatedEndpoint() gin.HandlerFunc {
	return m.AuthFlow()
}

// ForEventManagement middleware completo para gestión de eventos
func (m *AuthMiddleware) ForEventManagement() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.AuthFlow()(c)
		if c.IsAborted() {
			return
		}

		m.RequirePermissionEnhanced(permissions.WriteEvent)(c)
		if c.IsAborted() {
			return
		}

		m.RequireEventOwnership()(c)
		if c.IsAborted() {
			return
		}

		c.Next()
	}
}

// ForOrganizationManagement middleware completo para gestión de organizaciones
func (m *AuthMiddleware) ForOrganizationManagement() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.AuthFlow()(c)
		if c.IsAborted() {
			return
		}

		m.RequirePermissionEnhanced(permissions.WriteOrganization)(c)
		if c.IsAborted() {
			return
		}

		m.RequireOrganizationOwnership()(c)
		if c.IsAborted() {
			return
		}

		c.Next()
	}
}

// ForAdminOnly middleware para endpoints de solo administrador
func (m *AuthMiddleware) ForAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.AuthFlow()(c)
		if c.IsAborted() {
			return
		}

		userContext := m.extractUserContext(c)
		if !userContext.IsAdmin() {
			common.ErrorResponse(c, common.NewBusinessError("admin_required",
				"Solo administradores pueden acceder a este recurso"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// ForReadOnlyEndpoint middleware para endpoints de solo lectura
func (m *AuthMiddleware) ForReadOnlyEndpoint(resourceType string) gin.HandlerFunc {
	var permission permissions.Permission

	switch resourceType {
	case "event":
		permission = permissions.ReadEvent
	case "organization":
		permission = permissions.ReadOrganization
	case "user":
		permission = permissions.ReadProfile
	default:
		permission = permissions.Permission{Resource: resourceType, Action: "read"}
	}

	return func(c *gin.Context) {
		m.OptionalAuthFlow()(c)
		if c.IsAborted() {
			return
		}

		// Si hay usuario autenticado, verificar permisos
		userContext := m.extractUserContext(c)
		if userContext != nil && !userContext.HasPermission(permission.Resource, permission.Action) {
			common.ErrorResponse(c, common.NewBusinessError("read_permission_denied",
				"No tienes permisos para leer este tipo de recurso"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// setUserContext establece la información del usuario en el contexto
func (m *AuthMiddleware) setUserContext(c *gin.Context, claims *auth.Claims, user *models.User) {
	c.Set("user_id", claims.UserID)
	c.Set("user_email", claims.Email)
	c.Set("user_role", claims.Role)
	c.Set("token_id", claims.TokenID)
	c.Set("user_active", user.IsActive)
	c.Set("user_verified", user.IsVerified)

	if user.OrganizationID != nil {
		c.Set("organization_id", *user.OrganizationID)
	}
}

// extractUserContext convierte contexto de Gin a UserContext unificado
func (m *AuthMiddleware) extractUserContext(c *gin.Context) *common.UserContext {
	userID, hasUser := c.Get("user_id")
	if !hasUser {
		return nil
	}

	userEmail, _ := c.Get("user_email")
	userRole, _ := c.Get("user_role")
	isActive, _ := c.Get("user_active")
	isVerified, _ := c.Get("user_verified")
	orgID, hasOrg := c.Get("organization_id")

	role := models.UserRole(userRole.(string))

	userContext := &common.UserContext{
		ID:           userID.(string),
		Email:        userEmail.(string),
		Role:         role,
		IsActive:     isActive.(bool),
		IsVerified:   isVerified.(bool),
		Permissions:  m.permissionChecker.GetRolePermissions(role),
		Capabilities: m.permissionChecker.GetRoleCapabilities(role),
	}

	if hasOrg {
		orgIDStr := orgID.(string)
		userContext.OrganizationID = &orgIDStr
	}

	return userContext
}

// checkEventCreationPermissions verifica permisos para crear eventos
func (m *AuthMiddleware) checkEventCreationPermissions(c *gin.Context, userContext *common.UserContext) {
	if !userContext.IsOrganizer() || userContext.OrganizationID == nil {
		common.ErrorResponse(c, common.NewBusinessError("no_organization",
			"Debes pertenecer a una organización para crear eventos"))
		c.Abort()
		return
	}

	// Verificar que la organización puede crear eventos
	var org models.Organization
	if err := m.db.First(&org, "id = ?", *userContext.OrganizationID).Error; err != nil {
		common.ErrorResponse(c, common.MapGormError(err))
		c.Abort()
		return
	}

	if !org.CanCreateEvent() {
		common.ErrorResponse(c, common.NewBusinessError("organization_cannot_create_events",
			"Tu organización no puede crear eventos en este momento"))
		c.Abort()
		return
	}

	c.Next()
}

// checkEventOwnership verifica ownership de eventos existentes
func (m *AuthMiddleware) checkEventOwnership(c *gin.Context, userContext *common.UserContext, eventID string) {
	var event models.Event
	if err := m.db.First(&event, "id = ?", eventID).Error; err != nil {
		common.ErrorResponse(c, common.MapGormError(err))
		c.Abort()
		return
	}

	if !userContext.CanManageOrganization(event.OrganizationID) {
		common.ErrorResponse(c, common.NewBusinessError("event_access_denied",
			"Solo puedes gestionar eventos de tu organización"))
		c.Abort()
		return
	}

	c.Next()
}

// =============================================================================
// MÉTODOS LEGACY (mantenidos por compatibilidad)
// =============================================================================

// RequireRole middleware que requiere roles específicos (LEGACY)
func (m *AuthMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			common.ErrorResponse(c, common.ErrUnauthorized)
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		common.ErrorResponse(c, common.NewBusinessError("insufficient_permissions",
			"Permisos insuficientes para esta acción"))
		c.Abort()
	}
}

// RequireAdmin middleware que requiere rol de administrador (LEGACY)
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}

// RequireOrganizerOrAdmin middleware que requiere rol de organizador o admin (LEGACY)
func (m *AuthMiddleware) RequireOrganizerOrAdmin() gin.HandlerFunc {
	return m.RequireRole("admin", "organizer")
}

// RequirePermission middleware que verifica permisos específicos (LEGACY)
func (m *AuthMiddleware) RequirePermission(permission permissions.Permission) gin.HandlerFunc {
	return m.RequirePermissionEnhanced(permission)
}

// RequireResourceAccess verifica acceso a recurso específico (LEGACY)
func (m *AuthMiddleware) RequireResourceAccess(resourceType, resourceIDParam string, permission permissions.Permission) gin.HandlerFunc {
	return m.GuardResource(resourceType, resourceIDParam, permission)
}

// RequireOrganizationMember verifica membresía de organización (LEGACY)
func (m *AuthMiddleware) RequireOrganizationMember(orgIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userContext := m.extractUserContext(c)
		if userContext == nil {
			common.ErrorResponse(c, common.ErrUnauthorized)
			c.Abort()
			return
		}

		orgID := c.Param(orgIDParam)

		// Admin puede acceder a todo
		if userContext.IsAdmin() {
			c.Next()
			return
		}

		if userContext.OrganizationID != nil && *userContext.OrganizationID == orgID {
			c.Next()
			return
		}

		common.ErrorResponse(c, common.NewBusinessError("organization_access_denied",
			"No perteneces a esta organización"))
		c.Abort()
	}
}

// AuditLog middleware para logging de auditoría (MANTENIDO)
func (m *AuthMiddleware) AuditLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		userID, _ := c.Get("user_id")
		userEmail, _ := c.Get("user_email")
		userRole, _ := c.Get("user_role")

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		if userID != nil {
			logger.Infof("Audit: User %s (%s, %s) %s %s - Status: %d - Duration: %v - IP: %s",
				userID, userEmail, userRole, method, path, statusCode, duration, c.ClientIP())

			if m.isCriticalAction(method, path, statusCode) {
				go m.logCriticalAction(userID.(string), method, path, statusCode, c.ClientIP(), c.GetHeader("User-Agent"))
			}
		} else {
			logger.Debugf("Public: %s %s - Status: %d - Duration: %v - IP: %s",
				method, path, statusCode, duration, c.ClientIP())
		}
	}
}

// InjectCapabilities inyecta capacidades del usuario (LEGACY pero actualizado)
func (m *AuthMiddleware) InjectCapabilities() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.Next()
			return
		}

		role := models.UserRole(userRole.(string))
		capabilities := m.permissionChecker.GetRoleCapabilities(role)
		permissions := m.permissionChecker.GetRolePermissions(role)

		c.Set("user_capabilities", capabilities)
		c.Set("user_permissions", permissions)
		c.Next()
	}
}

// GetJWTManager returns the JWT manager instance (MANTENIDO)
func (m *AuthMiddleware) GetJWTManager() *auth.JWTManager {
	return m.jwtManager
}

// Helper methods (MANTENIDOS)
func (m *AuthMiddleware) isCriticalAction(method, path string, statusCode int) bool {
	criticalPaths := []string{"/users/", "/admin/", "/auth/logout-all"}
	criticalMethods := []string{"POST", "PUT", "DELETE"}

	for _, criticalPath := range criticalPaths {
		if strings.Contains(path, criticalPath) {
			for _, criticalMethod := range criticalMethods {
				if method == criticalMethod {
					return true
				}
			}
		}
	}

	return statusCode == 401 || statusCode == 403
}

func (m *AuthMiddleware) logCriticalAction(userID, method, path string, statusCode int, ipAddress, userAgent string) {
	auditLog := models.AuditLog{
		UserID:    userID,
		Action:    method,
		Resource:  path,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Timestamp: time.Now(),
		Status:    statusCode,
	}

	if err := m.db.Create(&auditLog).Error; err != nil {
		logger.Errorf("Failed to log critical action: %v", err)
	}
}
