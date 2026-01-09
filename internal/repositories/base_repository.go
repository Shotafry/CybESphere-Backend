package repositories

import (
	"context"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/query"
	"cybesphere-backend/pkg/database"

	"gorm.io/gorm"
)

// BaseRepository implementación base para repositorios
type BaseRepository[T common.BaseEntity] struct {
	db      *gorm.DB
	builder *query.Builder[T]
}

// NewBaseRepository crea una nueva instancia base
func NewBaseRepository[T common.BaseEntity]() *BaseRepository[T] {
	db := database.GetDB()
	return &BaseRepository[T]{
		db:      db,
		builder: query.NewBuilder[T](db),
	}
}

// GetByID obtiene una entidad por ID
func (r *BaseRepository[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, common.MapGormError(err)
	}
	return &entity, nil
}

// GetAll obtiene todas las entidades con opciones
func (r *BaseRepository[T]) GetAll(ctx context.Context, opts common.QueryOptions) ([]*T, *common.PaginationMeta, error) {
	r.builder.ApplyOptions(opts)
	return r.builder.FindWithPagination(ctx, opts)
}

// Create crea una nueva entidad
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Create(entity).Error
	return common.MapGormError(err)
}

// Update actualiza una entidad existente
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Save(entity).Error
	return common.MapGormError(err)
}

// Delete elimina una entidad (soft delete)
func (r *BaseRepository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	err := r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
	return common.MapGormError(err)
}

// Count cuenta entidades con opciones
func (r *BaseRepository[T]) Count(ctx context.Context, opts common.QueryOptions) (int64, error) {
	r.builder.ApplyOptions(opts)
	return r.builder.Count(ctx)
}
