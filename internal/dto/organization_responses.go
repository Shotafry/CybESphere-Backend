package dto

import (
	"cybesphere-backend/internal/common"
	"time"
)

// OrganizationResponse DTO de respuesta para una organización
type OrganizationResponse struct {
	// Identificación
	ID   string `json:"id"`
	Slug string `json:"slug"`

	// Información básica
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website,omitempty"`

	// Información de contacto (solo para admins y miembros)
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`

	// Geolocalización
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	// Branding y medios
	LogoURL        string `json:"logo_url,omitempty"`
	BannerURL      string `json:"banner_url,omitempty"`
	PrimaryColor   string `json:"primary_color,omitempty"`
	SecondaryColor string `json:"secondary_color,omitempty"`

	// Redes sociales
	SocialMedia *SocialMediaLinks `json:"social_media,omitempty"`

	// Estado y verificación
	Status     string     `json:"status"`
	IsVerified bool       `json:"is_verified"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`

	// Estadísticas
	EventsCount    int `json:"events_count"`
	MembersCount   int `json:"members_count,omitempty"`
	UpcomingEvents int `json:"upcoming_events,omitempty"`

	// Configuración (solo para admins)
	MaxEvents       *int `json:"max_events,omitempty"`
	CanCreateEvents bool `json:"can_create_events,omitempty"`

	// Información del usuario (si está autenticado)
	IsMember  bool `json:"is_member,omitempty"`
	CanEdit   bool `json:"can_edit,omitempty"`
	CanManage bool `json:"can_manage,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrganizationDetailResponse DTO de respuesta detallada para una organización
type OrganizationDetailResponse struct {
	OrganizationResponse

	// Información adicional de verificación (solo admin)
	TaxID            string `json:"tax_id,omitempty"`
	LegalName        string `json:"legal_name,omitempty"`
	RegistrationDocs string `json:"registration_docs,omitempty"`
	VerifiedBy       string `json:"verified_by,omitempty"`

	// Estadísticas detalladas (para miembros)
	Statistics *OrganizationStatistics `json:"statistics,omitempty"`

	// Eventos recientes
	RecentEvents []EventSummaryResponse `json:"recent_events,omitempty"`

	// Miembros destacados (públicos)
	FeaturedMembers []UserSummaryResponse `json:"featured_members,omitempty"`
}

// OrganizationSummaryResponse DTO resumido para listados
type OrganizationSummaryResponse struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	LogoURL     string `json:"logo_url,omitempty"`
	IsVerified  bool   `json:"is_verified"`
	EventsCount int    `json:"events_count"`
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
}

// OrganizationListResponse DTO para lista de organizaciones
type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
	Pagination    common.PaginationMeta  `json:"pagination"`
	Filters       AppliedFilters         `json:"filters,omitempty"`
}

// OrganizationStatistics estadísticas de la organización
type OrganizationStatistics struct {
	TotalEvents      int            `json:"total_events"`
	PublishedEvents  int            `json:"published_events"`
	CompletedEvents  int            `json:"completed_events"`
	CanceledEvents   int            `json:"canceled_events"`
	TotalAttendees   int            `json:"total_attendees"`
	AverageAttendees float64        `json:"average_attendees"`
	EventsByType     map[string]int `json:"events_by_type"`
	EventsByMonth    []MonthlyCount `json:"events_by_month"`
	TopCities        []CityCount    `json:"top_cities"`
}

// OrganizationMemberResponse miembro de organización
type OrganizationMemberResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Position  string    `json:"position,omitempty"`
	JoinedAt  time.Time `json:"joined_at"`
	IsActive  bool      `json:"is_active"`
}

// OrganizationMembersListResponse lista de miembros
type OrganizationMembersListResponse struct {
	OrganizationID   string                       `json:"organization_id"`
	OrganizationName string                       `json:"organization_name"`
	Members          []OrganizationMemberResponse `json:"members"`
	Pagination       common.PaginationMeta        `json:"pagination"`
}

// SocialMediaLinks enlaces de redes sociales
type SocialMediaLinks struct {
	LinkedIn  string `json:"linkedin,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	Instagram string `json:"instagram,omitempty"`
	YouTube   string `json:"youtube,omitempty"`
}

// MonthlyCount conteo mensual
type MonthlyCount struct {
	Month string `json:"month"` // YYYY-MM
	Count int    `json:"count"`
}

// CityCount conteo por ciudad
type CityCount struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Count   int    `json:"count"`
}
