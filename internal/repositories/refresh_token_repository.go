package repositories

import (
	"context"
	"time"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/models"
)

// RefreshTokenRepository repositorio para refresh tokens
type RefreshTokenRepository struct {
	*BaseRepository[models.RefreshToken]
}

// NewRefreshTokenRepository crea una nueva instancia
func NewRefreshTokenRepository() *RefreshTokenRepository {
	base := NewBaseRepository[models.RefreshToken]()

	// Configurar filtros específicos
	base.builder.SetAllowedFilters(map[string]string{
		"user_id":    "=",
		"is_revoked": "=",
		"ip_address": "=",
	})

	base.builder.SetAllowedSorts([]string{
		"created_at", "updated_at", "expires_at", "last_used_at",
	})

	return &RefreshTokenRepository{BaseRepository: base}
}

// GetByTokenHash obtiene token por hash
func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, common.MapGormError(err)
	}
	return &token, nil
}

// GetActiveByUserID obtiene tokens activos de un usuario
func (r *RefreshTokenRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*models.RefreshToken, error) {
	var tokens []*models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&tokens).Error
	if err != nil {
		return nil, common.MapGormError(err)
	}
	return tokens, nil
}

// RevokeByTokenHash revoca un token por hash
func (r *RefreshTokenRepository) RevokeByTokenHash(ctx context.Context, tokenHash string) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
	return common.MapGormError(err)
}

// RevokeAllByUserID revoca todos los tokens de un usuario
func (r *RefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID string) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ? AND is_revoked = false", userID).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
	return common.MapGormError(err)
}

// DeleteExpired elimina tokens expirados
func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	err := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.RefreshToken{}).Error
	return common.MapGormError(err)
}

// CleanupOldTokens limpia tokens antiguos revocados
func (r *RefreshTokenRepository) CleanupOldTokens(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	err := r.db.WithContext(ctx).
		Where("is_revoked = true AND revoked_at < ?", cutoff).
		Delete(&models.RefreshToken{}).Error
	return common.MapGormError(err)
}

// CountActiveByUserID cuenta tokens activos de un usuario
func (r *RefreshTokenRepository) CountActiveByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ? AND is_revoked = false AND expires_at > ?", userID, time.Now()).
		Count(&count).Error
	return count, common.MapGormError(err)
}
