package mappers

import (
	"strings"
	"time"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/models"
)

// EventMapperImpl implementación del mapper de eventos
type EventMapperImpl struct{}

// NewEventMapper crea nueva instancia del mapper
func NewEventMapper() EventMapperImpl {
	return EventMapperImpl{}
}

// CreateEventRequestToModel convierte CreateEventRequest a modelo Event
func (m EventMapperImpl) CreateEventRequestToModel(req *dto.CreateEventRequest, userCtx *common.UserContext) (*models.Event, error) {
	// Determinar organización según el rol del usuario
	organizationID := req.OrganizationID
	if !userCtx.IsAdmin() && userCtx.OrganizationID != nil {
		// Organizadores usan su propia organización
		organizationID = *userCtx.OrganizationID
	}

	if organizationID == "" {
		return nil, common.NewBusinessError("missing_organization", "ID de organización requerido")
	}

	event := &models.Event{
		// Información básica
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		ShortDesc:   strings.TrimSpace(req.ShortDesc),
		Type:        models.EventType(req.Type),
		Category:    strings.TrimSpace(req.Category),
		Level:       strings.TrimSpace(req.Level),

		// Fechas y horarios
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Timezone:  req.Timezone,

		// Ubicación
		IsOnline:     req.IsOnline,
		VenueAddress: strings.TrimSpace(req.VenueAddress),
		VenueName:    strings.TrimSpace(req.VenueName),
		VenueCity:    strings.TrimSpace(req.VenueCity),
		VenueCountry: strings.TrimSpace(req.VenueCountry),
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		OnlineURL:    strings.TrimSpace(req.OnlineURL),
		StreamingURL: strings.TrimSpace(req.StreamingURL),

		// Capacidad y registro
		MaxAttendees:    req.MaxAttendees,
		IsFree:          req.IsFree,
		Price:           req.Price,
		Currency:        req.Currency,
		RegistrationURL: strings.TrimSpace(req.RegistrationURL),

		// Contenido
		ImageURL:     strings.TrimSpace(req.ImageURL),
		BannerURL:    strings.TrimSpace(req.BannerURL),
		Requirements: strings.TrimSpace(req.Requirements),
		Agenda:       strings.TrimSpace(req.Agenda),

		// Fechas importantes
		RegistrationStartDate: req.RegistrationStartDate,
		RegistrationEndDate:   req.RegistrationEndDate,

		// Información de contacto
		ContactEmail: strings.ToLower(strings.TrimSpace(req.ContactEmail)),
		ContactPhone: strings.TrimSpace(req.ContactPhone),

		// Metadatos SEO
		MetaTitle:       strings.TrimSpace(req.MetaTitle),
		MetaDescription: strings.TrimSpace(req.MetaDescription),

		// Organización
		OrganizationID: organizationID,

		// Estado inicial
		Status:     models.EventStatusDraft,
		IsPublic:   true,
		IsFeatured: false, // Solo admin puede destacar en creación
	}

	// Aplicar valores por defecto
	if event.Timezone == "" {
		event.Timezone = "Europe/Madrid"
	}
	if event.Currency == "" {
		event.Currency = "EUR"
	}

	// Solo admin puede establecer ciertos campos en creación
	if userCtx.IsAdmin() {
		if req.Status != "" {
			event.Status = models.EventStatus(req.Status)
		}
		if req.IsPublic != nil {
			event.IsPublic = *req.IsPublic
		}
		if req.IsFeatured != nil {
			event.IsFeatured = *req.IsFeatured
		}
	}

	// Configurar tags si se proporcionaron
	if len(req.Tags) > 0 {
		if err := event.SetTags(req.Tags); err != nil {
			return nil, common.NewBusinessError("invalid_tags", "Error al procesar tags")
		}
	}

	return event, nil
}

// UpdateEventRequestToModel aplica cambios de UpdateEventRequest a evento existente
func (m EventMapperImpl) UpdateEventRequestToModel(event *models.Event, req *dto.UpdateEventRequest, userCtx *common.UserContext) error {
	// Aplicar cambios básicos permitidos para organizadores y admin
	if req.Title != nil {
		event.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		event.Description = strings.TrimSpace(*req.Description)
	}
	if req.ShortDesc != nil {
		event.ShortDesc = strings.TrimSpace(*req.ShortDesc)
	}
	if req.Category != nil {
		event.Category = strings.TrimSpace(*req.Category)
	}
	if req.Level != nil {
		event.Level = strings.TrimSpace(*req.Level)
	}

	// Fechas
	if req.StartDate != nil {
		event.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		event.EndDate = *req.EndDate
	}
	if req.Timezone != nil {
		event.Timezone = strings.TrimSpace(*req.Timezone)
	}

	// Ubicación
	if req.IsOnline != nil {
		event.IsOnline = *req.IsOnline
	}
	if req.VenueAddress != nil {
		event.VenueAddress = strings.TrimSpace(*req.VenueAddress)
	}
	if req.VenueName != nil {
		event.VenueName = strings.TrimSpace(*req.VenueName)
	}
	if req.VenueCity != nil {
		event.VenueCity = strings.TrimSpace(*req.VenueCity)
	}
	if req.VenueCountry != nil {
		event.VenueCountry = strings.TrimSpace(*req.VenueCountry)
	}
	if req.Latitude != nil {
		event.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		event.Longitude = req.Longitude
	}
	if req.OnlineURL != nil {
		event.OnlineURL = strings.TrimSpace(*req.OnlineURL)
	}
	if req.StreamingURL != nil {
		event.StreamingURL = strings.TrimSpace(*req.StreamingURL)
	}

	// Capacidad y registro
	if req.MaxAttendees != nil {
		event.MaxAttendees = req.MaxAttendees
	}
	if req.IsFree != nil {
		event.IsFree = *req.IsFree
	}
	if req.Price != nil {
		event.Price = req.Price
	}
	if req.Currency != nil {
		event.Currency = strings.TrimSpace(*req.Currency)
	}
	if req.RegistrationURL != nil {
		event.RegistrationURL = strings.TrimSpace(*req.RegistrationURL)
	}

	// Contenido
	if req.ImageURL != nil {
		event.ImageURL = strings.TrimSpace(*req.ImageURL)
	}
	if req.BannerURL != nil {
		event.BannerURL = strings.TrimSpace(*req.BannerURL)
	}
	if req.Requirements != nil {
		event.Requirements = strings.TrimSpace(*req.Requirements)
	}
	if req.Agenda != nil {
		event.Agenda = strings.TrimSpace(*req.Agenda)
	}

	// Fechas importantes
	if req.RegistrationStartDate != nil {
		event.RegistrationStartDate = req.RegistrationStartDate
	}
	if req.RegistrationEndDate != nil {
		event.RegistrationEndDate = req.RegistrationEndDate
	}

	// Información de contacto
	if req.ContactEmail != nil {
		event.ContactEmail = strings.ToLower(strings.TrimSpace(*req.ContactEmail))
	}
	if req.ContactPhone != nil {
		event.ContactPhone = strings.TrimSpace(*req.ContactPhone)
	}

	// Metadatos SEO
	if req.MetaTitle != nil {
		event.MetaTitle = strings.TrimSpace(*req.MetaTitle)
	}
	if req.MetaDescription != nil {
		event.MetaDescription = strings.TrimSpace(*req.MetaDescription)
	}

	// Tags
	if len(req.Tags) > 0 {
		if err := event.SetTags(req.Tags); err != nil {
			return common.NewBusinessError("invalid_tags", "Error al procesar tags")
		}
	}

	// Campos que solo admin puede cambiar
	if userCtx.IsAdmin() {
		if req.Status != nil {
			event.Status = models.EventStatus(*req.Status)
		}
		if req.IsPublic != nil {
			event.IsPublic = *req.IsPublic
		}
		if req.IsFeatured != nil {
			event.IsFeatured = *req.IsFeatured
		}
	}

	return nil
}

// EventToResponse convierte Event model a EventResponse
func (m EventMapperImpl) EventToResponse(event *models.Event, userCtx *common.UserContext) dto.EventResponse {
	response := dto.EventResponse{
		// Identificación
		ID:   event.ID.String(),
		Slug: event.Slug,

		// Información básica
		Title:       event.Title,
		Description: event.Description,
		ShortDesc:   event.ShortDesc,
		Type:        string(event.Type),
		Category:    event.Category,
		Level:       event.Level,

		// Fechas y horarios
		StartDate: event.StartDate,
		EndDate:   event.EndDate,
		Timezone:  event.Timezone,
		Duration:  event.Duration,

		// Ubicación
		IsOnline:     event.IsOnline,
		VenueName:    event.VenueName,
		VenueAddress: event.VenueAddress,
		VenueCity:    event.VenueCity,
		VenueCountry: event.VenueCountry,
		Latitude:     event.Latitude,
		Longitude:    event.Longitude,
		OnlineURL:    event.OnlineURL,
		StreamingURL: event.StreamingURL,

		// Capacidad y registro
		MaxAttendees:     event.MaxAttendees,
		CurrentAttendees: event.CurrentAttendees,
		AvailableSpots:   event.GetAvailableSpots(),
		IsFree:           event.IsFree,
		Price:            event.Price,
		Currency:         event.Currency,
		RegistrationURL:  event.RegistrationURL,

		// Estado y visibilidad
		Status:     string(event.Status),
		IsPublic:   event.IsPublic,
		IsFeatured: event.IsFeatured,
		ViewsCount: event.ViewsCount,

		// Contenido
		ImageURL:  event.ImageURL,
		BannerURL: event.BannerURL,
		Tags:      event.GetTags(),

		// Información de registro
		RegistrationOpen:      event.IsRegistrationOpen(),
		RegistrationStartDate: event.RegistrationStartDate,
		RegistrationEndDate:   event.RegistrationEndDate,

		// Metadatos SEO (filtrados según permisos)
		MetaTitle:       m.filterMetadata(event.MetaTitle, userCtx),
		MetaDescription: m.filterMetadata(event.MetaDescription, userCtx),

		// Organización
		Organization: m.mapOrganizationSummary(event.Organization),

		// Estados computados
		IsUpcoming: event.IsUpcoming(),
		IsPast:     event.IsPast(),
		IsOngoing:  m.isOngoing(event),

		// Información del usuario (si está autenticado)
		IsFavorite:   false, // TODO: implementar lógica de favoritos
		IsRegistered: false, // TODO: implementar lógica de registro
		CanEdit:      m.canEdit(event, userCtx),
		CanManage:    m.canManage(event, userCtx),

		// Timestamps
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
		PublishedAt: event.PublishedAt,
		CanceledAt:  event.CanceledAt,
		CompletedAt: event.CompletedAt,
	}

	return response
}

// EventToDetailResponse convierte Event a respuesta detallada
func (m EventMapperImpl) EventToDetailResponse(event *models.Event, userCtx *common.UserContext) dto.EventDetailResponse {
	baseResponse := m.EventToResponse(event, userCtx)

	detailResponse := dto.EventDetailResponse{
		EventResponse: baseResponse,
		// Información adicional solo para usuarios autorizados
		Requirements: m.filterSensitiveField(event.Requirements, userCtx, event),
		Agenda:       m.filterSensitiveField(event.Agenda, userCtx, event),
		ContactEmail: m.filterContactInfo(event.ContactEmail, userCtx, event),
		ContactPhone: m.filterContactInfo(event.ContactPhone, userCtx, event),
	}

	// Estadísticas para organizadores/admin
	if m.canViewStatistics(event, userCtx) {
		detailResponse.Statistics = m.buildEventStatistics(event)
	}

	return detailResponse
}

// EventToSummaryResponse convierte a respuesta resumida
func (m EventMapperImpl) EventToSummaryResponse(event *models.Event) dto.EventSummaryResponse {
	return dto.EventSummaryResponse{
		ID:               event.ID.String(),
		Slug:             event.Slug,
		Title:            event.Title,
		ShortDesc:        event.ShortDesc,
		Type:             string(event.Type),
		StartDate:        event.StartDate,
		EndDate:          event.EndDate,
		IsOnline:         event.IsOnline,
		VenueCity:        event.VenueCity,
		IsFree:           event.IsFree,
		Price:            event.Price,
		ImageURL:         event.ImageURL,
		CurrentAttendees: event.CurrentAttendees,
		MaxAttendees:     event.MaxAttendees,
		IsFeatured:       event.IsFeatured,
		Organization:     m.getOrganizationName(event.Organization),
		Tags:             event.GetTags(),
	}
}

// EventsToListResponse convierte lista de eventos a respuesta con paginación
func (m EventMapperImpl) EventsToListResponse(events []*models.Event, pagination *common.PaginationMeta, userCtx *common.UserContext) dto.EventListResponse {
	eventResponses := make([]dto.EventResponse, 0, len(events))

	for _, event := range events {
		eventResponses = append(eventResponses, m.EventToResponse(event, userCtx))
	}

	return dto.EventListResponse{
		Events:     eventResponses,
		Pagination: *pagination,
		Filters:    dto.AppliedFilters{}, // TODO: implementar filtros aplicados
	}
}

// =============================================================================
// MÉTODOS HELPER PRIVADOS
// =============================================================================

// mapOrganizationSummary convierte Organization a OrganizationSummaryResponse
func (m EventMapperImpl) mapOrganizationSummary(org *models.Organization) dto.OrganizationSummaryResponse {
	if org == nil {
		return dto.OrganizationSummaryResponse{}
	}

	return dto.OrganizationSummaryResponse{
		ID:          org.ID.String(),
		Slug:        org.Slug,
		Name:        org.Name,
		LogoURL:     org.LogoURL,
		IsVerified:  org.IsVerified,
		EventsCount: org.EventsCount,
		City:        org.City,
		Country:     org.Country,
	}
}

// isOngoing determina si un evento está en curso
func (m EventMapperImpl) isOngoing(event *models.Event) bool {
	now := time.Now()
	return now.After(event.StartDate) && now.Before(event.EndDate) && event.IsActive()
}

// canEdit verifica si el usuario puede editar el evento
func (m EventMapperImpl) canEdit(event *models.Event, userCtx *common.UserContext) bool {
	if userCtx == nil {
		return false
	}

	if userCtx.IsAdmin() {
		return true
	}

	return userCtx.IsOrganizer() && userCtx.CanManageOrganization(event.OrganizationID)
}

// canManage verifica si el usuario puede gestionar el evento (permisos completos)
func (m EventMapperImpl) canManage(event *models.Event, userCtx *common.UserContext) bool {
	return m.canEdit(event, userCtx) // Por ahora son iguales, puede diferir en el futuro
}

// canViewStatistics verifica si puede ver estadísticas del evento
func (m EventMapperImpl) canViewStatistics(event *models.Event, userCtx *common.UserContext) bool {
	return m.canManage(event, userCtx)
}

// filterMetadata filtra metadatos SEO según permisos
func (m EventMapperImpl) filterMetadata(metadata string, userCtx *common.UserContext) string {
	// Los metadatos SEO son públicos si el evento es público
	return metadata
}

// filterSensitiveField filtra campos sensibles según permisos
func (m EventMapperImpl) filterSensitiveField(field string, userCtx *common.UserContext, event *models.Event) string {
	if userCtx == nil {
		return ""
	}

	// Solo usuarios autenticados pueden ver información detallada
	if event.IsActive() && event.IsPublic {
		return field
	}

	// O usuarios que pueden gestionar el evento
	if m.canManage(event, userCtx) {
		return field
	}

	return ""
}

// filterContactInfo filtra información de contacto
func (m EventMapperImpl) filterContactInfo(contact string, userCtx *common.UserContext, event *models.Event) string {
	if userCtx == nil {
		return ""
	}

	// Solo usuarios registrados o gestores pueden ver info de contacto
	if m.canManage(event, userCtx) {
		return contact
	}

	// TODO: implementar lógica para usuarios registrados al evento
	return ""
}

// buildEventStatistics construye estadísticas del evento
func (m EventMapperImpl) buildEventStatistics(event *models.Event) *dto.EventStatistics {
	// Implementación básica - en producción consultarías métricas reales
	return &dto.EventStatistics{
		TotalViews:         event.ViewsCount,
		UniqueViews:        event.ViewsCount, // Placeholder
		TotalRegistrations: event.CurrentAttendees,
		ConfirmedAttendees: event.CurrentAttendees,
		PendingAttendees:   0,
		CanceledAttendees:  0,
		ConversionRate:     0.0, // Calcular views to registration
		OccupancyRate:      m.calculateOccupancyRate(event),
	}
}

// calculateOccupancyRate calcula tasa de ocupación
func (m EventMapperImpl) calculateOccupancyRate(event *models.Event) float64 {
	if event.MaxAttendees == nil || *event.MaxAttendees == 0 {
		return 0.0
	}

	return float64(event.CurrentAttendees) / float64(*event.MaxAttendees) * 100.0
}

// getOrganizationName obtiene nombre de la organización de forma segura
func (m EventMapperImpl) getOrganizationName(org *models.Organization) string {
	if org == nil {
		return ""
	}
	return org.Name
}
