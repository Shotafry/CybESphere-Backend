package services

import (
	"context"

	"cybesphere-backend/internal/common"
)

// BaseService servicio base genérico
type BaseService[T common.BaseEntity, CreateDTO any, UpdateDTO any] struct {
	repo   common.Repository[T]
	mapper ResponseMapper
	auth   AuthorizationService
}

// NewBaseService crea una nueva instancia base
func NewBaseService[T common.BaseEntity, CreateDTO any, UpdateDTO any](
	repo common.Repository[T],
	mapper ResponseMapper,
	auth AuthorizationService,
) *BaseService[T, CreateDTO, UpdateDTO] {
	return &BaseService[T, CreateDTO, UpdateDTO]{
		repo:   repo,
		mapper: mapper,
		auth:   auth,
	}
}

// Get obtiene una entidad por ID con verificación de permisos
func (s *BaseService[T, CreateDTO, UpdateDTO]) Get(
	ctx context.Context,
	id string,
	userCtx *common.UserContext,
) (*T, error) {
	// Verificar permisos de lectura
	if err := s.auth.CheckReadPermission(userCtx, getResourceType[T](), id); err != nil {
		return nil, err
	}

	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// GetAll obtiene entidades con filtros y permisos
func (s *BaseService[T, CreateDTO, UpdateDTO]) GetAll(
	ctx context.Context,
	opts common.QueryOptions,
	userCtx *common.UserContext,
) ([]*T, *common.PaginationMeta, error) {
	// Aplicar filtros de seguridad según el rol del usuario
	s.auth.ApplySecurityFilters(&opts, userCtx, getResourceType[T]())

	return s.repo.GetAll(ctx, opts)
}

// Create crea una nueva entidad
func (s *BaseService[T, CreateDTO, UpdateDTO]) Create(
	ctx context.Context,
	dto CreateDTO,
	userCtx *common.UserContext,
) (*T, error) {
	// Verificar permisos de creación
	if err := s.auth.CheckCreatePermission(userCtx, getResourceType[T]()); err != nil {
		return nil, err
	}

	// Mapear DTO a entidad
	entity, err := s.mapper.DTOToEntity(dto, userCtx)
	if err != nil {
		return nil, err
	}

	typedEntity := entity.(T)

	// Crear en base de datos
	if err := s.repo.Create(ctx, &typedEntity); err != nil {
		return nil, err
	}

	return &typedEntity, nil
}

// Update actualiza una entidad existente
func (s *BaseService[T, CreateDTO, UpdateDTO]) Update(
	ctx context.Context,
	id string,
	dto UpdateDTO,
	userCtx *common.UserContext,
) (*T, error) {
	// Verificar que existe y verificar permisos
	existing, err := s.Get(ctx, id, userCtx)
	if err != nil {
		return nil, err
	}

	// Verificar permisos de actualización
	if err := s.auth.CheckUpdatePermission(userCtx, getResourceType[T](), id); err != nil {
		return nil, err
	}

	// Aplicar cambios del DTO
	updatedEntity, err := s.mapper.ApplyUpdateDTO(*existing, dto, userCtx)
	if err != nil {
		return nil, err
	}

	typedUpdated := updatedEntity.(T)

	// Actualizar en base de datos
	if err := s.repo.Update(ctx, &typedUpdated); err != nil {
		return nil, err
	}

	return &typedUpdated, nil
}

// Delete elimina una entidad
func (s *BaseService[T, CreateDTO, UpdateDTO]) Delete(
	ctx context.Context,
	id string,
	userCtx *common.UserContext,
) error {
	// Verificar que existe y verificar permisos
	_, err := s.Get(ctx, id, userCtx)
	if err != nil {
		return err
	}

	// Verificar permisos de eliminación
	if err := s.auth.CheckDeletePermission(userCtx, getResourceType[T](), id); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
