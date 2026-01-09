// internal/services/organization_service.go
package services

import (
	"context"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/repositories"
	"cybesphere-backend/pkg/logger"
)

// OrganizationServiceImpl implementación concreta del servicio de organizaciones
type OrganizationServiceImpl struct {
	*BaseService[models.Organization, dto.CreateOrganizationRequest, dto.UpdateOrganizationRequest]
	orgRepo  *repositories.OrganizationRepository
	userRepo *repositories.UserRepository
	auth     AuthorizationService
}

// Verificación en tiempo de compilación
var _ OrganizationService = (*OrganizationServiceImpl)(nil)

// NewOrganizationService crea nueva instancia del servicio de organizaciones
func NewOrganizationService(
	orgRepo *repositories.OrganizationRepository,
	userRepo *repositories.UserRepository,
	mapper ResponseMapper,
	auth AuthorizationService,
) OrganizationService {
	base := NewBaseService[models.Organization, dto.CreateOrganizationRequest, dto.UpdateOrganizationRequest](
		orgRepo, mapper, auth,
	)

	return &OrganizationServiceImpl{
		BaseService: base,
		orgRepo:     orgRepo,
		userRepo:    userRepo,
		auth:        auth,
	}
}

// CreateOrganization crea una nueva organización (wrapper para compatibilidad)
func (s *OrganizationServiceImpl) CreateOrganization(ctx context.Context, req dto.CreateOrganizationRequest, userCtx *common.UserContext) (*models.Organization, error) {
	// Validaciones específicas
	if err := s.validateOrganizationCreation(userCtx); err != nil {
		return nil, err
	}

	// Crear organización usando el método base
	org, err := s.Create(ctx, req, userCtx)
	if err != nil {
		return nil, err
	}

	// Post-procesamiento: asignar usuario como organizador si no es admin
	if !userCtx.IsAdmin() {
		go func() {
			if err := s.userRepo.UpdateRole(context.Background(), userCtx.ID, models.RoleOrganizer); err != nil {
				logger.Error("Error actualizando rol del usuario: ", err)
			}
			user, err := s.userRepo.GetByID(context.Background(), userCtx.ID)
			if err != nil {
				logger.Error("Error obteniendo usuario por ID: ", err)
				return
			}
			if user != nil {
				orgIDStr := org.GetID()
				user.OrganizationID = &orgIDStr
				if err := s.userRepo.Update(context.Background(), user); err != nil {
					logger.Error("Error actualizando organización del usuario: ", err)
				}
			}
		}()
	}

	return org, nil
}

// GetActiveOrganizations obtiene organizaciones activas
func (s *OrganizationServiceImpl) GetActiveOrganizations(ctx context.Context, opts common.QueryOptions) ([]*models.Organization, *common.PaginationMeta, error) {
	return s.orgRepo.GetActive(ctx, opts)
}

// VerifyOrganization verifica una organización (solo admin)
func (s *OrganizationServiceImpl) VerifyOrganization(ctx context.Context, id string, userCtx *common.UserContext) (*models.Organization, error) {
	if !userCtx.IsAdmin() {
		return nil, common.NewBusinessError("admin_required", "Solo administradores pueden verificar organizaciones")
	}

	// Verificar organización
	if err := s.orgRepo.Verify(ctx, id, userCtx.ID); err != nil {
		return nil, err
	}

	return s.orgRepo.GetByID(ctx, id)
}

// GetMembers obtiene miembros de una organización
func (s *OrganizationServiceImpl) GetMembers(ctx context.Context, organizationID string, opts common.QueryOptions, userCtx *common.UserContext) ([]*models.User, *common.PaginationMeta, error) {
	// Verificar permisos para ver miembros
	if !userCtx.IsAdmin() && !userCtx.CanManageOrganization(organizationID) {
		return nil, nil, common.NewBusinessError("access_denied", "No tienes permisos para ver los miembros")
	}

	return s.userRepo.GetByOrganization(ctx, organizationID, opts)
}

// validateOrganizationCreation valida creación de organización
func (s *OrganizationServiceImpl) validateOrganizationCreation(userCtx *common.UserContext) error {
	// Solo usuarios verificados pueden crear organizaciones (excepto admin)
	if !userCtx.IsAdmin() && !userCtx.IsVerified {
		return common.NewBusinessError("verification_required",
			"Debes tener una cuenta verificada para crear organizaciones")
	}

	// Verificar que no existe organización con el mismo nombre
	// (Esto lo haría el repository con unique constraint, pero agregamos validación)

	return nil
}
