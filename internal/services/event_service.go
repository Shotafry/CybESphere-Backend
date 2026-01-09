// internal/services/event_service.go
package services

import (
	"context"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/repositories"
	"cybesphere-backend/pkg/logger"
)

// EventServiceImpl implementación concreta del servicio de eventos
type EventServiceImpl struct {
	*BaseService[models.Event, dto.CreateEventRequest, dto.UpdateEventRequest]
	eventRepo *repositories.EventRepository
	orgRepo   *repositories.OrganizationRepository
	userRepo  *repositories.UserRepository
	auth      AuthorizationService
}

// Verificación en tiempo de compilación de que EventServiceImpl implementa EventService
var _ EventService = (*EventServiceImpl)(nil)

// NewEventService crea una nueva instancia del servicio de eventos
func NewEventService(
	eventRepo *repositories.EventRepository,
	orgRepo *repositories.OrganizationRepository,
	userRepo *repositories.UserRepository,
	mapper ResponseMapper,
	auth AuthorizationService,
) EventService {
	base := NewBaseService[models.Event, dto.CreateEventRequest, dto.UpdateEventRequest](
		eventRepo, mapper, auth,
	)
	return &EventServiceImpl{
		BaseService: base,
		eventRepo:   eventRepo,
		orgRepo:     orgRepo,
		userRepo:    userRepo,
		auth:        auth,
	}
}

// CreateEvent crea un nuevo evento con validaciones de negocio (wrapper para compatibilidad)
func (s *EventServiceImpl) CreateEvent(ctx context.Context, req dto.CreateEventRequest, userCtx *common.UserContext) (*models.Event, error) {
	// Validaciones de negocio específicas
	if err := s.validateEventCreation(ctx, &req, userCtx); err != nil {
		return nil, err
	}

	// Usar el método base para crear
	event, err := s.Create(ctx, req, userCtx)
	if err != nil {
		return nil, err
	}

	// Post-procesamiento: incrementar contador de eventos de la organización
	go func() {
		if err := s.orgRepo.IncrementEventsCount(context.Background(), event.OrganizationID); err != nil {
			logger.Error("Error incrementando contador de eventos de la organización: ", err)
		}
	}()

	return event, nil
}

// GetPublicEvents obtiene eventos públicos (para usuarios no autenticados)
func (s *EventServiceImpl) GetPublicEvents(ctx context.Context, opts common.QueryOptions) ([]*models.Event, *common.PaginationMeta, error) {
	return s.eventRepo.GetPublicEvents(ctx, opts)
}

// GetEventBySlug obtiene un evento por su slug
func (s *EventServiceImpl) GetEventBySlug(ctx context.Context, slug string) (*models.Event, error) {
	return s.eventRepo.GetBySlug(ctx, slug)
}

// GetUpcomingEvents obtiene eventos futuros
func (s *EventServiceImpl) GetUpcomingEvents(ctx context.Context, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.Event, *common.PaginationMeta, error) {
	// Aplicar filtros de seguridad
	s.auth.ApplySecurityFilters(&opts, userCtx, "event")

	return s.eventRepo.GetUpcoming(ctx, opts)
}

// GetEventsByOrganization obtiene eventos de una organización
func (s *EventServiceImpl) GetEventsByOrganization(ctx context.Context, organizationID string, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.Event, *common.PaginationMeta, error) {
	// Verificar que la organización existe
	_, err := s.orgRepo.GetByID(ctx, organizationID)
	if err != nil {
		return nil, nil, err
	}

	// Aplicar filtros de seguridad
	s.auth.ApplySecurityFilters(&opts, userCtx, "event")

	return s.eventRepo.GetByOrganization(ctx, organizationID, opts)
}

// PublishEvent publica un evento
func (s *EventServiceImpl) PublishEvent(ctx context.Context, id string, userCtx *common.UserContext) (*models.Event, error) {
	// Verificar permisos
	if err := s.auth.CheckUpdatePermission(userCtx, "event", id); err != nil {
		return nil, err
	}

	// Obtener evento
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validar que se puede publicar
	if err := event.Publish(); err != nil {
		return nil, common.NewBusinessError("publish_failed", err.Error())
	}

	// Actualizar en base de datos
	if err := s.eventRepo.UpdateStatus(ctx, id, models.EventStatusPublished); err != nil {
		return nil, err
	}

	// Retornar evento actualizado
	return s.eventRepo.GetByID(ctx, id)
}

// CancelEvent cancela un evento
func (s *EventServiceImpl) CancelEvent(ctx context.Context, id string, userCtx *common.UserContext) (*models.Event, error) {
	// Verificar permisos
	if err := s.auth.CheckUpdatePermission(userCtx, "event", id); err != nil {
		return nil, err
	}

	// Obtener evento
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validar que se puede cancelar
	if err := event.Cancel(); err != nil {
		return nil, common.NewBusinessError("cancel_failed", err.Error())
	}

	// Actualizar en base de datos
	if err := s.eventRepo.UpdateStatus(ctx, id, models.EventStatusCanceled); err != nil {
		return nil, err
	}

	return s.eventRepo.GetByID(ctx, id)
}

// IncrementViews incrementa las visualizaciones de un evento
func (s *EventServiceImpl) IncrementViews(ctx context.Context, id string) error {
	return s.eventRepo.IncrementViews(ctx, id)
}

// GetFeaturedEvents obtiene eventos destacados
func (s *EventServiceImpl) GetFeaturedEvents(ctx context.Context, limit int) ([]*models.Event, error) {
	return s.eventRepo.GetFeatured(ctx, limit)
}

// GetEventsByTags busca eventos por tags
func (s *EventServiceImpl) GetEventsByTags(ctx context.Context, tags []string, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.Event, *common.PaginationMeta, error) {
	s.auth.ApplySecurityFilters(&opts, userCtx, "event")
	return s.eventRepo.GetEventsByTags(ctx, tags, opts)
}

// GetNearEvents busca eventos cerca de una ubicación
func (s *EventServiceImpl) GetNearEvents(ctx context.Context, latitude, longitude float64, radiusKm int, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.Event, *common.PaginationMeta, error) {
	s.auth.ApplySecurityFilters(&opts, userCtx, "event")
	return s.eventRepo.GetNearEvents(ctx, latitude, longitude, radiusKm, opts)
}

// AddToFavorites agrega evento a favoritos del usuario
func (s *EventServiceImpl) AddToFavorites(ctx context.Context, eventID string, userCtx *common.UserContext) error {
	// Verificar que el evento existe y es público
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if !event.IsActive() || !event.IsPublic {
		return common.NewBusinessError("event_not_available", "El evento no está disponible")
	}

	// Agregar a favoritos (implementación simplificada)
	// En la implementación real, tendrías una tabla many-to-many
	return nil
}

// RemoveFromFavorites remueve evento de favoritos
func (s *EventServiceImpl) RemoveFromFavorites(ctx context.Context, eventID string, userCtx *common.UserContext) error {
	// Implementación de remoción de favoritos
	return nil
}

// validateEventCreation valida reglas de negocio para creación de eventos
func (s *EventServiceImpl) validateEventCreation(ctx context.Context, req *dto.CreateEventRequest, userCtx *common.UserContext) error {
	// Determinar organización
	var organizationID string
	if userCtx.IsAdmin() && req.OrganizationID != "" {
		organizationID = req.OrganizationID
	} else if userCtx.IsOrganizer() && userCtx.OrganizationID != nil {
		organizationID = *userCtx.OrganizationID
	} else {
		return common.NewBusinessError("no_organization", "Se requiere una organización para crear eventos")
	}

	// Verificar que la organización puede crear eventos
	org, err := s.orgRepo.GetByID(ctx, organizationID)
	if err != nil {
		return err
	}

	if !org.CanCreateEvent() {
		return common.NewBusinessError("organization_cannot_create_events",
			"La organización no puede crear eventos en este momento")
	}

	// Establecer la organización en la request
	req.OrganizationID = organizationID

	// Validaciones adicionales de fechas
	if req.EndDate.Before(req.StartDate) {
		return common.NewBusinessError("invalid_dates", "La fecha de fin debe ser posterior a la de inicio")
	}

	// Validar ubicación según tipo
	if req.IsOnline && req.OnlineURL == "" {
		return common.NewBusinessError("missing_online_url", "URL requerida para eventos online")
	}

	if !req.IsOnline && req.VenueAddress == "" {
		return common.NewBusinessError("missing_venue", "Dirección requerida para eventos presenciales")
	}

	return nil
}
