package common

import (
	"context"
	"time"
)

// BaseEntity interfaz para todas las entidades del dominio
type BaseEntity interface {
	GetID() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

// Repository interfaz genérica para acceso a datos
type Repository[T BaseEntity] interface {
	GetByID(ctx context.Context, id string) (*T, error)
	GetAll(ctx context.Context, opts QueryOptions) ([]*T, *PaginationMeta, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context, opts QueryOptions) (int64, error)
}

// Service interfaz genérica para lógica de negocio
type Service[T BaseEntity, CreateDTO any, UpdateDTO any] interface {
	Get(ctx context.Context, id string, userCtx *UserContext) (*T, error)
	GetAll(ctx context.Context, opts QueryOptions, userCtx *UserContext) ([]*T, *PaginationMeta, error)
	Create(ctx context.Context, dto CreateDTO, userCtx *UserContext) (*T, error)
	Update(ctx context.Context, id string, dto UpdateDTO, userCtx *UserContext) (*T, error)
	Delete(ctx context.Context, id string, userCtx *UserContext) error
}
