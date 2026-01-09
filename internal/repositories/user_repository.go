package repositories

import (
	"context"
	"time"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/models"
)

// UserRepository repositorio específico para usuarios
type UserRepository struct {
	*BaseRepository[models.User]
}

// NewUserRepository crea una nueva instancia
func NewUserRepository() *UserRepository {
	base := NewBaseRepository[models.User]()

	// Configurar filtros permitidos
	base.builder.SetAllowedFilters(map[string]string{
		"role":            "=",
		"is_active":       "=",
		"is_verified":     "=",
		"city":            "LIKE",
		"country":         "LIKE",
		"email":           "LIKE",
		"organization_id": "=",
	})

	// Configurar ordenamiento
	base.builder.SetAllowedSorts([]string{
		"created_at", "updated_at", "email", "first_name", "last_name", "last_login_at",
	})

	// Configurar búsqueda
	base.builder.SetSearchFields([]string{
		"first_name", "last_name", "email", "company",
	})

	return &UserRepository{BaseRepository: base}
}

// GetByEmail obtiene usuario por email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, common.MapGormError(err)
	}
	return &user, nil
}

// GetActive obtiene usuarios activos
func (r *UserRepository) GetActive(ctx context.Context, opts common.QueryOptions) ([]*models.User, *common.PaginationMeta, error) {
	opts.AddFilter("is_active", true)
	return r.GetAll(ctx, opts)
}

// GetByRole obtiene usuarios por rol
func (r *UserRepository) GetByRole(ctx context.Context, role models.UserRole, opts common.QueryOptions) ([]*models.User, *common.PaginationMeta, error) {
	opts.AddFilter("role", role)
	return r.GetAll(ctx, opts)
}

// GetByOrganization obtiene usuarios de una organización
func (r *UserRepository) GetByOrganization(ctx context.Context, organizationID string, opts common.QueryOptions) ([]*models.User, *common.PaginationMeta, error) {
	opts.AddFilter("organization_id", organizationID)
	return r.GetAll(ctx, opts)
}

// GetWithOrganization obtiene usuario con su organización
func (r *UserRepository) GetWithOrganization(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Organization").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, common.MapGormError(err)
	}
	return &user, nil
}

// UpdateLastLogin actualiza timestamp de último login
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn("last_login_at", time.Now()).Error
	return common.MapGormError(err)
}

// ActivateUser activa un usuario
func (r *UserRepository) ActivateUser(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn("is_active", true).Error
	return common.MapGormError(err)
}

// DeactivateUser desactiva un usuario
func (r *UserRepository) DeactivateUser(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn("is_active", false).Error
	return common.MapGormError(err)
}

// UpdateRole actualiza el rol de un usuario
func (r *UserRepository) UpdateRole(ctx context.Context, id string, role models.UserRole) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		UpdateColumn("role", role).Error
	return common.MapGormError(err)
}
