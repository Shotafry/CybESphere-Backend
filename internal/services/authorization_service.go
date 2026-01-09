// internal/services/authorization_service.go
package services

import (
	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/permissions"
)

// AuthorizationServiceImpl implementación del servicio de autorización
type AuthorizationServiceImpl struct {
	permissionChecker *permissions.PermissionChecker
}

// NewAuthorizationService crea nueva instancia del servicio de autorización
func NewAuthorizationService() AuthorizationService {
	return &AuthorizationServiceImpl{
		permissionChecker: permissions.NewPermissionChecker(),
	}
}

// CheckReadPermission verifica permisos de lectura
func (s *AuthorizationServiceImpl) CheckReadPermission(userCtx *common.UserContext, resourceType, resourceID string) error {
	if userCtx == nil {
		// Permitir lectura pública para ciertos recursos
		if s.isPublicReadResource(resourceType) {
			return nil
		}
		return common.ErrUnauthorized
	}

	var permission permissions.Permission
	switch resourceType {
	case "event":
		permission = permissions.ReadEvent
	case "organization":
		permission = permissions.ReadOrganization
	case "user":
		permission = permissions.ReadProfile
	default:
		return common.NewBusinessError("unknown_resource", "Tipo de recurso desconocido")
	}

	if !s.permissionChecker.CanAccessResource(userCtx.ID, userCtx.Role, permission, resourceType, resourceID) {
		reason := s.permissionChecker.GetAccessDenialReason(userCtx.Role, permission, resourceType, resourceID)
		return common.NewBusinessError("access_denied", reason)
	}

	return nil
}

// CheckCreatePermission verifica permisos de creación
func (s *AuthorizationServiceImpl) CheckCreatePermission(userCtx *common.UserContext, resourceType string) error {
	if userCtx == nil {
		return common.ErrUnauthorized
	}

	var permission permissions.Permission
	switch resourceType {
	case "event":
		permission = permissions.WriteEvent
	case "organization":
		permission = permissions.WriteOrganization
	case "user":
		permission = permissions.WriteProfile
	default:
		return common.NewBusinessError("unknown_resource", "Tipo de recurso desconocido")
	}

	if !s.permissionChecker.CanAccessResource(userCtx.ID, userCtx.Role, permission, resourceType, "") {
		reason := s.permissionChecker.GetAccessDenialReason(userCtx.Role, permission, resourceType, "")
		return common.NewBusinessError("access_denied", reason)
	}

	return nil
}

// CheckUpdatePermission verifica permisos de actualización
func (s *AuthorizationServiceImpl) CheckUpdatePermission(userCtx *common.UserContext, resourceType, resourceID string) error {
	if userCtx == nil {
		return common.ErrUnauthorized
	}

	var permission permissions.Permission
	switch resourceType {
	case "event":
		permission = permissions.WriteEvent
	case "organization":
		permission = permissions.WriteOrganization
	case "user":
		permission = permissions.WriteProfile
	default:
		return common.NewBusinessError("unknown_resource", "Tipo de recurso desconocido")
	}

	if !s.permissionChecker.CanAccessResource(userCtx.ID, userCtx.Role, permission, resourceType, resourceID) {
		reason := s.permissionChecker.GetAccessDenialReason(userCtx.Role, permission, resourceType, resourceID)
		return common.NewBusinessError("access_denied", reason)
	}

	return nil
}

// CheckDeletePermission verifica permisos de eliminación
func (s *AuthorizationServiceImpl) CheckDeletePermission(userCtx *common.UserContext, resourceType, resourceID string) error {
	if userCtx == nil {
		return common.ErrUnauthorized
	}

	var permission permissions.Permission
	switch resourceType {
	case "event":
		permission = permissions.DeleteEvent
	case "organization":
		permission = permissions.DeleteOrganization
	case "user":
		permission = permissions.DeleteProfile
	default:
		return common.NewBusinessError("unknown_resource", "Tipo de recurso desconocido")
	}

	if !s.permissionChecker.CanAccessResource(userCtx.ID, userCtx.Role, permission, resourceType, resourceID) {
		reason := s.permissionChecker.GetAccessDenialReason(userCtx.Role, permission, resourceType, resourceID)
		return common.NewBusinessError("access_denied", reason)
	}

	return nil
}

// ApplySecurityFilters aplica filtros de seguridad según el contexto del usuario
func (s *AuthorizationServiceImpl) ApplySecurityFilters(opts *common.QueryOptions, userCtx *common.UserContext, resourceType string) {
	// Si no hay usuario, aplicar filtros públicos
	if userCtx == nil {
		s.applyPublicFilters(opts, resourceType)
		return
	}

	// Aplicar filtros según el rol
	switch userCtx.Role {
	case models.RoleAdmin:
		// Admin puede ver todo, no aplicar filtros adicionales
		return

	case models.RoleOrganizer:
		s.applyOrganizerFilters(opts, userCtx, resourceType)

	case models.RoleUser:
		s.applyUserFilters(opts, userCtx, resourceType)

	default:
		// Rol desconocido, aplicar filtros más restrictivos
		s.applyPublicFilters(opts, resourceType)
	}
}

// applyPublicFilters aplica filtros para acceso público
func (s *AuthorizationServiceImpl) applyPublicFilters(opts *common.QueryOptions, resourceType string) {
	switch resourceType {
	case "event":
		// Solo eventos publicados y públicos
		opts.AddFilter("status", models.EventStatusPublished)
		opts.AddFilter("is_public", true)

	case "organization":
		// Solo organizaciones activas
		opts.AddFilter("status", models.OrgStatusActive)

	case "user":
		// Usuarios públicos no tienen acceso a listados
		opts.AddFilter("id", "00000000-0000-0000-0000-000000000000") // ID que no existe
	}
}

// applyOrganizerFilters aplica filtros para organizadores
func (s *AuthorizationServiceImpl) applyOrganizerFilters(opts *common.QueryOptions, userCtx *common.UserContext, resourceType string) {
	switch resourceType {
	case "event":
		// Eventos públicos + sus propios eventos
		if userCtx.OrganizationID != nil {
			// Esta lógica se podría mejorar con una query más compleja
			// Por ahora, no aplicar filtros restrictivos para organizadores
			opts.AddFilter("status", models.EventStatusPublished) // Al menos publicados
		} else {
			// Organizador sin organización, solo eventos públicos
			s.applyPublicFilters(opts, resourceType)
		}

	case "organization":
		// Puede ver organizaciones activas
		opts.AddFilter("status", models.OrgStatusActive)

	case "user":
		// Solo usuarios de su organización si tiene una
		if userCtx.OrganizationID != nil {
			opts.AddFilter("organization_id", *userCtx.OrganizationID)
		} else {
			opts.AddFilter("id", "00000000-0000-0000-0000-000000000000") // No ver usuarios
		}
	}
}

// applyUserFilters aplica filtros para usuarios regulares
func (s *AuthorizationServiceImpl) applyUserFilters(opts *common.QueryOptions, userCtx *common.UserContext, resourceType string) {
	switch resourceType {
	case "event":
		// Solo eventos publicados y públicos
		opts.AddFilter("status", models.EventStatusPublished)
		opts.AddFilter("is_public", true)

	case "organization":
		// Solo organizaciones activas y verificadas
		opts.AddFilter("status", models.OrgStatusActive)
		opts.AddFilter("is_verified", true)

	case "user":
		// Usuario regular no puede ver listados de usuarios
		opts.AddFilter("id", userCtx.ID) // Solo su propio perfil
	}
}

// isPublicReadResource determina si un recurso permite lectura pública
func (s *AuthorizationServiceImpl) isPublicReadResource(resourceType string) bool {
	publicResources := map[string]bool{
		"event":        true,  // Eventos públicos son visibles
		"organization": true,  // Organizaciones activas son visibles
		"user":         false, // Perfiles de usuario requieren auth
	}

	return publicResources[resourceType]
}

// =============================================================================
// MÉTODOS HELPER PARA VALIDACIONES ESPECÍFICAS
// =============================================================================

// CanUserManageEvent verifica si un usuario puede gestionar un evento específico
func (s *AuthorizationServiceImpl) CanUserManageEvent(userCtx *common.UserContext, eventID string) bool {
	if userCtx == nil {
		return false
	}

	return s.permissionChecker.CanAccessResource(
		userCtx.ID,
		userCtx.Role,
		permissions.WriteEvent,
		"event",
		eventID,
	)
}

// CanUserManageOrganization verifica si un usuario puede gestionar una organización específica
func (s *AuthorizationServiceImpl) CanUserManageOrganization(userCtx *common.UserContext, orgID string) bool {
	if userCtx == nil {
		return false
	}

	return s.permissionChecker.CanAccessResource(
		userCtx.ID,
		userCtx.Role,
		permissions.WriteOrganization,
		"organization",
		orgID,
	)
}

// RequireEventOwnership valida que el usuario pueda gestionar el evento
func (s *AuthorizationServiceImpl) RequireEventOwnership(userCtx *common.UserContext, eventID string) error {
	if !s.CanUserManageEvent(userCtx, eventID) {
		return common.NewBusinessError("event_access_denied", "No tienes permisos para gestionar este evento")
	}
	return nil
}

// RequireOrganizationOwnership valida que el usuario pueda gestionar la organización
func (s *AuthorizationServiceImpl) RequireOrganizationOwnership(userCtx *common.UserContext, orgID string) error {
	if !s.CanUserManageOrganization(userCtx, orgID) {
		return common.NewBusinessError("organization_access_denied", "No tienes permisos para gestionar esta organización")
	}
	return nil
}

// ValidateRoleTransition valida transición de roles (para admin)
func (s *AuthorizationServiceImpl) ValidateRoleTransition(userCtx *common.UserContext, fromRole, toRole models.UserRole) error {
	if !userCtx.IsAdmin() {
		return common.NewBusinessError("admin_required", "Solo administradores pueden cambiar roles")
	}

	return s.permissionChecker.ValidateRoleTransition(fromRole, toRole, userCtx.Role)
}

// GetUserCapabilities obtiene las capacidades de un usuario
func (s *AuthorizationServiceImpl) GetUserCapabilities(userCtx *common.UserContext) map[string]bool {
	if userCtx == nil {
		return map[string]bool{}
	}

	return s.permissionChecker.GetRoleCapabilities(userCtx.Role)
}

// GetUserPermissions obtiene los permisos de un usuario
func (s *AuthorizationServiceImpl) GetUserPermissions(userCtx *common.UserContext) []permissions.Permission {
	if userCtx == nil {
		return []permissions.Permission{}
	}

	return s.permissionChecker.GetRolePermissions(userCtx.Role)
}

// =============================================================================
// VALIDACIONES DE DOMINIO ESPECÍFICAS
// =============================================================================

// ValidateEventCreation valida reglas específicas para creación de eventos
func (s *AuthorizationServiceImpl) ValidateEventCreation(userCtx *common.UserContext, organizationID string) error {
	if userCtx == nil {
		return common.ErrUnauthorized
	}

	// Admin puede crear eventos para cualquier organización
	if userCtx.IsAdmin() {
		return nil
	}

	// Organizador debe tener organización
	if userCtx.OrganizationID == nil {
		return common.NewBusinessError("no_organization", "Debes pertenecer a una organización para crear eventos")
	}

	// Organizador solo puede crear eventos para su organización
	if !userCtx.IsAdmin() && *userCtx.OrganizationID != organizationID {
		return common.NewBusinessError("organization_mismatch", "Solo puedes crear eventos para tu organización")
	}

	return nil
}

// ValidateOrganizationCreation valida reglas específicas para creación de organizaciones
func (s *AuthorizationServiceImpl) ValidateOrganizationCreation(userCtx *common.UserContext) error {
	if userCtx == nil {
		return common.ErrUnauthorized
	}

	// Admin puede crear cualquier organización
	if userCtx.IsAdmin() {
		return nil
	}

	// Usuario debe estar verificado
	if !userCtx.IsVerified {
		return common.NewBusinessError("verification_required", "Debes tener una cuenta verificada para crear organizaciones")
	}

	return nil
}

// CanAccessAuditLogs verifica si el usuario puede acceder a logs de auditoría
func (s *AuthorizationServiceImpl) CanAccessAuditLogs(userCtx *common.UserContext) bool {
	if userCtx == nil {
		return false
	}

	return userCtx.HasPermission("system", "view_audit_logs")
}

// CanManageSystem verifica si el usuario puede gestionar el sistema
func (s *AuthorizationServiceImpl) CanManageSystem(userCtx *common.UserContext) bool {
	if userCtx == nil {
		return false
	}

	return userCtx.HasPermission("system", "manage_system")
}

// =============================================================================
// BUILDER PATTERN PARA CONFIGURACIÓN DE AUTORIZACION
// =============================================================================

// AuthorizationBuilder builder para configurar autorizaciones complejas
type AuthorizationBuilder struct {
	service *AuthorizationServiceImpl
	rules   []AuthorizationRule
}

// AuthorizationRule regla de autorización personalizada
type AuthorizationRule struct {
	ResourceType string
	Action       string
	Condition    func(*common.UserContext, string) bool
	ErrorMessage string
}

// NewAuthorizationBuilder crea nuevo builder
func NewAuthorizationBuilder(service *AuthorizationServiceImpl) *AuthorizationBuilder {
	return &AuthorizationBuilder{
		service: service,
		rules:   make([]AuthorizationRule, 0),
	}
}

// AddRule agrega regla personalizada
func (b *AuthorizationBuilder) AddRule(resourceType, action string, condition func(*common.UserContext, string) bool, errorMessage string) *AuthorizationBuilder {
	b.rules = append(b.rules, AuthorizationRule{
		ResourceType: resourceType,
		Action:       action,
		Condition:    condition,
		ErrorMessage: errorMessage,
	})
	return b
}

// ValidateRules valida todas las reglas configuradas
func (b *AuthorizationBuilder) ValidateRules(userCtx *common.UserContext, resourceType, action, resourceID string) error {
	for _, rule := range b.rules {
		if rule.ResourceType == resourceType && rule.Action == action {
			if !rule.Condition(userCtx, resourceID) {
				return common.NewBusinessError("authorization_failed", rule.ErrorMessage)
			}
		}
	}
	return nil
}
