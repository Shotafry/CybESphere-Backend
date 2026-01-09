package permissions

import (
	"fmt"

	"cybesphere-backend/internal/models"
	"cybesphere-backend/pkg/database"

	"gorm.io/gorm"
)

// Permission define una acción específica sobre un recurso
type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// Permisos predefinidos del sistema
var (
	// User permissions
	ReadProfile   = Permission{"user", "read"}
	WriteProfile  = Permission{"user", "write"}
	DeleteProfile = Permission{"user", "delete"}

	// Organization permissions
	ReadOrganization   = Permission{"organization", "read"}
	WriteOrganization  = Permission{"organization", "write"}
	DeleteOrganization = Permission{"organization", "delete"}
	ManageOrganization = Permission{"organization", "manage"}

	// Event permissions
	ReadEvent       = Permission{"event", "read"}
	WriteEvent      = Permission{"event", "write"}
	DeleteEvent     = Permission{"event", "delete"}
	PublishEvent    = Permission{"event", "publish"}
	ManageAttendees = Permission{"event", "manage_attendees"}

	// System permissions
	ManageUsers       = Permission{"system", "manage_users"}
	ManageSystem      = Permission{"system", "manage_system"}
	ViewAuditLogs     = Permission{"system", "view_audit_logs"}
	ManagePermissions = Permission{"system", "manage_permissions"}
)

// RolePermissions define los permisos para cada rol
var RolePermissions = map[models.UserRole][]Permission{
	models.RoleUser: {
		ReadProfile, WriteProfile, ReadEvent, ReadOrganization,
	},
	models.RoleOrganizer: {
		ReadProfile, WriteProfile, ReadEvent, WriteEvent, DeleteEvent,
		PublishEvent, ManageAttendees, ReadOrganization, WriteOrganization, ManageOrganization,
	},
	models.RoleAdmin: {
		ReadProfile, WriteProfile, DeleteProfile, ReadEvent, WriteEvent, DeleteEvent,
		PublishEvent, ManageAttendees, ReadOrganization, WriteOrganization, DeleteOrganization,
		ManageOrganization, ManageUsers, ManageSystem, ViewAuditLogs, ManagePermissions,
	},
}

// PermissionChecker interfaz para verificar permisos
type PermissionChecker struct {
	db *gorm.DB
}

// NewPermissionChecker crea una nueva instancia
func NewPermissionChecker() *PermissionChecker {
	return &PermissionChecker{
		db: database.GetDB(),
	}
}

// HasPermission verifica si un rol tiene un permiso específico
func (pc *PermissionChecker) HasPermission(role models.UserRole, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p.Resource == permission.Resource && p.Action == permission.Action {
			return true
		}
	}
	return false
}

// CanAccessResource verifica si un usuario puede acceder a un recurso específico
func (pc *PermissionChecker) CanAccessResource(userID string, userRole models.UserRole, permission Permission, resourceType, resourceID string) bool {
	// Verificar permiso básico
	if !pc.HasPermission(userRole, permission) {
		return false
	}

	// Admin puede acceder a todo
	if userRole == models.RoleAdmin {
		return true
	}

	// Para recursos con ownership, verificar ownership
	switch resourceType {
	case "event":
		return pc.checkEventAccess(userID, resourceID, userRole, permission.Action)
	case "organization":
		return pc.checkOrganizationAccess(userID, resourceID, userRole, permission.Action)
	case "user":
		// Los usuarios solo pueden gestionar su propio perfil (excepto admin)
		return userID == resourceID
	default:
		// Para recursos no-ownership, el permiso básico es suficiente
		return true
	}
}

// checkEventAccess verifica acceso a eventos
func (pc *PermissionChecker) checkEventAccess(userID, eventID string, userRole models.UserRole, action string) bool {
	if userRole != models.RoleOrganizer {
		return false
	}

	// Para creación de eventos (eventID vacío), verificar que el usuario tenga organización
	if eventID == "" {
		var user models.User
		if err := pc.db.First(&user, "id = ?", userID).Error; err != nil {
			return false
		}

		// Debe tener organización activa
		if user.OrganizationID == nil {
			return false
		}

		// Verificar que la organización puede crear eventos
		var org models.Organization
		if err := pc.db.First(&org, "id = ?", *user.OrganizationID).Error; err != nil {
			return false
		}

		return org.CanCreateEvents
	}

	// Para eventos existentes, verificar ownership
	var event models.Event
	if err := pc.db.First(&event, "id = ?", eventID).Error; err != nil {
		return false
	}

	var user models.User
	if err := pc.db.First(&user, "id = ?", userID).Error; err != nil {
		return false
	}

	return user.OrganizationID != nil && *user.OrganizationID == event.OrganizationID
}

// checkOrganizationAccess verifica acceso a organizaciones
func (pc *PermissionChecker) checkOrganizationAccess(userID, orgID string, userRole models.UserRole, action string) bool {
	if userRole != models.RoleOrganizer {
		return false
	}

	// Para creación de organizaciones (orgID vacío), verificar que el usuario esté verificado
	if orgID == "" {
		var user models.User
		if err := pc.db.First(&user, "id = ?", userID).Error; err != nil {
			return false
		}
		return user.IsVerified
	}

	// Para organizaciones existentes, verificar membership
	var user models.User
	if err := pc.db.First(&user, "id = ?", userID).Error; err != nil {
		return false
	}

	return user.OrganizationID != nil && *user.OrganizationID == orgID
}

// GetRolePermissions obtiene todos los permisos de un rol
func (pc *PermissionChecker) GetRolePermissions(role models.UserRole) []Permission {
	if permissions, exists := RolePermissions[role]; exists {
		return permissions
	}
	return []Permission{}
}

// GetRoleHierarchy returns the role hierarchy level
func (pc *PermissionChecker) GetRoleHierarchy(role models.UserRole) int {
	hierarchy := map[models.UserRole]int{
		models.RoleUser:      1,
		models.RoleOrganizer: 2,
		models.RoleAdmin:     3,
	}

	if level, exists := hierarchy[role]; exists {
		return level
	}
	return 0
}

// CanAccessRole checks if requester role can access target role
func (pc *PermissionChecker) CanAccessRole(requesterRole, targetRole models.UserRole) bool {
	return pc.GetRoleHierarchy(requesterRole) >= pc.GetRoleHierarchy(targetRole)
}

// GetAvailableRoles returns roles that the current role can assign
func (pc *PermissionChecker) GetAvailableRoles(currentRole models.UserRole) []models.UserRole {
	currentLevel := pc.GetRoleHierarchy(currentRole)
	var availableRoles []models.UserRole

	allRoles := []models.UserRole{
		models.RoleUser,
		models.RoleOrganizer,
		models.RoleAdmin,
	}

	for _, role := range allRoles {
		if pc.GetRoleHierarchy(role) <= currentLevel {
			availableRoles = append(availableRoles, role)
		}
	}

	return availableRoles
}

// ValidateRoleTransition checks if role transition is allowed
func (pc *PermissionChecker) ValidateRoleTransition(fromRole, toRole, requesterRole models.UserRole) error {
	// Only admin can change roles to/from admin
	if (fromRole == models.RoleAdmin || toRole == models.RoleAdmin) && requesterRole != models.RoleAdmin {
		return fmt.Errorf("only admin can manage admin role")
	}

	// Requester must have higher or equal hierarchy to both roles
	requesterLevel := pc.GetRoleHierarchy(requesterRole)
	fromLevel := pc.GetRoleHierarchy(fromRole)
	toLevel := pc.GetRoleHierarchy(toRole)

	if requesterLevel < fromLevel || requesterLevel < toLevel {
		return fmt.Errorf("insufficient privileges for role transition")
	}

	return nil
}

// GetRoleDescription returns human-readable role description
func (pc *PermissionChecker) GetRoleDescription(role models.UserRole, language string) string {
	descriptions := map[models.UserRole]map[string]string{
		models.RoleUser: {
			"en": "Regular user with basic access to view events and manage personal profile",
			"es": "Usuario regular con acceso básico para ver eventos y gestionar perfil personal",
		},
		models.RoleOrganizer: {
			"en": "Organization member who can create and manage events for their organization",
			"es": "Miembro de organización que puede crear y gestionar eventos de su organización",
		},
		models.RoleAdmin: {
			"en": "System administrator with full access to all platform features",
			"es": "Administrador del sistema con acceso completo a todas las funciones",
		},
	}

	if desc, exists := descriptions[role]; exists {
		if text, langExists := desc[language]; langExists {
			return text
		}
		// Fallback to English
		if text, langExists := desc["en"]; langExists {
			return text
		}
	}

	return string(role)
}

// UserContext provides complete user context for permission checks
type UserContext struct {
	ID             string          `json:"id"`
	Role           models.UserRole `json:"role"`
	OrganizationID *string         `json:"organization_id,omitempty"`
	Capabilities   map[string]bool `json:"capabilities"`
	Permissions    []Permission    `json:"permissions"`
}

// GetUserContext builds complete user context
func (pc *PermissionChecker) GetUserContext(userID string) (*UserContext, error) {
	var user models.User
	if err := pc.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &UserContext{
		ID:             user.ID.String(),
		Role:           user.Role,
		OrganizationID: user.OrganizationID,
		Capabilities:   pc.GetRoleCapabilities(user.Role),
		Permissions:    pc.GetRolePermissions(user.Role),
	}, nil
}

// GetRoleCapabilities obtiene las capacidades de un rol
func (pc *PermissionChecker) GetRoleCapabilities(role models.UserRole) map[string]bool {
	capabilities := map[string]bool{
		// User capabilities
		"can_read_profile":   pc.HasPermission(role, ReadProfile),
		"can_write_profile":  pc.HasPermission(role, WriteProfile),
		"can_delete_profile": pc.HasPermission(role, DeleteProfile),

		// Event capabilities
		"can_read_events":      pc.HasPermission(role, ReadEvent),
		"can_write_events":     pc.HasPermission(role, WriteEvent),
		"can_delete_events":    pc.HasPermission(role, DeleteEvent),
		"can_publish_events":   pc.HasPermission(role, PublishEvent),
		"can_manage_attendees": pc.HasPermission(role, ManageAttendees),

		// Organization capabilities
		"can_read_organizations":   pc.HasPermission(role, ReadOrganization),
		"can_write_organizations":  pc.HasPermission(role, WriteOrganization),
		"can_delete_organizations": pc.HasPermission(role, DeleteOrganization),
		"can_manage_organizations": pc.HasPermission(role, ManageOrganization),

		// System capabilities
		"can_manage_users":       pc.HasPermission(role, ManageUsers),
		"can_manage_system":      pc.HasPermission(role, ManageSystem),
		"can_view_audit_logs":    pc.HasPermission(role, ViewAuditLogs),
		"can_manage_permissions": pc.HasPermission(role, ManagePermissions),
	}

	return capabilities
}

// GetAccessDenialReason obtiene la razón de denegación de acceso
func (pc *PermissionChecker) GetAccessDenialReason(role models.UserRole, permission Permission, resourceType, resourceID string) string {
	// Verificar si el rol tiene el permiso básico
	if !pc.HasPermission(role, permission) {
		return "Tu rol no tiene permisos para esta acción"
	}

	// Si tiene el permiso pero no puede acceder al recurso específico
	switch resourceType {
	case "event":
		if role == models.RoleOrganizer {
			if resourceID == "" {
				return "Debes pertenecer a una organización activa para crear eventos"
			}
			return "Solo puedes gestionar eventos de tu organización"
		}
	case "organization":
		if role == models.RoleOrganizer {
			if resourceID == "" {
				return "Debes tener una cuenta verificada para crear organizaciones"
			}
			return "Solo puedes gestionar tu propia organización"
		}
	case "user":
		if role != models.RoleAdmin {
			return "Solo puedes gestionar tu propio perfil"
		}
	}

	return "Acceso denegado"
}
