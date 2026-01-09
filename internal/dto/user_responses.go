package dto

import (
	"cybesphere-backend/internal/common"
	"time"
)

// UserResponse DTO de respuesta para un usuario
type UserResponse struct {
	// Identificación
	ID string `json:"id"`

	// Información básica
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`

	// Perfil profesional
	Company  string `json:"company,omitempty"`
	Position string `json:"position,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Website  string `json:"website,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
	Twitter  string `json:"twitter,omitempty"`

	// Sistema de roles y estado
	Role       string `json:"role"`
	IsActive   bool   `json:"is_active"`
	IsVerified bool   `json:"is_verified"`

	// Ubicación
	City      string   `json:"city,omitempty"`
	Country   string   `json:"country,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	// Configuraciones
	Timezone          string `json:"timezone"`
	Language          string `json:"language"`
	NewsletterEnabled bool   `json:"newsletter_enabled"`

	// Organización
	Organization *OrganizationSummaryResponse `json:"organization,omitempty"`

	// Estadísticas (solo visible para el propio usuario o admin)
	Statistics *UserStatistics `json:"statistics,omitempty"`

	// Información del usuario actual (permisos)
	CanEdit   bool `json:"can_edit,omitempty"`
	CanManage bool `json:"can_manage,omitempty"`

	// Timestamps
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserDetailResponse DTO de respuesta detallada para un usuario
type UserDetailResponse struct {
	UserResponse

	// Eventos favoritos
	FavoriteEvents []EventSummaryResponse `json:"favorite_events,omitempty"`

	// Eventos registrados
	RegisteredEvents []EventSummaryResponse `json:"registered_events,omitempty"`

	// Sesiones activas (solo para el propio usuario)
	ActiveSessions []SessionResponse `json:"active_sessions,omitempty"`
}

// UserSummaryResponse DTO resumido para listados
type UserSummaryResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email,omitempty"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FullName     string `json:"full_name"`
	Company      string `json:"company,omitempty"`
	Position     string `json:"position,omitempty"`
	Role         string `json:"role"`
	IsVerified   bool   `json:"is_verified"`
	Organization string `json:"organization_name,omitempty"`
}

// UserListResponse DTO para lista de usuarios
type UserListResponse struct {
	Users      []UserResponse        `json:"users"`
	Pagination common.PaginationMeta `json:"pagination"`
	Filters    AppliedFilters        `json:"filters,omitempty"`
}

// UserStatistics estadísticas del usuario
type UserStatistics struct {
	EventsAttended   int       `json:"events_attended"`
	EventsOrganized  int       `json:"events_organized,omitempty"`
	FavoriteEvents   int       `json:"favorite_events"`
	TotalConnections int       `json:"total_connections,omitempty"`
	MemberSince      time.Time `json:"member_since"`
	LastActive       time.Time `json:"last_active"`
}

// SessionResponse información de sesión
type SessionResponse struct {
	ID         string     `json:"id"`
	TokenID    string     `json:"token_id"`
	DeviceInfo string     `json:"device_info"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  time.Time  `json:"expires_at"`
	IsCurrent  bool       `json:"is_current"`
}

// UserProfileResponse perfil público del usuario
type UserProfileResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Company   string `json:"company,omitempty"`
	Position  string `json:"position,omitempty"`
	Bio       string `json:"bio,omitempty"`
	City      string `json:"city,omitempty"`
	Country   string `json:"country,omitempty"`
	LinkedIn  string `json:"linkedin,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Website   string `json:"website,omitempty"`

	// Solo si el usuario ha hecho públicos sus eventos
	PublicEvents []EventSummaryResponse `json:"public_events,omitempty"`
}

// UserCapabilitiesResponse capacidades del usuario
type UserCapabilitiesResponse struct {
	UserID       string                       `json:"user_id"`
	Role         string                       `json:"role"`
	Permissions  []PermissionResponse         `json:"permissions"`
	Capabilities map[string]bool              `json:"capabilities"`
	Organization *OrganizationSummaryResponse `json:"organization,omitempty"`
}

// PermissionResponse permiso individual
type PermissionResponse struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Display  string `json:"display"`
}
