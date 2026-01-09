// internal/handlers/helpers.go
package handlers

import (
	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/pkg/logger"
)

func GetUserContext(c *gin.Context) *common.UserContext {
	return extractUserContext(c)
}

// extractQueryOptions extrae opciones de query del contexto o las construye
func extractQueryOptions(c *gin.Context) *common.QueryOptions {
	// Primero intentar obtener las opciones ya parseadas
	if opts, exists := c.Get("query_options"); exists {
		if queryOpts, ok := opts.(*common.QueryOptions); ok {
			return queryOpts
		}
	}

	// Si no existen, construir opciones básicas
	opts := &common.QueryOptions{
		Page:     1,
		Limit:    20,
		OrderBy:  "created_at",
		OrderDir: "desc",
		Filters:  make(map[string]interface{}),
	}

	// Parsear de query params
	if err := c.ShouldBindQuery(opts); err == nil {
		if err := opts.Validate(); err != nil {
			logger.Warn("Error validando opciones de query: ", err)
		}
	}

	// Agregar contexto de usuario si existe
	opts.UserContext = extractUserContext(c)

	return opts
}
