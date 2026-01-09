// internal/handlers/event_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/mappers"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/services"
)

type EventHandler struct {
	*BaseHandler[models.Event, dto.CreateEventRequest, dto.UpdateEventRequest, dto.EventResponse]
	eventService services.EventService
	mapper       *mappers.UnifiedMapper
}

func NewEventHandler(
	eventService services.EventService,
	mapper *mappers.UnifiedMapper,
) *EventHandler {
	// Crear handler base
	base := NewBaseHandler[models.Event, dto.CreateEventRequest, dto.UpdateEventRequest, dto.EventResponse](
		eventService,
		mapper,
		"event",
	)

	// Configurar límites de paginación específicos
	base.SetLimits(20, 100)

	return &EventHandler{
		BaseHandler:  base,
		eventService: eventService,
		mapper:       mapper,
	}
}

// Los métodos CRUD básicos ya están implementados en BaseHandler
// Solo agregamos métodos específicos de eventos

// PublishEvent método específico de eventos
func (h *EventHandler) PublishEvent(c *gin.Context) {
	eventID := c.Param("eventId")
	userCtx := extractUserContext(c)

	event, err := h.eventService.PublishEvent(c.Request.Context(), eventID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	response := h.mapper.EventToResponse(event, userCtx)
	common.SuccessResponse(c, http.StatusOK, "Evento publicado exitosamente", response)
}

// GetPublicEvent obtiene un evento público por ID o Slug
func (h *EventHandler) GetPublicEvent(c *gin.Context) {
	identifier := c.Param("id")
	userCtx := extractUserContext(c)

	var event *models.Event
	var err error

	// Verificar si es un UUID válido
	if _, uuidErr := uuid.Parse(identifier); uuidErr == nil {
		// Es un UUID, buscar por ID
		event, err = h.eventService.Get(c.Request.Context(), identifier, userCtx)
	} else {
		// No es UUID, asumir que es Slug
		event, err = h.eventService.GetEventBySlug(c.Request.Context(), identifier)
	}

	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Verificar qu el evento sea público si el usuario no tiene permisos especiales
	// NOTA: Esto ya debería estar filtrado por el repositorio/servicio, pero doble check no hace daño
	if !event.IsPublic && (userCtx == nil || !userCtx.IsAdmin()) {
		common.ErrorResponse(c, common.NewBusinessError("event_not_found", "Evento no encontrado"))
		return
	}

	response := h.mapper.EventToResponse(event, userCtx)
	common.SuccessResponse(c, http.StatusOK, "Detalle del evento", response)
}

// CancelEvent método específico
func (h *EventHandler) CancelEvent(c *gin.Context) {
	eventID := c.Param("eventId")
	userCtx := extractUserContext(c)

	var req dto.CancelEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, common.NewValidationError("request", err.Error()))
		return
	}

	event, err := h.eventService.CancelEvent(c.Request.Context(), eventID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	response := h.mapper.EventToResponse(event, userCtx)
	common.SuccessResponse(c, http.StatusOK, "Evento cancelado", response)
}

// GetUpcomingEvents eventos futuros
func (h *EventHandler) GetUpcomingEvents(c *gin.Context) {
	opts := extractQueryOptions(c)
	userCtx := extractUserContext(c)

	events, pagination, err := h.eventService.GetUpcomingEvents(
		c.Request.Context(),
		*opts,
		userCtx,
	)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	response := h.mapper.EventsToListResponse(events, pagination, userCtx)
	common.SuccessWithPagination(c, "Próximos eventos", response.Events, pagination)
}

// GetEventsByOrganization eventos de una organización
func (h *EventHandler) GetEventsByOrganization(c *gin.Context) {
	orgID := c.Param("orgId")
	opts := extractQueryOptions(c)
	userCtx := extractUserContext(c)

	events, pagination, err := h.eventService.GetEventsByOrganization(
		c.Request.Context(),
		orgID,
		*opts,
		userCtx,
	)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	response := h.mapper.EventsToListResponse(events, pagination, userCtx)
	common.SuccessWithPagination(c, "Eventos de la organización", response.Events, pagination)
}

// AddToFavorites agregar a favoritos
func (h *EventHandler) AddToFavorites(c *gin.Context) {
	eventID := c.Param("eventId")
	userCtx := extractUserContext(c)

	err := h.eventService.AddToFavorites(c.Request.Context(), eventID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Agregado a favoritos", nil)
}

// RemoveFromFavorites remover de favoritos
func (h *EventHandler) RemoveFromFavorites(c *gin.Context) {
	eventID := c.Param("eventId")
	userCtx := extractUserContext(c)

	err := h.eventService.RemoveFromFavorites(c.Request.Context(), eventID, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Removido de favoritos", nil)
}

// GetFeaturedEvents eventos destacados
func (h *EventHandler) GetFeaturedEvents(c *gin.Context) {
	events, err := h.eventService.GetFeaturedEvents(c.Request.Context(), 10)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Convertir a respuestas resumidas
	summaries := make([]dto.EventSummaryResponse, 0, len(events))
	for _, event := range events {
		summaries = append(summaries, h.mapper.EventToSummaryResponse(event))
	}

	common.SuccessResponse(c, http.StatusOK, "Eventos destacados", summaries)
}
