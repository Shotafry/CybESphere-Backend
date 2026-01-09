package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// RefreshToken modelo para almacenar refresh tokens en la base de datos
type RefreshToken struct {
	BaseModel

	// Relación con usuario
	UserID string `json:"user_id" gorm:"not null;size:36;index" validate:"required"`
	User   User   `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`

	// Token data
	TokenHash string `json:"-" gorm:"not null;size:255;uniqueIndex"` // Hash del token, no el token real
	TokenID   string `json:"token_id" gorm:"not null;size:36;index"` // ID único del token JWT

	// Expiración y estado
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" gorm:"index"`
	IsRevoked bool       `json:"is_revoked" gorm:"not null;default:false;index"`

	// Información de sesión
	UserAgent string `json:"user_agent" gorm:"size:500"`
	IPAddress string `json:"ip_address" gorm:"size:45;index"`

	// Metadata opcional
	DeviceInfo string     `json:"device_info,omitempty" gorm:"size:200"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

// TableName especifica el nombre de tabla
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// BeforeCreate hook de GORM para validación
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if err := rt.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Validar datos requeridos
	if err := rt.ValidateRefreshToken(); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate hook de GORM para validaciones
func (rt *RefreshToken) BeforeUpdate(tx *gorm.DB) error {
	if err := rt.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	return nil
}

// ValidateRefreshToken valida los datos del refresh token
func (rt *RefreshToken) ValidateRefreshToken() error {
	if rt.UserID == "" {
		return errors.New("user ID is required")
	}

	if rt.TokenHash == "" {
		return errors.New("token hash is required")
	}

	if rt.TokenID == "" {
		return errors.New("token ID is required")
	}

	if rt.ExpiresAt.IsZero() {
		return errors.New("expires at is required")
	}

	if rt.ExpiresAt.Before(time.Now()) {
		return errors.New("token cannot be created with past expiration date")
	}

	return nil
}

// IsExpired verifica si el token ha expirado
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsActive verifica si el token está activo (no revocado ni expirado)
func (rt *RefreshToken) IsActive() bool {
	return !rt.IsRevoked && !rt.IsExpired()
}

// Revoke revoca el refresh token
func (rt *RefreshToken) Revoke() {
	if !rt.IsRevoked {
		now := time.Now()
		rt.IsRevoked = true
		rt.RevokedAt = &now
	}
}

// UpdateLastUsed actualiza el timestamp de último uso
func (rt *RefreshToken) UpdateLastUsed() {
	now := time.Now()
	rt.LastUsedAt = &now
}

// GetAuditData implementa AuditableModel
func (rt *RefreshToken) GetAuditData() map[string]interface{} {
	return map[string]interface{}{
		"id":         rt.ID,
		"user_id":    rt.UserID,
		"token_id":   rt.TokenID,
		"expires_at": rt.ExpiresAt,
		"is_revoked": rt.IsRevoked,
		"ip_address": rt.IPAddress,
	}
}

// RefreshTokenRepository define las operaciones de refresh tokens
type RefreshTokenRepository interface {
	// Create almacena un nuevo refresh token
	Create(token *RefreshToken) error

	// FindByTokenHash busca un refresh token por su hash
	FindByTokenHash(tokenHash string) (*RefreshToken, error)

	// FindByTokenID busca un refresh token por su TokenID
	FindByTokenID(tokenID string) (*RefreshToken, error)

	// FindActiveByUserID busca todos los tokens activos de un usuario
	FindActiveByUserID(userID string) ([]*RefreshToken, error)

	// RevokeByTokenHash revoca un token por su hash
	RevokeByTokenHash(tokenHash string) error

	// RevokeByTokenID revoca un token por su TokenID
	RevokeByTokenID(tokenID string) error

	// RevokeAllByUserID revoca todos los tokens de un usuario
	RevokeAllByUserID(userID string) error

	// DeleteExpired elimina tokens expirados (para limpieza)
	DeleteExpired() error

	// DeleteRevokedOlderThan elimina tokens revocados más antiguos que la fecha especificada
	DeleteRevokedOlderThan(before time.Time) error

	// CountActiveByUserID cuenta tokens activos de un usuario
	CountActiveByUserID(userID string) (int64, error)
}

// RefreshTokenService define la lógica de negocio para refresh tokens
type RefreshTokenService interface {
	// StoreRefreshToken almacena un nuevo refresh token
	StoreRefreshToken(userID, tokenHash, tokenID string, expiresAt time.Time, ipAddress, userAgent string) (*RefreshToken, error)

	// ValidateRefreshToken valida un refresh token y lo retorna si es válido
	ValidateRefreshToken(tokenHash string) (*RefreshToken, error)

	// UseRefreshToken marca un refresh token como usado y lo revoca opcionalmente
	UseRefreshToken(tokenHash string, revokeAfterUse bool) (*RefreshToken, error)

	// RevokeUserTokens revoca todos los tokens de un usuario
	RevokeUserTokens(userID string) error

	// RevokeToken revoca un token específico
	RevokeToken(tokenHash string) error

	// CleanupExpiredTokens limpia tokens expirados
	CleanupExpiredTokens() error

	// GetUserSessions obtiene información de sesiones activas de un usuario
	GetUserSessions(userID string) ([]*RefreshToken, error)

	// ValidateTokenLimit verifica si el usuario puede crear más tokens
	ValidateTokenLimit(userID string, maxTokensPerUser int) error
}

// Constantes para configuración de refresh tokens
const (
	// MaxRefreshTokensPerUser límite máximo de refresh tokens activos por usuario
	MaxRefreshTokensPerUser = 5

	// RefreshTokenCleanupInterval intervalo para limpiar tokens expirados
	RefreshTokenCleanupInterval = 24 * time.Hour

	// RevokedTokenRetentionPeriod tiempo que se mantienen tokens revocados para auditoría
	RevokedTokenRetentionPeriod = 30 * 24 * time.Hour // 30 días
)

// SessionInfo información de sesión para mostrar al usuario
type SessionInfo struct {
	ID         string     `json:"id"`
	DeviceInfo string     `json:"device_info"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	IsRevoked  bool       `json:"is_revoked"`
}

// ToSessionInfo convierte RefreshToken a SessionInfo
func (rt *RefreshToken) ToSessionInfo() *SessionInfo {
	return &SessionInfo{
		ID:         rt.ID.String(),
		DeviceInfo: rt.DeviceInfo,
		IPAddress:  rt.IPAddress,
		UserAgent:  rt.UserAgent,
		CreatedAt:  rt.CreatedAt,
		LastUsedAt: rt.LastUsedAt,
		ExpiresAt:  rt.ExpiresAt,
		IsRevoked:  rt.IsRevoked,
	}
}

// Métodos de base model implementados
func (rt RefreshToken) GetID() string           { return rt.ID.String() }
func (rt RefreshToken) GetCreatedAt() time.Time { return rt.CreatedAt }
func (rt RefreshToken) GetUpdatedAt() time.Time { return rt.UpdatedAt }
