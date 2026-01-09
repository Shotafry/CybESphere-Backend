package common

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

// Errores de dominio
var (
	ErrNotFound       = errors.New("recurso no encontrado")
	ErrUnauthorized   = errors.New("no autorizado")
	ErrForbidden      = errors.New("acceso denegado")
	ErrValidation     = errors.New("datos inválidos")
	ErrConflict       = errors.New("conflicto de datos")
	ErrInternalError  = errors.New("error interno del servidor")
	ErrDuplicateEntry = errors.New("entrada duplicada")
	ErrInvalidInput   = errors.New("entrada inválida")
)

// BusinessError error de negocio con contexto
type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Field   string `json:"field,omitempty"`
}

func (e BusinessError) Error() string {
	return e.Message
}

// NewValidationError crea un error de validación
func NewValidationError(field, message string) *BusinessError {
	return &BusinessError{
		Code:    "validation_error",
		Message: message,
		Field:   field,
	}
}

// NewBusinessError crea un error de negocio genérico
func NewBusinessError(code, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// MapGormError mapea errores de GORM a errores de dominio
func MapGormError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	// Otros mapeos de errores GORM según sea necesario
	return ErrInternalError
}

// HTTPStatusFromError mapea errores a códigos HTTP
func HTTPStatusFromError(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrValidation):
		return http.StatusBadRequest
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
