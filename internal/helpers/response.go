// internal/helpers/response.go
package helpers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"cybesphere-backend/internal/common"
)

// ============================================
// Funciones de Formateo de Respuesta
// ============================================

// FormatPaginationResponse formatea una respuesta con paginación estándar
func FormatPaginationResponse(c *gin.Context, data interface{}, meta *common.PaginationMeta, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}

// FormatSuccessResponse formatea una respuesta exitosa simple
func FormatSuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// FormatErrorResponse formatea una respuesta de error estándar
func FormatErrorResponse(c *gin.Context, statusCode int, errorCode, message string) {
	c.JSON(statusCode, gin.H{
		"success":   false,
		"error":     errorCode,
		"message":   message,
		"timestamp": time.Now().UTC(),
	})
}

// FormatValidationErrorResponse formatea errores de validación
func FormatValidationErrorResponse(c *gin.Context, errorDetail interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "validation_error",
		"message": "La validación de los datos ha fallado",
		"details": errorDetail,
	})
}

// FormatCreatedResponse formatea una respuesta de recurso creado
func FormatCreatedResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// FormatNoContentResponse responde sin contenido (204)
func FormatNoContentResponse(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// FormatAcceptedResponse formatea una respuesta de proceso aceptado (202)
func FormatAcceptedResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"status":  "processing",
	})
}

// ============================================
// Funciones de Paginación
// ============================================

// ValidatePaginationParams valida y normaliza parámetros de paginación
func ValidatePaginationParams(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

// ExtractPaginationParams extrae parámetros de paginación del contexto
func ExtractPaginationParams(c *gin.Context) (page, limit int) {
	page = 1
	limit = 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	return ValidatePaginationParams(page, limit)
}

// CalculateOffset calcula el offset para paginación
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculateTotalPages calcula el número total de páginas
func CalculateTotalPages(total int64, limit int) int64 {
	if limit <= 0 {
		return 0
	}
	return (total + int64(limit) - 1) / int64(limit)
}

// BuildPaginationMeta construye metadata de paginación
func BuildPaginationMeta(page, limit int, total int64) *common.PaginationMeta {
	totalPages := CalculateTotalPages(total, limit)

	return &common.PaginationMeta{
		Page:    page,
		Limit:   limit,
		Total:   total,
		Pages:   totalPages,
		HasNext: int64(page) < totalPages,
		HasPrev: page > 1,
	}
}

// ============================================
// Funciones de Ordenamiento y Filtrado
// ============================================

// ExtractSortParams extrae y valida parámetros de ordenamiento
func ExtractSortParams(c *gin.Context, validFields map[string]bool) (orderBy, orderDir string) {
	orderBy = c.DefaultQuery("order_by", "created_at")
	orderDir = c.DefaultQuery("order_dir", "desc")

	// Validar dirección de ordenamiento
	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "desc"
	}

	// Si no se proporcionan campos válidos, usar conjunto por defecto
	if validFields == nil {
		validFields = DefaultSortFields()
	}

	// Validar campo de ordenamiento contra lista blanca
	if !validFields[orderBy] {
		orderBy = "created_at"
	}

	return orderBy, orderDir
}

// DefaultSortFields retorna los campos de ordenamiento por defecto
func DefaultSortFields() map[string]bool {
	return map[string]bool{
		"created_at":  true,
		"updated_at":  true,
		"name":        true,
		"title":       true,
		"start_date":  true,
		"end_date":    true,
		"status":      true,
		"is_public":   true,
		"is_verified": true,
		"is_active":   true,
		"email":       true,
		"role":        true,
	}
}

// ExtractSearchQuery extrae el parámetro de búsqueda
func ExtractSearchQuery(c *gin.Context) string {
	return c.Query("search")
}

// ExtractFilters extrae filtros comunes de la query
func ExtractFilters(c *gin.Context, allowedFilters []string) map[string]interface{} {
	filters := make(map[string]interface{})

	// Crear mapa de filtros permitidos para búsqueda rápida
	allowed := make(map[string]bool)
	for _, filter := range allowedFilters {
		allowed[filter] = true
	}

	// Extraer solo los filtros permitidos
	for key := range allowed {
		if value := c.Query(key); value != "" {
			// Convertir valores booleanos
			if value == "true" || value == "false" {
				filters[key] = value == "true"
			} else {
				filters[key] = value
			}
		}
	}

	return filters
}

// ============================================
// Funciones de Manejo de Errores
// ============================================

// HandleError maneja errores comunes y responde apropiadamente
func HandleError(c *gin.Context, err error) {
	switch err {
	case gorm.ErrRecordNotFound:
		FormatErrorResponse(c, http.StatusNotFound, "not_found", "Recurso no encontrado")
	case common.ErrUnauthorized:
		FormatErrorResponse(c, http.StatusUnauthorized, "unauthorized", "No autorizado")
	case common.ErrForbidden:
		FormatErrorResponse(c, http.StatusForbidden, "forbidden", "Acceso denegado")
	case common.ErrValidation:
		FormatErrorResponse(c, http.StatusBadRequest, "validation_error", err.Error())
	case common.ErrDuplicateEntry:
		FormatErrorResponse(c, http.StatusConflict, "duplicate_entry", "El recurso ya existe")
	case common.ErrInvalidInput:
		FormatErrorResponse(c, http.StatusBadRequest, "invalid_input", err.Error())
	default:
		// Para errores no esperados, loggear pero no exponer detalles
		// log.Printf("Internal error: %v", err)
		FormatErrorResponse(c, http.StatusInternalServerError, "internal_error", "Error interno del servidor")
	}
}

// HandleDatabaseError maneja errores específicos de base de datos
func HandleDatabaseError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Verificar si es error de registro no encontrado
	if err == gorm.ErrRecordNotFound {
		FormatErrorResponse(c, http.StatusNotFound, "not_found", "Registro no encontrado")
		return
	}

	// Para otros errores de DB, responder con error genérico
	FormatErrorResponse(c, http.StatusInternalServerError, "database_error", "Error al procesar la solicitud")
}

// ============================================
// Funciones de Respuesta para Listas
// ============================================

// FormatListResponse formatea una respuesta de lista con paginación opcional
func FormatListResponse(c *gin.Context, data interface{}, total int64, page, limit int, message string) {
	if page > 0 && limit > 0 {
		// Con paginación
		meta := BuildPaginationMeta(page, limit, total)
		FormatPaginationResponse(c, data, meta, message)
	} else {
		// Sin paginación
		FormatSuccessResponse(c, data, message)
	}
}

// FormatEmptyListResponse formatea una respuesta para lista vacía
func FormatEmptyListResponse(c *gin.Context, message string) {
	FormatSuccessResponse(c, []interface{}{}, message)
}

// ============================================
// Funciones de Utilidad
// ============================================

// ExtractIDParam extrae y valida un parámetro ID de la URL
func ExtractIDParam(c *gin.Context, paramName string) (string, error) {
	id := c.Param(paramName)
	if id == "" {
		return "", common.ErrInvalidInput
	}
	return id, nil
}

// ExtractBoolQuery extrae un parámetro booleano de la query
func ExtractBoolQuery(c *gin.Context, paramName string, defaultValue bool) bool {
	value := c.Query(paramName)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1"
}

// ExtractIntQuery extrae un parámetro entero de la query
func ExtractIntQuery(c *gin.Context, paramName string, defaultValue int) int {
	value := c.Query(paramName)
	if value == "" {
		return defaultValue
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}

	return defaultValue
}

// BuildFilterQuery construye una query de filtrado para GORM
func BuildFilterQuery(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}
	return query
}

// BuildSearchQuery construye una query de búsqueda para GORM
func BuildSearchQuery(query *gorm.DB, search string, searchFields []string) *gorm.DB {
	if search == "" || len(searchFields) == 0 {
		return query
	}

	searchPattern := "%" + search + "%"
	condition := ""
	values := []interface{}{}

	for i, field := range searchFields {
		if i > 0 {
			condition += " OR "
		}
		condition += field + " ILIKE ?"
		values = append(values, searchPattern)
	}

	return query.Where(condition, values...)
}

// GetUserIDFromContext extrae el ID del usuario del contexto
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// GetOrganizationIDFromContext extrae el ID de la organización del contexto
func GetOrganizationIDFromContext(c *gin.Context) (string, bool) {
	orgID, exists := c.Get("organization_id")
	if !exists {
		return "", false
	}

	id, ok := orgID.(string)
	return id, ok
}
