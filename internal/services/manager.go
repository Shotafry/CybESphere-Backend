// internal/services/manager.go
package services

import (
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/repositories"
)

// ServiceManager centraliza todos los servicios
type ServiceManager struct {
	Events        EventService
	Organizations OrganizationService
	Users         UserService
	mapper        ResponseMapper
	auth          AuthorizationService
}

// NewServiceManager crea nueva instancia del manager con interfaces
func NewServiceManager(
	repoManager *repositories.RepositoryManager,
	mapper ResponseMapper,
	auth AuthorizationService,
) *ServiceManager {
	// Los constructores ahora devuelven interfaces directamente
	return &ServiceManager{
		Events: NewEventService(
			repoManager.Events,
			repoManager.Organizations,
			repoManager.Users,
			mapper,
			auth,
		),
		Organizations: NewOrganizationService(
			repoManager.Organizations,
			repoManager.Users,
			mapper,
			auth,
		),
		Users: NewUserService(
			repoManager.Users,
			repoManager.RefreshTokens,
			mapper,
			auth,
		),
		mapper: mapper,
		auth:   auth,
	}
}

// GetEventService retorna el servicio de eventos
func (sm *ServiceManager) GetEventService() EventService {
	return sm.Events
}

// GetOrganizationService retorna el servicio de organizaciones
func (sm *ServiceManager) GetOrganizationService() OrganizationService {
	return sm.Organizations
}

// GetUserService retorna el servicio de usuarios
func (sm *ServiceManager) GetUserService() UserService {
	return sm.Users
}

// GetAuthorizationService retorna el servicio de autorización
func (sm *ServiceManager) GetAuthorizationService() AuthorizationService {
	return sm.auth
}

// GetMapper retorna el mapper
func (sm *ServiceManager) GetMapper() ResponseMapper {
	return sm.mapper
}

// getResourceType helper para obtener tipo de recurso genérico
func getResourceType[T any]() string {
	var entity T
	switch any(entity).(type) {
	case models.Event, *models.Event:
		return "event"
	case models.Organization, *models.Organization:
		return "organization"
	case models.User, *models.User:
		return "user"
	default:
		return "unknown"
	}
}
