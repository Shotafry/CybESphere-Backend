package repositories

import (
	"context"
	"time"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/models"

	"gorm.io/gorm"
)

// OrganizationRepository repositorio específico para organizaciones
type OrganizationRepository struct {
	*BaseRepository[models.Organization]
}

// NewOrganizationRepository crea una nueva instancia
func NewOrganizationRepository() *OrganizationRepository {
	base := NewBaseRepository[models.Organization]()

	// Configurar filtros permitidos
	base.builder.SetAllowedFilters(map[string]string{
		"status":      "=",
		"is_verified": "=",
		"city":        "LIKE",
		"country":     "LIKE",
	})

	// Configurar ordenamiento
	base.builder.SetAllowedSorts([]string{
		"name", "created_at", "updated_at", "events_count", "city",
	})

	// Configurar búsqueda
	base.builder.SetSearchFields([]string{
		"name", "description",
	})

	return &OrganizationRepository{BaseRepository: base}
}

// GetActive obtiene organizaciones activas
func (r *OrganizationRepository) GetActive(ctx context.Context, opts common.QueryOptions) ([]*models.Organization, *common.PaginationMeta, error) {
	opts.AddFilter("status", models.OrgStatusActive)
	return r.GetAll(ctx, opts)
}

// GetVerified obtiene organizaciones verificadas
func (r *OrganizationRepository) GetVerified(ctx context.Context, opts common.QueryOptions) ([]*models.Organization, *common.PaginationMeta, error) {
	opts.AddFilter("is_verified", true)
	return r.GetAll(ctx, opts)
}

// GetByCity obtiene organizaciones por ciudad
func (r *OrganizationRepository) GetByCity(ctx context.Context, city string, opts common.QueryOptions) ([]*models.Organization, *common.PaginationMeta, error) {
	opts.AddFilter("city", city)
	return r.GetAll(ctx, opts)
}

// GetWithMembers obtiene organización con sus miembros
func (r *OrganizationRepository) GetWithMembers(ctx context.Context, id string) (*models.Organization, error) {
	return r.GetWithPreloads(ctx, id, []string{"Users"})
}

// GetWithEvents obtiene organización con sus eventos
func (r *OrganizationRepository) GetWithEvents(ctx context.Context, id string) (*models.Organization, error) {
	return r.GetWithPreloads(ctx, id, []string{"Events"})
}

// GetWithPreloads obtiene organización con preloads específicos
func (r *OrganizationRepository) GetWithPreloads(ctx context.Context, id string, preloads []string) (*models.Organization, error) {
	var org models.Organization
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&org, "id = ?", id).Error
	if err != nil {
		return nil, common.MapGormError(err)
	}

	return &org, nil
}

// IncrementEventsCount incrementa contador de eventos
func (r *OrganizationRepository) IncrementEventsCount(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Model(&models.Organization{}).
		Where("id = ?", id).
		UpdateColumn("events_count", gorm.Expr("events_count + 1")).Error
	return common.MapGormError(err)
}

// DecrementEventsCount decrementa contador de eventos
func (r *OrganizationRepository) DecrementEventsCount(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Model(&models.Organization{}).
		Where("id = ?", id).
		UpdateColumn("events_count", gorm.Expr("GREATEST(events_count - 1, 0)")).Error
	return common.MapGormError(err)
}

// UpdateStatus actualiza el estado de una organización
func (r *OrganizationRepository) UpdateStatus(ctx context.Context, id string, status models.OrganizationStatus) error {
	err := r.db.WithContext(ctx).Model(&models.Organization{}).
		Where("id = ?", id).
		UpdateColumn("status", status).Error
	return common.MapGormError(err)
}

// Verify marca una organización como verificada
func (r *OrganizationRepository) Verify(ctx context.Context, id string, verifierID string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"is_verified": true,
		"verified_by": verifierID,
		"verified_at": now.Format(time.RFC3339),
		"status":      models.OrgStatusActive,
	}

	err := r.db.WithContext(ctx).Model(&models.Organization{}).
		Where("id = ?", id).
		Updates(updates).Error
	return common.MapGormError(err)
}
