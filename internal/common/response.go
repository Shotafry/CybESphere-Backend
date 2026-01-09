package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse estructura estándar de respuesta
type APIResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message,omitempty"`
	Data       interface{}     `json:"data,omitempty"`
	Error      interface{}     `json:"error,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

// SuccessResponse respuesta de éxito
func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPagination respuesta de éxito con paginación
func SuccessWithPagination(c *gin.Context, message string, data interface{}, pagination *PaginationMeta) {
	c.JSON(http.StatusOK, APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// ErrorResponse respuesta de error
func ErrorResponse(c *gin.Context, err error) {
	status := HTTPStatusFromError(err)

	var errorData interface{} = map[string]string{
		"message": err.Error(),
	}

	// Si es un BusinessError, incluir detalles adicionales
	if businessErr, ok := err.(*BusinessError); ok {
		errorData = businessErr
	}

	c.JSON(status, APIResponse{
		Success: false,
		Error:   errorData,
	})
}
