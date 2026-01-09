package dto

import (
	"cybesphere-backend/internal/common"
	"time"
)

// EventResponse DTO de respuesta para un evento
type EventResponse struct {
	// Identificación
	ID   string `json:"id"`
	Slug string `json:"slug"`

	// Información básica
	Title       string `json:"title"`
	Description string `json:"description"`
	ShortDesc   string `json:"short_desc"`

	// Clasificación
	Type     string `json:"type"`
	Category string `json:"category"`
	Level    string `json:"level,omitempty"`

	// Fechas y horarios
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Timezone  string    `json:"timezone"`
	Duration  int       `json:"duration"` // En minutos

	// Ubicación
	IsOnline     bool     `json:"is_online"`
	VenueName    string   `json:"venue_name,omitempty"`
	VenueAddress string   `json:"venue_address,omitempty"`
	VenueCity    string   `json:"venue_city,omitempty"`
	VenueCountry string   `json:"venue_country,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	OnlineURL    string   `json:"online_url,omitempty"`
	StreamingURL string   `json:"streaming_url,omitempty"`

	// Capacidad y registro
	MaxAttendees     *int   `json:"max_attendees"`
	CurrentAttendees int    `json:"current_attendees"`
	AvailableSpots   *int   `json:"available_spots"`
	IsFree           bool   `json:"is_free"`
	Price            *int   `json:"price,omitempty"`
	Currency         string `json:"currency,omitempty"`
	RegistrationURL  string `json:"registration_url,omitempty"`

	// Estado y visibilidad
	Status     string `json:"status"`
	IsPublic   bool   `json:"is_public"`
	IsFeatured bool   `json:"is_featured"`
	ViewsCount int    `json:"views_count"`

	// Contenido y recursos
	ImageURL  string   `json:"image_url,omitempty"`
	BannerURL string   `json:"banner_url,omitempty"`
	Tags      []string `json:"tags,omitempty"`

	// Información de registro
	RegistrationOpen      bool       `json:"registration_open"`
	RegistrationStartDate *time.Time `json:"registration_start_date,omitempty"`
	RegistrationEndDate   *time.Time `json:"registration_end_date,omitempty"`

	// Metadatos
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`

	// Organización
	Organization OrganizationSummaryResponse `json:"organization"`

	// Estados computados
	IsUpcoming bool `json:"is_upcoming"`
	IsPast     bool `json:"is_past"`
	IsOngoing  bool `json:"is_ongoing"`

	// Información del usuario (si está autenticado)
	IsFavorite   bool `json:"is_favorite,omitempty"`
	IsRegistered bool `json:"is_registered,omitempty"`
	CanEdit      bool `json:"can_edit,omitempty"`
	CanManage    bool `json:"can_manage,omitempty"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CanceledAt  *time.Time `json:"canceled_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// EventDetailResponse DTO de respuesta detallada para un evento (incluye campos adicionales)
type EventDetailResponse struct {
	EventResponse

	// Información adicional no incluida en la respuesta básica
	Requirements string `json:"requirements,omitempty"`
	Agenda       string `json:"agenda,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`

	// Estadísticas adicionales (para organizadores/admin)
	Statistics *EventStatistics `json:"statistics,omitempty"`
}

// EventListResponse DTO para lista de eventos
type EventListResponse struct {
	Events     []EventResponse       `json:"events"`
	Pagination common.PaginationMeta `json:"pagination"`
	Filters    AppliedFilters        `json:"filters,omitempty"`
}

// EventSummaryResponse DTO resumido para listados
type EventSummaryResponse struct {
	ID               string    `json:"id"`
	Slug             string    `json:"slug"`
	Title            string    `json:"title"`
	ShortDesc        string    `json:"short_desc"`
	Type             string    `json:"type"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	IsOnline         bool      `json:"is_online"`
	VenueCity        string    `json:"venue_city,omitempty"`
	IsFree           bool      `json:"is_free"`
	Price            *int      `json:"price,omitempty"`
	ImageURL         string    `json:"image_url,omitempty"`
	CurrentAttendees int       `json:"current_attendees"`
	MaxAttendees     *int      `json:"max_attendees"`
	IsFeatured       bool      `json:"is_featured"`
	Organization     string    `json:"organization_name"`
	Tags             []string  `json:"tags,omitempty"`
}

// EventStatistics estadísticas del evento (para organizadores)
type EventStatistics struct {
	TotalViews         int            `json:"total_views"`
	UniqueViews        int            `json:"unique_views,omitempty"`
	TotalRegistrations int            `json:"total_registrations"`
	ConfirmedAttendees int            `json:"confirmed_attendees"`
	PendingAttendees   int            `json:"pending_attendees"`
	CanceledAttendees  int            `json:"canceled_attendees"`
	ConversionRate     float64        `json:"conversion_rate"` // Views to registration
	OccupancyRate      float64        `json:"occupancy_rate"`  // Current/Max attendees
	DailyViews         []DailyMetric  `json:"daily_views,omitempty"`
	SourceBreakdown    map[string]int `json:"source_breakdown,omitempty"`
}

// DailyMetric métrica diaria
type DailyMetric struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// EventStatusChangeResponse respuesta para cambios de estado
type EventStatusChangeResponse struct {
	ID        string    `json:"id"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	ChangedAt time.Time `json:"changed_at"`
	ChangedBy string    `json:"changed_by,omitempty"`
	Message   string    `json:"message"`
}

// EventAttendeeResponse información de asistente
type EventAttendeeResponse struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	Status           string     `json:"status"` // confirmed, pending, canceled
	RegistrationDate time.Time  `json:"registration_date"`
	AttendedAt       *time.Time `json:"attended_at,omitempty"`
	CheckedInBy      string     `json:"checked_in_by,omitempty"`
}

// EventAttendeesListResponse lista de asistentes
type EventAttendeesListResponse struct {
	EventID    string                  `json:"event_id"`
	EventTitle string                  `json:"event_title"`
	Attendees  []EventAttendeeResponse `json:"attendees"`
	Statistics AttendeeStatistics      `json:"statistics"`
	Pagination common.PaginationMeta   `json:"pagination"`
}

// AttendeeStatistics estadísticas de asistentes
type AttendeeStatistics struct {
	Total     int `json:"total"`
	Confirmed int `json:"confirmed"`
	Pending   int `json:"pending"`
	Canceled  int `json:"canceled"`
	CheckedIn int `json:"checked_in"`
}

// FavoriteEventResponse evento favorito del usuario
type FavoriteEventResponse struct {
	EventID     string    `json:"event_id"`
	EventTitle  string    `json:"event_title"`
	EventSlug   string    `json:"event_slug"`
	StartDate   time.Time `json:"start_date"`
	IsFree      bool      `json:"is_free"`
	ImageURL    string    `json:"image_url,omitempty"`
	FavoritedAt time.Time `json:"favorited_at"`
}

// UpcomingEventsResponse eventos próximos
type UpcomingEventsResponse struct {
	Today     []EventSummaryResponse `json:"today,omitempty"`
	ThisWeek  []EventSummaryResponse `json:"this_week,omitempty"`
	ThisMonth []EventSummaryResponse `json:"this_month,omitempty"`
	Later     []EventSummaryResponse `json:"later,omitempty"`
}

// NearbyEventsResponse eventos cercanos
type NearbyEventsResponse struct {
	CenterLatitude  float64             `json:"center_latitude"`
	CenterLongitude float64             `json:"center_longitude"`
	RadiusKm        int                 `json:"radius_km"`
	Events          []EventWithDistance `json:"events"`
	TotalFound      int                 `json:"total_found"`
}

// EventWithDistance evento con distancia
type EventWithDistance struct {
	EventSummaryResponse
	DistanceKm float64 `json:"distance_km"`
}
