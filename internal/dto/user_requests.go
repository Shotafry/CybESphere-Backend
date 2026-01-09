package dto

// CreateUserRequest DTO para crear un usuario (registro)
type CreateUserRequest struct {
	// Información básica
	Email     string `json:"email" binding:"required,email,max=255"`
	Password  string `json:"password" binding:"required,min=8,max=128"`
	FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string `json:"last_name" binding:"required,min=1,max=100"`

	// Perfil profesional
	Company  string `json:"company" binding:"max=200"`
	Position string `json:"position" binding:"max=200"`
	Bio      string `json:"bio" binding:"max=1000"`
	Website  string `json:"website" binding:"omitempty,url,max=255"`
	LinkedIn string `json:"linkedin" binding:"omitempty,url,max=255"`
	Twitter  string `json:"twitter" binding:"omitempty,url,max=255"`

	// Ubicación
	City      string   `json:"city" binding:"max=100"`
	Country   string   `json:"country" binding:"max=100"`
	Latitude  *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`

	// Configuraciones
	Timezone          string `json:"timezone" binding:"max=50"`
	Language          string `json:"language" binding:"omitempty,oneof=es en"`
	NewsletterEnabled bool   `json:"newsletter_enabled"`

	// Admin only (para crear usuarios desde admin)
	Role           string  `json:"role" binding:"omitempty,oneof=user organizer admin"`
	IsActive       *bool   `json:"is_active"`
	IsVerified     *bool   `json:"is_verified"`
	OrganizationID *string `json:"organization_id" binding:"omitempty,uuid"`
}

// UpdateUserRequest DTO para actualizar un usuario
type UpdateUserRequest struct {
	// Información básica (email no se puede cambiar)
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=1,max=100"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=1,max=100"`

	// Perfil profesional
	Company  *string `json:"company,omitempty" binding:"omitempty,max=200"`
	Position *string `json:"position,omitempty" binding:"omitempty,max=200"`
	Bio      *string `json:"bio,omitempty" binding:"omitempty,max=1000"`
	Website  *string `json:"website,omitempty" binding:"omitempty,url,max=255"`
	LinkedIn *string `json:"linkedin,omitempty" binding:"omitempty,url,max=255"`
	Twitter  *string `json:"twitter,omitempty" binding:"omitempty,url,max=255"`

	// Ubicación
	City      *string  `json:"city,omitempty" binding:"omitempty,max=100"`
	Country   *string  `json:"country,omitempty" binding:"omitempty,max=100"`
	Latitude  *float64 `json:"latitude,omitempty" binding:"omitempty,min=-90,max=90"`
	Longitude *float64 `json:"longitude,omitempty" binding:"omitempty,min=-180,max=180"`

	// Configuraciones
	Timezone          *string `json:"timezone,omitempty" binding:"omitempty,max=50"`
	Language          *string `json:"language,omitempty" binding:"omitempty,oneof=es en"`
	NewsletterEnabled *bool   `json:"newsletter_enabled,omitempty"`

	// Admin only
	Role           *string `json:"role,omitempty" binding:"omitempty,oneof=user organizer admin"`
	IsActive       *bool   `json:"is_active,omitempty"`
	IsVerified     *bool   `json:"is_verified,omitempty"`
	OrganizationID *string `json:"organization_id,omitempty" binding:"omitempty,uuid"`
}

// ChangePasswordRequest DTO para cambiar contraseña
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// UpdateUserRoleRequest DTO para actualizar rol de usuario (admin only)
type UpdateUserRoleRequest struct {
	Role   string `json:"role" binding:"required,oneof=user organizer admin"`
	Reason string `json:"reason" binding:"max=500"`
}

// UserFilterRequest DTO para filtrar usuarios
type UserFilterRequest struct {
	// Filtros básicos
	Role       string `form:"role" binding:"omitempty,oneof=user organizer admin"`
	IsActive   *bool  `form:"is_active"`
	IsVerified *bool  `form:"is_verified"`

	// Filtros de organización
	OrganizationID string `form:"organization_id" binding:"omitempty,uuid"`

	// Filtros de ubicación
	City    string `form:"city"`
	Country string `form:"country"`

	// Búsqueda y ordenamiento
	Search   string `form:"search"`
	OrderBy  string `form:"order_by" binding:"omitempty,oneof=created_at updated_at email first_name last_name last_login_at"`
	OrderDir string `form:"order_dir" binding:"omitempty,oneof=asc desc"`

	// Paginación
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

// ActivateUserRequest DTO para activar/desactivar usuario
type ActivateUserRequest struct {
	IsActive bool   `json:"is_active"`
	Reason   string `json:"reason" binding:"max=500"`
}

// VerifyUserRequest DTO para verificar usuario
type VerifyUserRequest struct {
	IsVerified bool   `json:"is_verified"`
	Notes      string `json:"notes" binding:"max=500"`
}
