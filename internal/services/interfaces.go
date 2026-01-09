// internal/services/interfaces.go
package services

import (
	"context"
	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/permissions"
)

// ResponseMapper interfaz para mapeo de responses
type ResponseMapper interface {
	DTOToEntity(dto interface{}, userCtx *common.UserContext) (interface{}, error)
	ApplyUpdateDTO(entity interface{}, dto interface{}, userCtx *common.UserContext) (interface{}, error)
	EntityToResponse(entity interface{}, userCtx *common.UserContext) (interface{}, error)
}

// AuthorizationService interfaz para autorización
type AuthorizationService interface {
	CheckReadPermission(userCtx *common.UserContext, resourceType, resourceID string) error
	CheckCreatePermission(userCtx *common.UserContext, resourceType string) error
	CheckUpdatePermission(userCtx *common.UserContext, resourceType, resourceID string) error
	CheckDeletePermission(userCtx *common.UserContext, resourceType, resourceID string) error
	ApplySecurityFilters(opts *common.QueryOptions, userCtx *common.UserContext, resourceType string)
	GetUserCapabilities(userCtx *common.UserContext) map[string]bool
	GetUserPermissions(userCtx *common.UserContext) []permissions.Permission
	CanUserManageEvent(userCtx *common.UserContext, eventID string) bool
	CanUserManageOrganization(userCtx *common.UserContext, orgID string) bool
	ValidateEventCreation(userCtx *common.UserContext, organizationID string) error
	ValidateOrganizationCreation(userCtx *common.UserContext) error
}

// EventService interfaz para servicio de eventos
type EventService interface {
	// Métodos CRUD base
	Get(ctx context.Context, id string, userCtx *common.UserContext) (*models.Event, error)
	GetAll(ctx context.Context, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.Event, *common.PaginationMeta, error)
	GetEventBySlug(ctx context.Context, slug string) (*models.Event, error)
	Create(ctx context.Context, dto dto.CreateEventRequest, userCtx *common.UserContext) (*models.Event, error)
	Update(ctx context.Context, id string, dto dto.UpdateEventRequest, userCtx *common.UserContext) (*models.Event, error)
	Delete(ctx context.Context, id string, userCtx *common.UserContext) error

	// Métodos específicos de eventos
	CreateEvent(ctx context.Context, req dto.CreateEventRequest, userCtx *common.UserContext) (*models.Event, error)
	PublishEvent(ctx context.Context, id string, userCtx *common.UserContext) (*models.Event, error)
	CancelEvent(ctx context.Context, id string, userCtx *common.UserContext) (*models.Event, error)
	GetFeaturedEvents(ctx context.Context, limit int) ([]*models.Event, error)
	GetUpcomingEvents(ctx context.Context, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.Event, *common.PaginationMeta, error)
	GetEventsByOrganization(ctx context.Context, orgID string, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.Event, *common.PaginationMeta, error)
	AddToFavorites(ctx context.Context, eventID string, userCtx *common.UserContext) error
	RemoveFromFavorites(ctx context.Context, eventID string, userCtx *common.UserContext) error
}

// OrganizationService interfaz para servicio de organizaciones
type OrganizationService interface {
	// Métodos CRUD base
	Get(ctx context.Context, id string, userCtx *common.UserContext) (*models.Organization, error)
	GetAll(ctx context.Context, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.Organization, *common.PaginationMeta, error)
	Create(ctx context.Context, dto dto.CreateOrganizationRequest, userCtx *common.UserContext) (*models.Organization, error)
	Update(ctx context.Context, id string, dto dto.UpdateOrganizationRequest, userCtx *common.UserContext) (*models.Organization, error)
	Delete(ctx context.Context, id string, userCtx *common.UserContext) error

	// Métodos específicos de organizaciones
	CreateOrganization(ctx context.Context, req dto.CreateOrganizationRequest, userCtx *common.UserContext) (*models.Organization, error)
	GetActiveOrganizations(ctx context.Context, opts common.QueryOptions) ([]*models.Organization, *common.PaginationMeta, error)
	VerifyOrganization(ctx context.Context, id string, userCtx *common.UserContext) (*models.Organization, error)
	GetMembers(ctx context.Context, orgID string, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.User, *common.PaginationMeta, error)
}

// UserService interfaz para servicio de usuarios
type UserService interface {
	// Métodos CRUD base
	Get(ctx context.Context, id string, userCtx *common.UserContext) (*models.User, error)
	GetAll(ctx context.Context, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.User, *common.PaginationMeta, error)
	Create(ctx context.Context, dto dto.CreateUserRequest, userCtx *common.UserContext) (*models.User, error)
	Update(ctx context.Context, id string, dto dto.UpdateUserRequest, userCtx *common.UserContext) (*models.User, error)
	Delete(ctx context.Context, id string, userCtx *common.UserContext) error

	// Métodos específicos de usuarios
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateRole(ctx context.Context, userID string, newRole models.UserRole, userCtx *common.UserContext) error
	ActivateUser(ctx context.Context, userID string, userCtx *common.UserContext) error
	DeactivateUser(ctx context.Context, userID string, userCtx *common.UserContext) error
	GetUserSessions(ctx context.Context, userID string, userCtx *common.UserContext) ([]*models.RefreshToken, error)
	RevokeUserSession(ctx context.Context, sessionID string, userCtx *common.UserContext) error
}
