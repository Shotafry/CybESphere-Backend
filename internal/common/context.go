package common

import (
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/permissions"
)

// UserContext contexto unificado del usuario
type UserContext struct {
	ID             string                   `json:"id"`
	Email          string                   `json:"email"`
	Role           models.UserRole          `json:"role"`
	OrganizationID *string                  `json:"organization_id,omitempty"`
	Permissions    []permissions.Permission `json:"permissions"`
	Capabilities   map[string]bool          `json:"capabilities"`
	IsActive       bool                     `json:"is_active"`
	IsVerified     bool                     `json:"is_verified"`
}

// HasPermission verifica si el usuario tiene un permiso específico
func (uc *UserContext) HasPermission(resource, action string) bool {
	for _, perm := range uc.Permissions {
		if perm.Resource == resource && perm.Action == action {
			return true
		}
	}
	return false
}

// IsAdmin verifica si el usuario es administrador
func (uc *UserContext) IsAdmin() bool {
	return uc.Role == models.RoleAdmin
}

// IsOrganizer verifica si el usuario es organizador
func (uc *UserContext) IsOrganizer() bool {
	return uc.Role == models.RoleOrganizer
}

// CanManageOrganization verifica si puede gestionar una organización específica
func (uc *UserContext) CanManageOrganization(orgID string) bool {
	if uc.IsAdmin() {
		return true
	}
	return uc.IsOrganizer() && uc.OrganizationID != nil && *uc.OrganizationID == orgID
}
