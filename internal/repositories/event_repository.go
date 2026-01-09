package repositories

import (
	"context"
	"time"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/models"

	"gorm.io/gorm"
)

// EventRepository repositorio específico para eventos
type EventRepository struct {
	*BaseRepository[models.Event]
}

// NewEventRepository crea una nueva instancia
func NewEventRepository() *EventRepository {
	base := NewBaseRepository[models.Event]()

	// Configurar filtros permitidos para eventos
	base.builder.SetAllowedFilters(map[string]string{
		"type":            "=",
		"category":        "LIKE",
		"level":           "=",
		"status":          "=",
		"is_online":       "=",
		"is_free":         "=",
		"is_featured":     "=",
		"is_public":       "=",
		"organization_id": "=",
		"venue_city":      "LIKE",
		"venue_country":   "LIKE",
		"start_date":      ">=",
		"end_date":        "<=",
	})

	// Configurar campos de ordenamiento permitidos
	base.builder.SetAllowedSorts([]string{
		"start_date", "end_date", "created_at", "updated_at",
		"title", "views_count", "current_attendees",
	})

	// Configurar campos de búsqueda
	base.builder.SetSearchFields([]string{
		"title", "description", "short_desc",
	})

	return &EventRepository{BaseRepository: base}
}

// GetPublicEvents obtiene eventos públicos para usuarios no autenticados
func (r *EventRepository) GetPublicEvents(ctx context.Context, opts common.QueryOptions) ([]*models.Event, *common.PaginationMeta, error) {
	// Agregar filtros específicos para eventos públicos
	opts.AddFilter("status", models.EventStatusPublished)
	opts.AddFilter("is_public", true)

	return r.GetAll(ctx, opts)
}

// GetByOrganization obtiene eventos de una organización específica
func (r *EventRepository) GetByOrganization(ctx context.Context, organizationID string, opts common.QueryOptions) ([]*models.Event, *common.PaginationMeta, error) {
	opts.AddFilter("organization_id", organizationID)
	return r.GetAll(ctx, opts)
}

// GetUpcoming obtiene eventos futuros
func (r *EventRepository) GetUpcoming(ctx context.Context, opts common.QueryOptions) ([]*models.Event, *common.PaginationMeta, error) {
	opts.AddFilter("start_date", time.Now())
	opts.AddFilter("status", models.EventStatusPublished)

	// Ordenar por fecha de inicio por defecto
	if opts.OrderBy == "created_at" {
		opts.OrderBy = "start_date"
		opts.OrderDir = "asc"
	}

	return r.GetAll(ctx, opts)
}

// GetPast obtiene eventos pasados
func (r *EventRepository) GetPast(ctx context.Context, opts common.QueryOptions) ([]*models.Event, *common.PaginationMeta, error) {
	// Usar una query personalizada para eventos pasados
	r.builder.AddWhere("end_date < ?", time.Now())
	opts.AddFilter("status", models.EventStatusPublished)

	return r.GetAll(ctx, opts)
}

// GetByDateRange obtiene eventos en un rango de fechas
func (r *EventRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, opts common.QueryOptions) ([]*models.Event, *common.PaginationMeta, error) {
	r.builder.AddWhere("start_date >= ? AND end_date <= ?", startDate, endDate)
	return r.GetAll(ctx, opts)
}

// GetFeatured obtiene eventos destacados
func (r *EventRepository) GetFeatured(ctx context.Context, limit int) ([]*models.Event, error) {
	opts := common.QueryOptions{
		Limit:    limit,
		OrderBy:  "created_at",
		OrderDir: "desc",
	}
	opts.AddFilter("is_featured", true)
	opts.AddFilter("status", models.EventStatusPublished)
	opts.AddFilter("is_public", true)

	events, _, err := r.GetAll(ctx, opts)
	return events, err
}

// GetWithPreloads obtiene eventos con preloads especificados
func (r *EventRepository) GetWithPreloads(ctx context.Context, id string, preloads []string) (*models.Event, error) {
	var event models.Event
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&event, "id = ?", id).Error
	if err != nil {
		return nil, common.MapGormError(err)
	}

	return &event, nil
}

// IncrementViews incrementa el contador de visualizaciones
func (r *EventRepository) IncrementViews(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Model(&models.Event{}).
		Where("id = ?", id).
		UpdateColumn("views_count", gorm.Expr("views_count + 1")).Error
	return common.MapGormError(err)
}

// UpdateStatus actualiza el estado de un evento
func (r *EventRepository) UpdateStatus(ctx context.Context, id string, status models.EventStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// Agregar timestamps específicos según el estado
	switch status {
	case models.EventStatusPublished:
		updates["published_at"] = time.Now()
	case models.EventStatusCanceled:
		updates["canceled_at"] = time.Now()
	case models.EventStatusCompleted:
		updates["completed_at"] = time.Now()
	}

	err := r.db.WithContext(ctx).Model(&models.Event{}).
		Where("id = ?", id).
		Updates(updates).Error
	return common.MapGormError(err)
}

// GetEventsByTags obtiene eventos por tags
func (r *EventRepository) GetEventsByTags(ctx context.Context, tags []string, opts common.QueryOptions) ([]*models.Event, *common.PaginationMeta, error) {
	// Usar query JSONB para tags
	r.builder.AddWhere("tags ?| array[?]", gorm.Expr("?", tags))
	opts.AddFilter("status", models.EventStatusPublished)
	opts.AddFilter("is_public", true)

	return r.GetAll(ctx, opts)
}

// GetNearEvents obtiene eventos cerca de una ubicación
func (r *EventRepository) GetNearEvents(ctx context.Context, latitude, longitude float64, radiusKm int, opts common.QueryOptions) ([]*models.Event, *common.PaginationMeta, error) {
	// Query geoespacial usando PostGIS
	r.builder.AddWhere(`
		ST_DWithin(
			ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography,
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			?
		)
	`, longitude, latitude, radiusKm*1000) // radiusKm a metros

	opts.AddFilter("status", models.EventStatusPublished)
	opts.AddFilter("is_public", true)
	opts.AddFilter("is_online", false) // Solo eventos presenciales

	return r.GetAll(ctx, opts)
}
