// internal/middleware/query_options.go
package middleware

import (
	"github.com/gin-gonic/gin"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/permissions"
	"cybesphere-backend/pkg/database"
)

// QueryOptions middleware para parsear opciones de query
func QueryOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		opts := &common.QueryOptions{
			Page:     1,
			Limit:    20,
			OrderBy:  "created_at",
			OrderDir: "desc",
			Filters:  make(map[string]interface{}),
		}

		// Parsear desde query params
		if err := c.ShouldBindQuery(opts); err == nil {
			// Manejar el error de validación
			if err := opts.Validate(); err != nil {
				// Log del error si es necesario, o simplemente usar valores por defecto
				// Si quieres que el error sea crítico, puedes hacer:
				// c.JSON(400, gin.H{"error": "Invalid query options: " + err.Error()})
				// c.Abort()
				// return

				// O simplemente resetear a valores por defecto seguros
				opts.Page = 1
				opts.Limit = 20
				opts.OrderBy = "created_at"
				opts.OrderDir = "desc"
			}
		}

		// Extraer filtros adicionales de query params
		parseFiltersFromQuery(c, opts)

		// Si hay un usuario autenticado, agregar el contexto
		if userCtx, exists := c.Get("user_context"); exists {
			if ctx, ok := userCtx.(*common.UserContext); ok {
				opts.UserContext = ctx
			}
		}

		c.Set("query_options", opts)
		c.Next()
	}
}

// EnhancedUserContext middleware para crear contexto completo de usuario desde DB
func EnhancedUserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Solo ejecutar si el usuario está autenticado
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok || userIDStr == "" {
			c.Next()
			return
		}

		userCtx := buildUserContextFromDB(c, userIDStr)
		if userCtx != nil {
			c.Set("user_context", userCtx)
		}

		c.Next()
	}
}

// buildUserContextFromDB construye un UserContext completo desde la base de datos
func buildUserContextFromDB(c *gin.Context, userID string) *common.UserContext {
	db := database.GetDB()
	if db == nil {
		return buildBasicUserContext(c)
	}

	permChecker := permissions.NewPermissionChecker()

	// Buscar usuario en DB para obtener info completa
	var user models.User
	if err := db.Preload("Organization").First(&user, "id = ?", userID).Error; err != nil {
		// Si hay cualquier error (incluido ErrRecordNotFound), usar info básica del contexto
		return buildBasicUserContext(c)
	}

	userCtx := &common.UserContext{
		ID:           user.ID.String(),
		Email:        user.Email,
		Role:         user.Role,
		IsActive:     user.IsActive,
		IsVerified:   user.IsVerified,
		Permissions:  permChecker.GetRolePermissions(user.Role),
		Capabilities: permChecker.GetRoleCapabilities(user.Role),
	}

	if user.OrganizationID != nil {
		userCtx.OrganizationID = user.OrganizationID
	}

	return userCtx
}

// buildBasicUserContext construye contexto básico desde valores de Gin
func buildBasicUserContext(c *gin.Context) *common.UserContext {
	permChecker := permissions.NewPermissionChecker()

	userCtx := &common.UserContext{
		Permissions:  []permissions.Permission{},
		Capabilities: make(map[string]bool),
	}

	hasData := false

	// Extraer ID de usuario
	if id, exists := c.Get("user_id"); exists {
		if userID, ok := id.(string); ok && userID != "" {
			userCtx.ID = userID
			hasData = true
		}
	}

	// Extraer email
	if email, exists := c.Get("user_email"); exists {
		if userEmail, ok := email.(string); ok && userEmail != "" {
			userCtx.Email = userEmail
			hasData = true
		}
	}

	// Extraer rol y cargar permisos
	if role, exists := c.Get("user_role"); exists {
		if userRole, ok := role.(string); ok && userRole != "" {
			userCtx.Role = models.UserRole(userRole)
			hasData = true

			// Cargar permisos y capacidades basados en el rol
			userCtx.Permissions = permChecker.GetRolePermissions(userCtx.Role)
			userCtx.Capabilities = permChecker.GetRoleCapabilities(userCtx.Role)
		}
	}

	// Extraer ID de organización
	if orgID, exists := c.Get("organization_id"); exists {
		if id, ok := orgID.(string); ok && id != "" {
			userCtx.OrganizationID = &id
		}
	}

	// Extraer estado activo
	if isActive, exists := c.Get("user_active"); exists {
		if active, ok := isActive.(bool); ok {
			userCtx.IsActive = active
		}
	}

	// Extraer estado verificado
	if isVerified, exists := c.Get("user_verified"); exists {
		if verified, ok := isVerified.(bool); ok {
			userCtx.IsVerified = verified
		}
	}

	if !hasData {
		return nil
	}

	return userCtx
}

// parseFiltersFromQuery extrae filtros adicionales de query params
func parseFiltersFromQuery(c *gin.Context, opts *common.QueryOptions) {
	// Lista de parámetros reservados que NO son filtros
	reservedParams := map[string]bool{
		"page":      true,
		"limit":     true,
		"order_by":  true,
		"order_dir": true,
		"search":    true,
		"expand":    true,
		"fields":    true,
	}

	queryParams := c.Request.URL.Query()

	for key, values := range queryParams {
		// Saltar parámetros reservados o sin valores
		if reservedParams[key] || len(values) == 0 {
			continue
		}

		value := values[0]

		// Ignorar valores vacíos o nulos
		if value == "" || value == "null" {
			continue
		}

		// Convertir valores booleanos
		switch value {
		case "true":
			opts.AddFilter(key, true)
		case "false":
			opts.AddFilter(key, false)
		default:
			opts.AddFilter(key, value)
		}
	}
}

// GetQueryOptions helper para extraer QueryOptions del contexto
func GetQueryOptions(c *gin.Context) *common.QueryOptions {
	if opts, exists := c.Get("query_options"); exists {
		if queryOpts, ok := opts.(*common.QueryOptions); ok {
			return queryOpts
		}
	}

	// Retornar opciones por defecto si no existen
	return &common.QueryOptions{
		Page:     1,
		Limit:    20,
		OrderBy:  "created_at",
		OrderDir: "desc",
		Filters:  make(map[string]interface{}),
	}
}

// GetUserContext helper para extraer UserContext del contexto
func GetUserContext(c *gin.Context) *common.UserContext {
	if userCtx, exists := c.Get("user_context"); exists {
		if ctx, ok := userCtx.(*common.UserContext); ok {
			return ctx
		}
	}
	return nil
}
