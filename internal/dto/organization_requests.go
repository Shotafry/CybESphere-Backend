package dto

// CreateOrganizationRequest DTO para crear una organización
type CreateOrganizationRequest struct {
	// Información básica
	Name        string `json:"name" binding:"required,min=3,max=200"`
	Description string `json:"description" binding:"required,min=10"`
	Website     string `json:"website" binding:"omitempty,url,max=255"`

	// Información de contacto
	Email      string `json:"email" binding:"required,email,max=255"`
	Phone      string `json:"phone" binding:"max=20"`
	Address    string `json:"address" binding:"max=500"`
	City       string `json:"city" binding:"max=100"`
	Country    string `json:"country" binding:"max=100"`
	PostalCode string `json:"postal_code" binding:"max=20"`

	// Geolocalización
	Latitude  *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`

	// Branding y medios
	LogoURL        string `json:"logo_url" binding:"omitempty,url,max=500"`
	BannerURL      string `json:"banner_url" binding:"omitempty,url,max=500"`
	PrimaryColor   string `json:"primary_color" binding:"omitempty,hexcolor"`
	SecondaryColor string `json:"secondary_color" binding:"omitempty,hexcolor"`

	// Redes sociales
	LinkedIn  string `json:"linkedin" binding:"omitempty,url,max=255"`
	Twitter   string `json:"twitter" binding:"omitempty,url,max=255"`
	Facebook  string `json:"facebook" binding:"omitempty,url,max=255"`
	Instagram string `json:"instagram" binding:"omitempty,url,max=255"`
	YouTube   string `json:"youtube" binding:"omitempty,url,max=255"`

	// Documentación para verificación
	TaxID            string `json:"tax_id" binding:"max=50"`
	LegalName        string `json:"legal_name" binding:"max=300"`
	RegistrationDocs string `json:"registration_docs" binding:"omitempty,url,max=500"`

	// Admin only
	Status          string `json:"status" binding:"omitempty,oneof=pending active suspended inactive"`
	IsVerified      *bool  `json:"is_verified"`
	MaxEvents       *int   `json:"max_events" binding:"omitempty,min=0"`
	CanCreateEvents *bool  `json:"can_create_events"`
}

// UpdateOrganizationRequest DTO para actualizar una organización
type UpdateOrganizationRequest struct {
	// Información básica
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=200"`
	Description *string `json:"description,omitempty" binding:"omitempty,min=10"`
	Website     *string `json:"website,omitempty" binding:"omitempty,url,max=255"`

	// Información de contacto
	Email      *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
	Phone      *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	Address    *string `json:"address,omitempty" binding:"omitempty,max=500"`
	City       *string `json:"city,omitempty" binding:"omitempty,max=100"`
	Country    *string `json:"country,omitempty" binding:"omitempty,max=100"`
	PostalCode *string `json:"postal_code,omitempty" binding:"omitempty,max=20"`

	// Geolocalización
	Latitude  *float64 `json:"latitude,omitempty" binding:"omitempty,min=-90,max=90"`
	Longitude *float64 `json:"longitude,omitempty" binding:"omitempty,min=-180,max=180"`

	// Branding y medios
	LogoURL        *string `json:"logo_url,omitempty" binding:"omitempty,url,max=500"`
	BannerURL      *string `json:"banner_url,omitempty" binding:"omitempty,url,max=500"`
	PrimaryColor   *string `json:"primary_color,omitempty" binding:"omitempty,hexcolor"`
	SecondaryColor *string `json:"secondary_color,omitempty" binding:"omitempty,hexcolor"`

	// Redes sociales
	LinkedIn  *string `json:"linkedin,omitempty" binding:"omitempty,url,max=255"`
	Twitter   *string `json:"twitter,omitempty" binding:"omitempty,url,max=255"`
	Facebook  *string `json:"facebook,omitempty" binding:"omitempty,url,max=255"`
	Instagram *string `json:"instagram,omitempty" binding:"omitempty,url,max=255"`
	YouTube   *string `json:"youtube,omitempty" binding:"omitempty,url,max=255"`

	// Admin only
	Status          *string `json:"status,omitempty" binding:"omitempty,oneof=pending active suspended inactive"`
	IsVerified      *bool   `json:"is_verified,omitempty"`
	MaxEvents       *int    `json:"max_events,omitempty" binding:"omitempty,min=0"`
	CanCreateEvents *bool   `json:"can_create_events,omitempty"`
}

// VerifyOrganizationRequest DTO para verificar una organización
type VerifyOrganizationRequest struct {
	Notes string `json:"notes" binding:"max=500"`
}

// OrganizationFilterRequest DTO para filtrar organizaciones
type OrganizationFilterRequest struct {
	// Filtros básicos
	Status     string `form:"status" binding:"omitempty,oneof=pending active suspended inactive"`
	IsVerified *bool  `form:"is_verified"`

	// Filtros de ubicación
	City    string `form:"city"`
	Country string `form:"country"`

	// Búsqueda y ordenamiento
	Search   string `form:"search"`
	OrderBy  string `form:"order_by" binding:"omitempty,oneof=name created_at updated_at events_count city"`
	OrderDir string `form:"order_dir" binding:"omitempty,oneof=asc desc"`

	// Paginación
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}
