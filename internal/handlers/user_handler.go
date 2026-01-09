// internal/handlers/user_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/mappers"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/services"
)

type UserHandler struct {
	*BaseHandler[models.User, dto.CreateUserRequest, dto.UpdateUserRequest, dto.UserResponse]
	userService services.UserService
	mapper      *mappers.UnifiedMapper
}

func NewUserHandler(
	userService services.UserService,
	mapper *mappers.UnifiedMapper,
) *UserHandler {
	base := NewBaseHandler[models.User, dto.CreateUserRequest, dto.UpdateUserRequest, dto.UserResponse](
		userService,
		mapper,
		"user",
	)

	base.SetLimits(20, 100)

	return &UserHandler{
		BaseHandler: base,
		userService: userService,
		mapper:      mapper,
	}
}

// UpdateRole actualiza el rol del usuario (admin only)
func (h *UserHandler) UpdateRole(c *gin.Context) {
	userID := c.Param("id")
	userCtx := extractUserContext(c)

	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, common.NewValidationError("request", err.Error()))
		return
	}

	err := h.userService.UpdateRole(c.Request.Context(), userID, models.UserRole(req.Role), userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Rol actualizado exitosamente", gin.H{
		"user_id":  userID,
		"new_role": req.Role,
	})
}

// ActivateUser activa un usuario (admin only)
func (h *UserHandler) ActivateUser(c *gin.Context) {
	userID := c.Param("id")
	userCtx := extractUserContext(c)

	err := h.userService.ActivateUser(c.Request.Context(), userID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Usuario activado exitosamente", nil)
}

// DeactivateUser desactiva un usuario (admin only)
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")
	userCtx := extractUserContext(c)

	err := h.userService.DeactivateUser(c.Request.Context(), userID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Usuario desactivado exitosamente", nil)
}

// GetUserSessions obtiene sesiones del usuario
func (h *UserHandler) GetUserSessions(c *gin.Context) {
	userID := c.Param("id")
	userCtx := extractUserContext(c)

	sessions, err := h.userService.GetUserSessions(c.Request.Context(), userID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Usar mapper para convertir a respuesta
	response := h.mapper.RefreshTokensToSessionList(sessions, "")
	common.SuccessResponse(c, http.StatusOK, "Sesiones del usuario", response)
}

// GetUserProfile obtiene perfil completo (override de GetByID para más detalle)
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("id")
	userCtx := extractUserContext(c)

	user, err := h.userService.Get(c.Request.Context(), userID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Usar respuesta detallada
	response := h.mapper.UserToDetailResponse(user, userCtx)
	common.SuccessResponse(c, http.StatusOK, "Perfil de usuario", response)
}
