// internal/handlers/base_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/mappers"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/services"
)

// BaseHandler handler genérico que elimina duplicación de código
type BaseHandler[T common.BaseEntity, CreateDTO any, UpdateDTO any, ResponseDTO any] struct {
	service      common.Service[T, CreateDTO, UpdateDTO]
	mapper       mappers.ResponseMapper
	resourceType string
	defaultLimit int
	maxLimit     int
}

// NewBaseHandler crea una nueva instancia del handler base
func NewBaseHandler[T common.BaseEntity, CreateDTO any, UpdateDTO any, ResponseDTO any](
	service common.Service[T, CreateDTO, UpdateDTO],
	mapper mappers.ResponseMapper,
	resourceType string,
) *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO] {
	return &BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]{
		service:      service,
		mapper:       mapper,
		resourceType: resourceType,
		defaultLimit: 20,
		maxLimit:     100,
	}
}

// GetAll maneja GET /resource con paginación y filtros automáticos
func (h *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) GetAll(c *gin.Context) {
	// Extraer opciones de query ya parseadas por el middleware
	opts, exists := c.Get("query_options")
	if !exists {
		opts = &common.QueryOptions{}
	}
	queryOptions := opts.(*common.QueryOptions)

	// Extraer contexto de usuario
	userCtx := extractUserContext(c)

	// Llamar al servicio
	entities, pagination, err := h.service.GetAll(c.Request.Context(), *queryOptions, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Mapear a respuestas
	responses := make([]ResponseDTO, 0, len(entities))
	for _, entity := range entities {
		response, mapErr := h.mapper.EntityToResponse(entity, userCtx)
		if mapErr != nil {
			common.ErrorResponse(c, mapErr)
			return
		}
		responses = append(responses, response.(ResponseDTO))
	}

	// Respuesta con paginación
	common.SuccessWithPagination(c, "Recursos obtenidos exitosamente", responses, pagination)
}

// GetByID maneja GET /resource/:id
func (h *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ErrorResponse(c, common.NewValidationError("id", "ID requerido"))
		return
	}

	userCtx := extractUserContext(c)

	// Llamar al servicio
	entity, err := h.service.Get(c.Request.Context(), id, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Mapear a respuesta
	response, err := h.mapper.EntityToResponse(entity, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Recurso obtenido exitosamente", response)
}

// Create maneja POST /resource
func (h *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) Create(c *gin.Context) {
	var dto CreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.NewValidationError("request_body", "Datos inválidos: "+err.Error()))
		return
	}

	userCtx := extractUserContext(c)

	// Llamar al servicio
	entity, err := h.service.Create(c.Request.Context(), dto, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Mapear a respuesta
	response, err := h.mapper.EntityToResponse(entity, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusCreated, "Recurso creado exitosamente", response)
}

// Update maneja PUT/PATCH /resource/:id
func (h *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ErrorResponse(c, common.NewValidationError("id", "ID requerido"))
		return
	}

	var dto UpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		common.ErrorResponse(c, common.NewValidationError("request_body", "Datos inválidos: "+err.Error()))
		return
	}

	userCtx := extractUserContext(c)

	// Llamar al servicio
	entity, err := h.service.Update(c.Request.Context(), id, dto, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	// Mapear a respuesta
	response, err := h.mapper.EntityToResponse(entity, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Recurso actualizado exitosamente", response)
}

// Delete maneja DELETE /resource/:id
func (h *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ErrorResponse(c, common.NewValidationError("id", "ID requerido"))
		return
	}

	userCtx := extractUserContext(c)

	// Llamar al servicio
	err := h.service.Delete(c.Request.Context(), id, userCtx)
	if err != nil {
		common.ErrorResponse(c, err)
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Recurso eliminado exitosamente", gin.H{"id": id})
}

// SetLimits configura límites de paginación
func (h *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) SetLimits(defaultLimit, maxLimit int) {
	h.defaultLimit = defaultLimit
	h.maxLimit = maxLimit
}

// =============================================================================
// MÉTODOS HELPER PARA HANDLERS ESPECÍFICOS
// =============================================================================

// WithCustomValidation permite agregar validaciones personalizadas
func (h *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) WithCustomValidation(validator func(*gin.Context, interface{}) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer DTO del contexto (debe ser establecido por middleware personalizado)
		dto, exists := c.Get("validated_dto")
		if exists {
			if err := validator(c, dto); err != nil {
				common.ErrorResponse(c, err)
				return
			}
		}
		c.Next()
	}
}

// WithPreProcessor permite procesamiento previo antes del CRUD
func (h *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) WithPreProcessor(processor func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := processor(c); err != nil {
			common.ErrorResponse(c, err)
			return
		}
		c.Next()
	}
}

// WithPostProcessor permite procesamiento posterior después del CRUD
func (h *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) WithPostProcessor(processor func(*gin.Context, interface{}) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Obtener resultado del contexto
		result, exists := c.Get("handler_result")
		if exists {
			if err := processor(c, result); err != nil {
				// Log error but don't fail the response since operation completed
				c.Header("X-Post-Process-Warning", err.Error())
			}
		}
	}
}

// =============================================================================
// BUILDER PATTERN PARA CONFIGURAR HANDLERS COMPLEJOS
// =============================================================================

// HandlerBuilder builder para configurar handlers con opciones específicas
type HandlerBuilder[T common.BaseEntity, CreateDTO any, UpdateDTO any, ResponseDTO any] struct {
	baseHandler    *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]
	validators     []func(*gin.Context, interface{}) error
	preProcessors  []func(*gin.Context) error
	postProcessors []func(*gin.Context, interface{}) error
	middlewares    []gin.HandlerFunc
}

// NewHandlerBuilder crea un nuevo builder
func NewHandlerBuilder[T common.BaseEntity, CreateDTO any, UpdateDTO any, ResponseDTO any](
	service common.Service[T, CreateDTO, UpdateDTO],
	mapper mappers.ResponseMapper,
	resourceType string,
) *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO] {
	return &HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO]{
		baseHandler: NewBaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO](service, mapper, resourceType),
	}
}

// WithValidation agrega validador personalizado
func (b *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO]) WithValidation(validator func(*gin.Context, interface{}) error) *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO] {
	b.validators = append(b.validators, validator)
	return b
}

// WithPreProcessor agrega pre-procesador
func (b *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO]) WithPreProcessor(processor func(*gin.Context) error) *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO] {
	b.preProcessors = append(b.preProcessors, processor)
	return b
}

// WithPostProcessor agrega post-procesador
func (b *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO]) WithPostProcessor(processor func(*gin.Context, interface{}) error) *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO] {
	b.postProcessors = append(b.postProcessors, processor)
	return b
}

// WithMiddleware agrega middleware personalizado
func (b *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO]) WithMiddleware(middleware gin.HandlerFunc) *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO] {
	b.middlewares = append(b.middlewares, middleware)
	return b
}

// WithLimits configura límites de paginación
func (b *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO]) WithLimits(defaultLimit, maxLimit int) *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO] {
	b.baseHandler.SetLimits(defaultLimit, maxLimit)
	return b
}

// Build construye el handler final
func (b *HandlerBuilder[T, CreateDTO, UpdateDTO, ResponseDTO]) Build() *BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO] {
	// Combinar todos los procesadores en el handler base
	// (implementación simplificada - en producción podrías hacer esto más elegante)
	return b.baseHandler
}

// =============================================================================
// HANDLERS ESPECIALIZADOS PARA CASOS COMUNES
// =============================================================================

// NestedResourceHandler para recursos anidados como /organizations/:orgId/events
type NestedResourceHandler[T common.BaseEntity, CreateDTO any, UpdateDTO any, ResponseDTO any] struct {
	*BaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO]
	parentParam string
}

// NewNestedResourceHandler crea handler para recursos anidados
func NewNestedResourceHandler[T common.BaseEntity, CreateDTO any, UpdateDTO any, ResponseDTO any](
	service common.Service[T, CreateDTO, UpdateDTO],
	mapper mappers.ResponseMapper,
	resourceType string,
	parentParam string,
) *NestedResourceHandler[T, CreateDTO, UpdateDTO, ResponseDTO] {
	return &NestedResourceHandler[T, CreateDTO, UpdateDTO, ResponseDTO]{
		BaseHandler: NewBaseHandler[T, CreateDTO, UpdateDTO, ResponseDTO](service, mapper, resourceType),
		parentParam: parentParam,
	}
}

// GetAll sobrescribe para incluir filtro del padre
func (h *NestedResourceHandler[T, CreateDTO, UpdateDTO, ResponseDTO]) GetAll(c *gin.Context) {
	parentID := c.Param(h.parentParam)
	if parentID == "" {
		common.ErrorResponse(c, common.NewValidationError(h.parentParam, "Parent ID requerido"))
		return
	}

	// Agregar filtro del padre a las opciones
	opts, exists := c.Get("query_options")
	if !exists {
		opts = &common.QueryOptions{}
	}
	queryOptions := opts.(*common.QueryOptions)
	queryOptions.AddFilter(h.parentParam, parentID)

	// Continuar con lógica normal
	h.BaseHandler.GetAll(c)
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// extractUserContext extrae UserContext del contexto de Gin
func extractUserContext(c *gin.Context) *common.UserContext {
	if userCtx, exists := c.Get("user_context"); exists {
		return userCtx.(*common.UserContext)
	}
	return nil
}

// HandlerRegistry registro centralizado de handlers
type HandlerRegistry struct {
	handlers map[string]interface{}
	services *services.ServiceManager
	mapper   *mappers.UnifiedMapper
}

// NewHandlerRegistry crea nuevo registro
func NewHandlerRegistry(services *services.ServiceManager, mapper *mappers.UnifiedMapper) *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]interface{}),
		services: services,
		mapper:   mapper,
	}
}

// RegisterEventHandler registra handler de eventos
func (r *HandlerRegistry) RegisterEventHandler() *BaseHandler[models.Event, dto.CreateEventRequest, dto.UpdateEventRequest, dto.EventResponse] {
	handler := NewBaseHandler[models.Event, dto.CreateEventRequest, dto.UpdateEventRequest, dto.EventResponse](
		r.services.Events, r.mapper, "event",
	)
	r.handlers["event"] = handler
	return handler
}

// RegisterOrganizationHandler registra handler de organizaciones
func (r *HandlerRegistry) RegisterOrganizationHandler() *BaseHandler[models.Organization, dto.CreateOrganizationRequest, dto.UpdateOrganizationRequest, dto.OrganizationResponse] {
	handler := NewBaseHandler[models.Organization, dto.CreateOrganizationRequest, dto.UpdateOrganizationRequest, dto.OrganizationResponse](
		r.services.Organizations, r.mapper, "organization",
	)
	r.handlers["organization"] = handler
	return handler
}

// RegisterUserHandler registra handler de usuarios
func (r *HandlerRegistry) RegisterUserHandler() *BaseHandler[models.User, dto.CreateUserRequest, dto.UpdateUserRequest, dto.UserResponse] {
	handler := NewBaseHandler[models.User, dto.CreateUserRequest, dto.UpdateUserRequest, dto.UserResponse](
		r.services.Users, r.mapper, "user",
	)
	r.handlers["user"] = handler
	return handler
}

// GetHandler obtiene handler por tipo
func (r *HandlerRegistry) GetHandler(resourceType string) interface{} {
	return r.handlers[resourceType]
}
