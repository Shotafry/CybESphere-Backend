package dto

import (
	"time"
)

// CreateEventRequest DTO para crear un evento
type CreateEventRequest struct {
	// Información básica
	Title       string `json:"title" binding:"required,min=5,max=300"`
	Description string `json:"description" binding:"required,min=10"`
	ShortDesc   string `json:"short_desc" binding:"max=500"`

	// Clasificación
	Type     string `json:"type" binding:"required,oneof=conference workshop meetup webinar training competition other"`
	Category string `json:"category" binding:"max=100"`
	Level    string `json:"level" binding:"omitempty,oneof=beginner intermediate advanced"`

	// Fechas y horarios
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
	Timezone  string    `json:"timezone" binding:"max=50"`

	// Ubicación
	IsOnline     bool     `json:"is_online"`
	VenueAddress string   `json:"venue_address" binding:"required_unless=IsOnline true,max=500"`
	VenueName    string   `json:"venue_name" binding:"max=200"`
	VenueCity    string   `json:"venue_city" binding:"max=100"`
	VenueCountry string   `json:"venue_country" binding:"max=100"`
	Latitude     *float64 `json:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude    *float64 `json:"longitude" binding:"omitempty,min=-180,max=180"`
	OnlineURL    string   `json:"online_url" binding:"required_if=IsOnline true,omitempty,url,max=500"`
	StreamingURL string   `json:"streaming_url" binding:"omitempty,url,max=500"`

	// Capacidad y registro
	MaxAttendees    *int   `json:"max_attendees" binding:"omitempty,min=1"`
	IsFree          bool   `json:"is_free"`
	Price           *int   `json:"price" binding:"omitempty,min=0"`
	Currency        string `json:"currency" binding:"max=3"`
	RegistrationURL string `json:"registration_url" binding:"omitempty,url,max=500"`

	// Contenido y recursos
	ImageURL     string   `json:"image_url" binding:"omitempty,url,max=500"`
	BannerURL    string   `json:"banner_url" binding:"omitempty,url,max=500"`
	Tags         []string `json:"tags" binding:"max=10,dive,min=1,max=50"`
	Requirements string   `json:"requirements" binding:"max=2000"`
	Agenda       string   `json:"agenda" binding:"max=5000"`

	// Fechas importantes
	RegistrationStartDate *time.Time `json:"registration_start_date"`
	RegistrationEndDate   *time.Time `json:"registration_end_date" binding:"omitempty,gtfield=RegistrationStartDate"`

	// Información de contacto
	ContactEmail string `json:"contact_email" binding:"omitempty,email,max=255"`
	ContactPhone string `json:"contact_phone" binding:"max=20"`

	// SEO
	MetaTitle       string `json:"meta_title" binding:"max=200"`
	MetaDescription string `json:"meta_description" binding:"max=500"`

	// Admin fields (only for admin users)
	OrganizationID string `json:"organization_id" binding:"omitempty,uuid"`
	Status         string `json:"status" binding:"omitempty,oneof=draft published"`
	IsPublic       *bool  `json:"is_public"`
	IsFeatured     *bool  `json:"is_featured"`
}

// UpdateEventRequest DTO para actualizar un evento
type UpdateEventRequest struct {
	// Información básica
	Title       *string `json:"title,omitempty" binding:"omitempty,min=5,max=300"`
	Description *string `json:"description,omitempty" binding:"omitempty,min=10"`
	ShortDesc   *string `json:"short_desc,omitempty" binding:"omitempty,max=500"`

	// Clasificación
	Category *string `json:"category,omitempty" binding:"omitempty,max=100"`
	Level    *string `json:"level,omitempty" binding:"omitempty,oneof=beginner intermediate advanced"`

	// Fechas y horarios (no se puede cambiar el tipo de evento)
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Timezone  *string    `json:"timezone,omitempty" binding:"omitempty,max=50"`

	// Ubicación
	IsOnline     *bool    `json:"is_online,omitempty"`
	VenueAddress *string  `json:"venue_address,omitempty" binding:"omitempty,max=500"`
	VenueName    *string  `json:"venue_name,omitempty" binding:"omitempty,max=200"`
	VenueCity    *string  `json:"venue_city,omitempty" binding:"omitempty,max=100"`
	VenueCountry *string  `json:"venue_country,omitempty" binding:"omitempty,max=100"`
	Latitude     *float64 `json:"latitude,omitempty" binding:"omitempty,min=-90,max=90"`
	Longitude    *float64 `json:"longitude,omitempty" binding:"omitempty,min=-180,max=180"`
	OnlineURL    *string  `json:"online_url,omitempty" binding:"omitempty,url,max=500"`
	StreamingURL *string  `json:"streaming_url,omitempty" binding:"omitempty,url,max=500"`

	// Capacidad y registro
	MaxAttendees    *int    `json:"max_attendees,omitempty" binding:"omitempty,min=1"`
	IsFree          *bool   `json:"is_free,omitempty"`
	Price           *int    `json:"price,omitempty" binding:"omitempty,min=0"`
	Currency        *string `json:"currency,omitempty" binding:"omitempty,max=3"`
	RegistrationURL *string `json:"registration_url,omitempty" binding:"omitempty,url,max=500"`

	// Contenido y recursos
	ImageURL     *string  `json:"image_url,omitempty" binding:"omitempty,url,max=500"`
	BannerURL    *string  `json:"banner_url,omitempty" binding:"omitempty,url,max=500"`
	Tags         []string `json:"tags,omitempty" binding:"omitempty,max=10,dive,min=1,max=50"`
	Requirements *string  `json:"requirements,omitempty" binding:"omitempty,max=2000"`
	Agenda       *string  `json:"agenda,omitempty" binding:"omitempty,max=5000"`

	// Fechas importantes
	RegistrationStartDate *time.Time `json:"registration_start_date,omitempty"`
	RegistrationEndDate   *time.Time `json:"registration_end_date,omitempty"`

	// Información de contacto
	ContactEmail *string `json:"contact_email,omitempty" binding:"omitempty,email,max=255"`
	ContactPhone *string `json:"contact_phone,omitempty" binding:"omitempty,max=20"`

	// SEO
	MetaTitle       *string `json:"meta_title,omitempty" binding:"omitempty,max=200"`
	MetaDescription *string `json:"meta_description,omitempty" binding:"omitempty,max=500"`

	// Admin only fields
	Status     *string `json:"status,omitempty" binding:"omitempty,oneof=draft published canceled completed"`
	IsPublic   *bool   `json:"is_public,omitempty"`
	IsFeatured *bool   `json:"is_featured,omitempty"`
}

// PublishEventRequest DTO para publicar un evento
type PublishEventRequest struct {
	SendNotifications bool `json:"send_notifications"`
}

// CancelEventRequest DTO para cancelar un evento
type CancelEventRequest struct {
	Reason            string `json:"reason" binding:"required,min=10,max=500"`
	SendNotifications bool   `json:"send_notifications"`
}

// EventFilterRequest DTO para filtrar eventos
type EventFilterRequest struct {
	// Filtros básicos
	Type     string `form:"type" binding:"omitempty,oneof=conference workshop meetup webinar training competition other"`
	Category string `form:"category"`
	Level    string `form:"level" binding:"omitempty,oneof=beginner intermediate advanced"`

	// Filtros de ubicación
	City      string  `form:"city"`
	Country   string  `form:"country"`
	IsOnline  *bool   `form:"is_online"`
	Latitude  float64 `form:"latitude" binding:"omitempty,min=-90,max=90"`
	Longitude float64 `form:"longitude" binding:"omitempty,min=-180,max=180"`
	RadiusKm  int     `form:"radius_km" binding:"omitempty,min=1,max=500"`

	// Filtros de precio
	IsFree   *bool `form:"is_free"`
	MinPrice *int  `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice *int  `form:"max_price" binding:"omitempty,min=0,gtefield=MinPrice"`

	// Filtros de fecha
	StartDateFrom *time.Time `form:"start_date_from" time_format:"2006-01-02"`
	StartDateTo   *time.Time `form:"start_date_to" time_format:"2006-01-02"`
	EndDateFrom   *time.Time `form:"end_date_from" time_format:"2006-01-02"`
	EndDateTo     *time.Time `form:"end_date_to" time_format:"2006-01-02"`

	// Filtros de estado (admin only)
	Status     string `form:"status" binding:"omitempty,oneof=draft published canceled completed"`
	IsPublic   *bool  `form:"is_public"`
	IsFeatured *bool  `form:"is_featured"`

	// Filtros de organización
	OrganizationID string `form:"organization_id" binding:"omitempty,uuid"`

	// Búsqueda y ordenamiento
	Search   string   `form:"search"`
	Tags     []string `form:"tags"`
	OrderBy  string   `form:"order_by" binding:"omitempty,oneof=start_date end_date created_at updated_at title views_count current_attendees"`
	OrderDir string   `form:"order_dir" binding:"omitempty,oneof=asc desc"`

	// Paginación
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

// AddToFavoritesRequest DTO para agregar a favoritos
type AddToFavoritesRequest struct {
	EventID string `json:"event_id" binding:"required,uuid"`
}

// EventAttendeesFilterRequest DTO para filtrar asistentes
type EventAttendeesFilterRequest struct {
	Status   string `form:"status" binding:"omitempty,oneof=confirmed pending canceled"`
	Search   string `form:"search"`
	OrderBy  string `form:"order_by" binding:"omitempty,oneof=created_at name email"`
	OrderDir string `form:"order_dir" binding:"omitempty,oneof=asc desc"`
	Page     int    `form:"page" binding:"min=1"`
	Limit    int    `form:"limit" binding:"min=1,max=100"`
}
