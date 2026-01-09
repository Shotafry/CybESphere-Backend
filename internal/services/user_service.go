// internal/services/user_service.go
package services

import (
	"context"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/repositories"
)

// UserServiceImpl implementación concreta del servicio de usuarios
type UserServiceImpl struct {
	*BaseService[models.User, dto.CreateUserRequest, dto.UpdateUserRequest]
	userRepo         *repositories.UserRepository
	refreshTokenRepo *repositories.RefreshTokenRepository
	auth             AuthorizationService
}

// Verificación en tiempo de compilación
var _ UserService = (*UserServiceImpl)(nil)

// NewUserService crea nueva instancia del servicio de usuarios
func NewUserService(
	userRepo *repositories.UserRepository,
	refreshTokenRepo *repositories.RefreshTokenRepository,
	mapper ResponseMapper,
	auth AuthorizationService,
) UserService {
	base := NewBaseService[models.User, dto.CreateUserRequest, dto.UpdateUserRequest](
		userRepo, mapper, auth,
	)

	return &UserServiceImpl{
		BaseService:      base,
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		auth:             auth,
	}
}

// GetByEmail obtiene usuario por email
func (s *UserServiceImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// UpdateRole actualiza el rol de un usuario (solo admin)
func (s *UserServiceImpl) UpdateRole(ctx context.Context, userID string, newRole models.UserRole, userCtx *common.UserContext) error {
	if !userCtx.IsAdmin() {
		return common.NewBusinessError("admin_required", "Solo administradores pueden cambiar roles")
	}

	// Validar transición de rol
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// No permitir auto-degradación de admin
	if userCtx.ID == userID && user.Role == models.RoleAdmin && newRole != models.RoleAdmin {
		return common.NewBusinessError("self_demotion_denied", "No puedes degradar tu propio rol de administrador")
	}

	return s.userRepo.UpdateRole(ctx, userID, newRole)
}

// DeactivateUser desactiva un usuario (solo admin)
func (s *UserServiceImpl) DeactivateUser(ctx context.Context, userID string, userCtx *common.UserContext) error {
	if !userCtx.IsAdmin() {
		return common.NewBusinessError("admin_required", "Solo administradores pueden desactivar usuarios")
	}

	// No permitir auto-desactivación
	if userCtx.ID == userID {
		return common.NewBusinessError("self_deactivation_denied", "No puedes desactivar tu propia cuenta")
	}

	// Desactivar usuario
	if err := s.userRepo.DeactivateUser(ctx, userID); err != nil {
		return err
	}

	// Revocar todos los tokens del usuario
	return s.refreshTokenRepo.RevokeAllByUserID(ctx, userID)
}

// ActivateUser reactiva un usuario
func (s *UserServiceImpl) ActivateUser(ctx context.Context, userID string, userCtx *common.UserContext) error {
	if !userCtx.IsAdmin() {
		return common.NewBusinessError("admin_required", "Solo administradores pueden activar usuarios")
	}

	return s.userRepo.ActivateUser(ctx, userID)
}

// GetUserSessions obtiene sesiones activas del usuario
func (s *UserServiceImpl) GetUserSessions(ctx context.Context, userID string, userCtx *common.UserContext) ([]*models.RefreshToken, error) {
	// Solo el propio usuario o admin pueden ver sesiones
	if userCtx.ID != userID && !userCtx.IsAdmin() {
		return nil, common.NewBusinessError("access_denied", "No puedes ver las sesiones de otro usuario")
	}

	return s.refreshTokenRepo.GetActiveByUserID(ctx, userID)
}

// RevokeUserSession revoca una sesión específica
func (s *UserServiceImpl) RevokeUserSession(ctx context.Context, sessionID string, userCtx *common.UserContext) error {
	// Verificar que la sesión pertenece al usuario o es admin
	token, err := s.refreshTokenRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}

	if token.UserID != userCtx.ID && !userCtx.IsAdmin() {
		return common.NewBusinessError("access_denied", "No puedes revocar sesiones de otro usuario")
	}

	return s.refreshTokenRepo.RevokeByTokenHash(ctx, token.TokenHash)
}
