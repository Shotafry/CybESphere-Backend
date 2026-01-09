// internal/handlers/auth_handler.go
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/mappers"
	"cybesphere-backend/internal/services"
	"cybesphere-backend/pkg/auth"
)

type AuthHandler struct {
	authService services.AuthService
	userService services.UserService
	mapper      *mappers.UnifiedMapper
	jwtManager  *auth.JWTManager
}

func NewAuthHandler(
	authService services.AuthService,
	userService services.UserService,
	mapper *mappers.UnifiedMapper,
	jwtManager *auth.JWTManager,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		mapper:      mapper,
		jwtManager:  jwtManager,
	}
}

// Register maneja el registro usando services
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, common.NewValidationError("request", err.Error()))
		return
	}

	// Delegar toda la lógica al service
	user, tokenPair, err := h.authService.Register(c.Request.Context(), &req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Usar mapper para la respuesta
	response := h.mapper.UserToRegisterResponse(
		user,
		tokenPair.AccessToken,
		tokenPair.RefreshToken,
		int(time.Until(tokenPair.AccessTokenExpiresAt).Seconds()),
	)

	common.SuccessResponse(c, http.StatusCreated, response.Message, response)
}

// Login simplificado delegando a services
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, common.NewValidationError("request", err.Error()))
		return
	}

	// Delegar al service
	user, tokenPair, err := h.authService.Login(
		c.Request.Context(),
		&req,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Usar mapper para respuesta
	response := h.mapper.UserToAuthResponse(
		user,
		tokenPair.AccessToken,
		tokenPair.RefreshToken,
		int(time.Until(tokenPair.AccessTokenExpiresAt).Seconds()),
	)

	common.SuccessResponse(c, http.StatusOK, "Login exitoso", response)
}

// RefreshToken simplificado
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, common.NewValidationError("request", err.Error()))
		return
	}

	// Delegar al service
	tokenPair, user, err := h.authService.RefreshTokens(
		c.Request.Context(),
		req.RefreshToken,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
	)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	response := h.mapper.BuildTokenResponse(
		tokenPair.AccessToken,
		tokenPair.RefreshToken,
		tokenPair.AccessTokenExpiresAt,
	)

	common.SuccessResponse(c, http.StatusOK, "Tokens renovados", gin.H{
		"data": response,
		"user": h.mapper.UserToSummaryResponse(user),
	})
}

// Me obtiene info del usuario actual usando service
func (h *AuthHandler) Me(c *gin.Context) {
	userCtx := extractUserContext(c)
	if userCtx == nil {
		common.ErrorResponse(c, common.ErrUnauthorized)
		return
	}

	user, err := h.userService.Get(c.Request.Context(), userCtx.ID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	response := h.mapper.UserToDetailResponse(user, userCtx)
	common.SuccessResponse(c, http.StatusOK, "Información del usuario", response)
}

// Logout usando service
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, common.NewValidationError("request", err.Error()))
		return
	}

	err := h.authService.Logout(c.Request.Context(), req.RefreshToken)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	response := h.mapper.BuildLogoutResponse(true)
	common.SuccessResponse(c, http.StatusOK, response.Message, nil)
}

// LogoutAll usando service
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userCtx := extractUserContext(c)
	if userCtx == nil {
		common.ErrorResponse(c, common.ErrUnauthorized)
		return
	}

	err := h.authService.LogoutAll(c.Request.Context(), userCtx.ID)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Todas las sesiones cerradas", nil)
}
