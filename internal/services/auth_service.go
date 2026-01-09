// internal/services/auth_service.go
package services

import (
	"context"
	"errors"
	"strings"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/mappers"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/repositories"
	"cybesphere-backend/pkg/auth"
	"cybesphere-backend/pkg/logger"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest, ipAddress, userAgent string) (*models.User, *auth.TokenPair, error)
	Login(ctx context.Context, req *dto.LoginRequest, ipAddress, userAgent string) (*models.User, *auth.TokenPair, error)
	RefreshTokens(ctx context.Context, refreshToken, ipAddress, userAgent string) (*auth.TokenPair, *models.User, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID string) error
}

type AuthServiceImpl struct {
	userRepo         *repositories.UserRepository
	refreshTokenRepo *repositories.RefreshTokenRepository
	jwtManager       *auth.JWTManager
	mapper           *mappers.UnifiedMapper
}

func NewAuthServiceImpl(
	userRepo *repositories.UserRepository,
	refreshTokenRepo *repositories.RefreshTokenRepository,
	jwtManager *auth.JWTManager,
	mapper *mappers.UnifiedMapper,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		mapper:           mapper,
	}
}

// Register maneja el registro completo
func (s *AuthServiceImpl) Register(ctx context.Context, req *dto.RegisterRequest, ipAddress, userAgent string) (*models.User, *auth.TokenPair, error) {
	// Log intento de registro
	logger.WithFields(map[string]interface{}{
		"email":      strings.ToLower(req.Email),
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"operation":  "register_attempt",
		"type":       "auth",
	}).Info("User registration attempt")

	// Verificar si usuario existe
	existingUser, _ := s.userRepo.GetByEmail(ctx, strings.ToLower(req.Email))
	if existingUser != nil {
		logger.LogAuth("", "register", false, "email_already_exists")
		return nil, nil, common.NewBusinessError("email_exists", "El email ya está registrado")
	}

	// Crear usuario usando mapper
	user, err := s.mapper.RegisterRequestToUser(req)
	if err != nil {
		logger.LogAuth("", "register", false, "mapping_error")
		return nil, nil, err
	}

	// Guardar usuario
	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.LogAuth("", "register", false, "database_error")
		return nil, nil, err
	}

	// Generar tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		string(user.Role),
	)
	if err != nil {
		logger.LogAuth(user.ID.String(), "register", false, "token_generation_error")
		return nil, nil, err
	}

	// Guardar refresh token
	if err := s.storeRefreshToken(ctx, user.ID.String(), tokenPair.RefreshToken, ipAddress, userAgent); err != nil {
		// Log the error but don't fail registration
		logger.WithFields(map[string]interface{}{
			"user_id":    user.ID.String(),
			"user_email": user.Email,
			"ip_address": ipAddress,
			"error":      err.Error(),
			"operation":  "store_refresh_token",
			"type":       "auth_warning",
		}).Warn("Failed to store refresh token during registration")
	}

	// Log registro exitoso
	logger.LogAuth(user.ID.String(), "register", true, "")

	return user, tokenPair, nil
}

// Login maneja el login completo
func (s *AuthServiceImpl) Login(ctx context.Context, req *dto.LoginRequest, ipAddress, userAgent string) (*models.User, *auth.TokenPair, error) {
	// Log intento de login
	logger.WithFields(map[string]interface{}{
		"email":      strings.ToLower(req.Email),
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"operation":  "login_attempt",
		"type":       "auth",
	}).Info("User login attempt")

	// Buscar usuario
	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			logger.LogAuth("", "login", false, "user_not_found")
			return nil, nil, common.NewBusinessError("invalid_credentials", "Credenciales inválidas")
		}
		logger.LogAuth("", "login", false, "database_error")
		return nil, nil, err
	}

	// Verificar password
	if !user.CheckPassword(req.Password) {
		logger.LogAuth(user.ID.String(), "login", false, "invalid_password")
		return nil, nil, common.NewBusinessError("invalid_credentials", "Credenciales inválidas")
	}

	// Verificar cuenta activa
	if !user.IsActive {
		logger.LogAuth(user.ID.String(), "login", false, "account_disabled")
		return nil, nil, common.NewBusinessError("account_disabled", "Cuenta deshabilitada")
	}

	// Actualizar último login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID.String()); err != nil {
		// Log the error but don't fail the login process
		logger.WithFields(map[string]interface{}{
			"user_id":    user.ID.String(),
			"user_email": user.Email,
			"error":      err.Error(),
			"operation":  "update_last_login",
			"type":       "auth_warning",
		}).Warn("Failed to update last login")
	}

	// Generar tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		string(user.Role),
	)
	if err != nil {
		logger.LogAuth(user.ID.String(), "login", false, "token_generation_error")
		return nil, nil, err
	}

	// Guardar refresh token
	if err := s.storeRefreshToken(ctx, user.ID.String(), tokenPair.RefreshToken, ipAddress, userAgent); err != nil {
		// Log error but don't fail login
		logger.WithFields(map[string]interface{}{
			"user_id":    user.ID.String(),
			"user_email": user.Email,
			"ip_address": ipAddress,
			"error":      err.Error(),
			"operation":  "store_refresh_token",
			"type":       "auth_warning",
		}).Warn("Failed to store refresh token during login")
	}

	// Log login exitoso
	logger.LogAuth(user.ID.String(), "login", true, "")

	return user, tokenPair, nil
}

// RefreshTokens renueva tokens
func (s *AuthServiceImpl) RefreshTokens(ctx context.Context, refreshToken, ipAddress, userAgent string) (*auth.TokenPair, *models.User, error) {
	// Validar refresh token
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		logger.LogAuth("", "refresh_token", false, "invalid_token")
		return nil, nil, common.NewBusinessError("invalid_token", "Refresh token inválido")
	}

	// Verificar token en DB
	tokenHash, _ := auth.HashRefreshToken(refreshToken)
	storedToken, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		logger.LogAuth(claims.UserID, "refresh_token", false, "token_not_found")
		return nil, nil, common.NewBusinessError("token_not_found", "Token no encontrado")
	}

	if !storedToken.IsActive() {
		logger.LogAuth(claims.UserID, "refresh_token", false, "token_expired")
		return nil, nil, common.NewBusinessError("token_expired", "Token expirado o revocado")
	}

	// Obtener usuario
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		logger.LogAuth(claims.UserID, "refresh_token", false, "user_not_found")
		return nil, nil, err
	}

	if !user.IsActive {
		logger.LogAuth(user.ID.String(), "refresh_token", false, "account_disabled")
		return nil, nil, common.NewBusinessError("account_disabled", "Cuenta deshabilitada")
	}

	// Revocar token actual
	if err := s.refreshTokenRepo.RevokeByTokenHash(ctx, tokenHash); err != nil {
		// Log the error but continue with token generation
		logger.WithFields(map[string]interface{}{
			"user_id":    user.ID.String(),
			"user_email": user.Email,
			"token_hash": tokenHash[:8] + "...", // Solo primeros 8 chars por seguridad
			"error":      err.Error(),
			"operation":  "revoke_refresh_token",
			"type":       "auth_warning",
		}).Warn("Failed to revoke old refresh token")
	}

	// Generar nuevos tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		string(user.Role),
	)
	if err != nil {
		logger.LogAuth(user.ID.String(), "refresh_token", false, "token_generation_error")
		return nil, nil, err
	}

	// Guardar nuevo refresh token
	if err := s.storeRefreshToken(ctx, user.ID.String(), tokenPair.RefreshToken, ipAddress, userAgent); err != nil {
		// Log error but don't fail the refresh process
		logger.WithFields(map[string]interface{}{
			"user_id":    user.ID.String(),
			"user_email": user.Email,
			"ip_address": ipAddress,
			"error":      err.Error(),
			"operation":  "store_new_refresh_token",
			"type":       "auth_warning",
		}).Warn("Failed to store new refresh token")
	}

	// Log refresh exitoso
	logger.LogAuth(user.ID.String(), "refresh_token", true, "")

	return tokenPair, user, nil
}

// Logout revoca un refresh token
func (s *AuthServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	tokenHash, err := auth.HashRefreshToken(refreshToken)
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"error":     err.Error(),
			"operation": "logout",
			"type":      "auth_error",
		}).Error("Failed to hash refresh token during logout")
		return err
	}

	if err := s.refreshTokenRepo.RevokeByTokenHash(ctx, tokenHash); err != nil {
		logger.WithFields(map[string]interface{}{
			"token_hash": tokenHash[:8] + "...",
			"error":      err.Error(),
			"operation":  "logout",
			"type":       "auth_error",
		}).Error("Failed to revoke refresh token during logout")
		return err
	}

	logger.WithFields(map[string]interface{}{
		"operation": "logout",
		"type":      "auth",
	}).Info("User logout successful")

	return nil
}

// LogoutAll revoca todos los tokens del usuario
func (s *AuthServiceImpl) LogoutAll(ctx context.Context, userID string) error {
	if err := s.refreshTokenRepo.RevokeAllByUserID(ctx, userID); err != nil {
		logger.WithFields(map[string]interface{}{
			"user_id":   userID,
			"error":     err.Error(),
			"operation": "logout_all",
			"type":      "auth_error",
		}).Error("Failed to revoke all refresh tokens")
		return err
	}

	logger.LogAuth(userID, "logout_all", true, "")
	return nil
}

// storeRefreshToken guarda un refresh token
func (s *AuthServiceImpl) storeRefreshToken(ctx context.Context, userID, refreshToken, ipAddress, userAgent string) error {
	tokenHash, err := auth.HashRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	claims, err := s.jwtManager.GetTokenClaims(refreshToken)
	if err != nil {
		return err
	}

	token := &models.RefreshToken{
		UserID:     userID,
		TokenHash:  tokenHash,
		TokenID:    claims.ID,
		ExpiresAt:  claims.ExpiresAt.Time,
		UserAgent:  userAgent,
		IPAddress:  ipAddress,
		DeviceInfo: s.extractDeviceInfo(userAgent),
	}

	return s.refreshTokenRepo.Create(ctx, token)
}

// extractDeviceInfo extrae información del dispositivo
func (s *AuthServiceImpl) extractDeviceInfo(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		return "Mobile"
	}
	if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		return "Tablet"
	}

	return "Desktop"
}
