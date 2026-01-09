// internal/handlers/user_capabilities_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/mappers"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/permissions"
	"cybesphere-backend/internal/services"
)

type UserCapabilitiesHandler struct {
	authService services.AuthorizationService
	userService services.UserService
	mapper      *mappers.UnifiedMapper
}

func NewUserCapabilitiesHandler(
	authService services.AuthorizationService,
	userService services.UserService,
	mapper *mappers.UnifiedMapper,
) *UserCapabilitiesHandler {
	return &UserCapabilitiesHandler{
		authService: authService,
		userService: userService,
		mapper:      mapper,
	}
}

// GetUserCapabilities obtiene capacidades del usuario
func (h *UserCapabilitiesHandler) GetUserCapabilities(c *gin.Context) {
	userCtx := extractUserContext(c)
	if userCtx == nil {
		common.ErrorResponse(c, common.ErrUnauthorized)
		return
	}

	capabilities := h.authService.GetUserCapabilities(userCtx)
	permissions := h.authService.GetUserPermissions(userCtx)

	response := dto.UserCapabilitiesResponse{
		UserID:       userCtx.ID,
		Role:         string(userCtx.Role),
		Capabilities: capabilities,
		Permissions:  h.formatPermissions(permissions),
	}

	if userCtx.OrganizationID != nil {
		response.Organization = &dto.OrganizationSummaryResponse{
			ID: *userCtx.OrganizationID,
		}
	}

	common.SuccessResponse(c, http.StatusOK, "Capacidades del usuario", response)
}

// CheckResourceAccess verifica acceso a recurso
func (h *UserCapabilitiesHandler) CheckResourceAccess(c *gin.Context) {
	userCtx := extractUserContext(c)

	var req struct {
		Resource   string `json:"resource" binding:"required"`
		Action     string `json:"action" binding:"required"`
		ResourceID string `json:"resource_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, common.NewValidationError("request", err.Error()))
		return
	}

	var err error
	switch req.Action {
	case "read":
		err = h.authService.CheckReadPermission(userCtx, req.Resource, req.ResourceID)
	case "create":
		err = h.authService.CheckCreatePermission(userCtx, req.Resource)
	case "update":
		err = h.authService.CheckUpdatePermission(userCtx, req.Resource, req.ResourceID)
	case "delete":
		err = h.authService.CheckDeletePermission(userCtx, req.Resource, req.ResourceID)
	default:
		err = common.NewBusinessError("invalid_action", "Acción no válida")
	}

	allowed := err == nil
	response := gin.H{
		"allowed":     allowed,
		"resource":    req.Resource,
		"action":      req.Action,
		"resource_id": req.ResourceID,
	}

	if !allowed {
		response["reason"] = err.Error()
	}

	common.SuccessResponse(c, http.StatusOK, "Verificación de acceso", response)
}

// GetAvailableActions obtiene las acciones disponibles para el usuario
func (h *UserCapabilitiesHandler) GetAvailableActions(c *gin.Context) {
	userCtx := extractUserContext(c)
	if userCtx == nil {
		common.ErrorResponse(c, common.ErrUnauthorized)
		return
	}

	capabilities := h.authService.GetUserCapabilities(userCtx)

	// Filtrar solo las capacidades que están habilitadas
	actions := make([]string, 0)
	for action, allowed := range capabilities {
		if allowed {
			actions = append(actions, action)
		}
	}

	common.SuccessResponse(c, http.StatusOK, "Acciones disponibles", gin.H{
		"actions": actions,
		"count":   len(actions),
	})
}

// GetUserSessions obtiene sesiones activas del usuario actual
func (h *UserCapabilitiesHandler) GetUserSessions(c *gin.Context) {
	userCtx := extractUserContext(c)
	if userCtx == nil {
		common.ErrorResponse(c, common.ErrUnauthorized)
		return
	}

	sessions, err := h.userService.GetUserSessions(c.Request.Context(), userCtx.ID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Obtener token ID actual del contexto para marcar la sesión actual
	currentTokenID := ""
	if tokenID, exists := c.Get("token_id"); exists {
		currentTokenID = tokenID.(string)
	}

	response := h.mapper.RefreshTokensToSessionList(sessions, currentTokenID)
	common.SuccessResponse(c, http.StatusOK, "Sesiones activas", response)
}

// RevokeSession revoca una sesión específica
func (h *UserCapabilitiesHandler) RevokeSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	userCtx := extractUserContext(c)
	if userCtx == nil {
		common.ErrorResponse(c, common.ErrUnauthorized)
		return
	}

	err := h.userService.RevokeUserSession(c.Request.Context(), sessionID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Sesión revocada exitosamente", nil)
}

// GetRoleInfo obtiene información sobre los roles del sistema
func (h *UserCapabilitiesHandler) GetRoleInfo(c *gin.Context) {
	userCtx := extractUserContext(c)
	if userCtx == nil {
		common.ErrorResponse(c, common.ErrUnauthorized)
		return
	}

	// Definir información de roles
	roles := []gin.H{
		{
			"value":       "user",
			"label":       "Usuario",
			"description": "Puede ver eventos públicos y gestionar sus favoritos",
			"permissions": []string{
				"Ver eventos públicos",
				"Gestionar favoritos",
				"Ver organizaciones",
				"Editar su perfil",
			},
		},
		{
			"value":       "organizer",
			"label":       "Organizador",
			"description": "Puede crear y gestionar eventos de su organización",
			"permissions": []string{
				"Crear eventos",
				"Gestionar eventos de su organización",
				"Ver miembros de su organización",
				"Editar información de su organización",
				"Todo lo que puede hacer un usuario",
			},
		},
		{
			"value":       "admin",
			"label":       "Administrador",
			"description": "Control total del sistema",
			"permissions": []string{
				"Gestionar todos los usuarios",
				"Gestionar todas las organizaciones",
				"Gestionar todos los eventos",
				"Ver logs de auditoría",
				"Configurar el sistema",
				"Verificar organizaciones",
				"Cambiar roles de usuarios",
			},
		},
	}

	response := gin.H{
		"current_role":       string(userCtx.Role),
		"current_role_label": h.getRoleLabel(userCtx.Role),
		"roles":              roles,
		"can_change_roles":   userCtx.IsAdmin(),
	}

	common.SuccessResponse(c, http.StatusOK, "Información de roles", response)
}

// formatPermissions formatea permisos para respuesta
func (h *UserCapabilitiesHandler) formatPermissions(perms []permissions.Permission) []dto.PermissionResponse {
	formatted := make([]dto.PermissionResponse, 0, len(perms))
	for _, p := range perms {
		formatted = append(formatted, dto.PermissionResponse{
			Resource: p.Resource,
			Action:   p.Action,
			Display:  h.getPermissionDisplay(p),
		})
	}
	return formatted
}

// getPermissionDisplay obtiene descripción legible del permiso
func (h *UserCapabilitiesHandler) getPermissionDisplay(p permissions.Permission) string {
	displays := map[string]map[string]string{
		"event": {
			"read":    "Ver eventos",
			"write":   "Crear/editar eventos",
			"delete":  "Eliminar eventos",
			"publish": "Publicar eventos",
		},
		"organization": {
			"read":   "Ver organizaciones",
			"write":  "Crear/editar organizaciones",
			"delete": "Eliminar organizaciones",
			"verify": "Verificar organizaciones",
		},
		"user": {
			"read":   "Ver usuarios",
			"write":  "Editar usuarios",
			"delete": "Eliminar usuarios",
			"manage": "Gestionar usuarios",
		},
		"system": {
			"view_audit_logs": "Ver logs de auditoría",
			"manage_system":   "Gestionar sistema",
		},
	}

	if resource, ok := displays[p.Resource]; ok {
		if display, ok := resource[p.Action]; ok {
			return display
		}
	}

	return p.Resource + ":" + p.Action
}

// getRoleLabel obtiene etiqueta legible del rol
func (h *UserCapabilitiesHandler) getRoleLabel(role models.UserRole) string {
	labels := map[models.UserRole]string{
		models.RoleUser:      "Usuario",
		models.RoleOrganizer: "Organizador",
		models.RoleAdmin:     "Administrador",
	}

	if label, ok := labels[role]; ok {
		return label
	}

	return string(role)
}
