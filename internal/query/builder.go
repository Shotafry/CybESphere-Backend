package query

import (
	"context"
	"fmt"
	"strings"

	"cybesphere-backend/internal/common"

	"gorm.io/gorm"
)

// Builder constructor de consultas genérico
type Builder[T common.BaseEntity] struct {
	db             *gorm.DB
	query          *gorm.DB
	entityType     T
	allowedFilters map[string]string // campo -> operador
	allowedSorts   []string
	defaultSort    string
	searchFields   []string
}

// NewBuilder crea un nuevo query builder
func NewBuilder[T common.BaseEntity](db *gorm.DB) *Builder[T] {
	var entity T
	return &Builder[T]{
		db:             db,
		query:          db,
		entityType:     entity,
		allowedFilters: make(map[string]string),
		allowedSorts:   []string{"created_at", "updated_at"},
		defaultSort:    "created_at",
		searchFields:   []string{},
	}
}

// SetAllowedFilters configura filtros permitidos
func (b *Builder[T]) SetAllowedFilters(filters map[string]string) *Builder[T] {
	b.allowedFilters = filters
	return b
}

// SetAllowedSorts configura campos de ordenamiento permitidos
func (b *Builder[T]) SetAllowedSorts(sorts []string) *Builder[T] {
	b.allowedSorts = sorts
	return b
}

// SetDefaultSort configura ordenamiento por defecto
func (b *Builder[T]) SetDefaultSort(sort string) *Builder[T] {
	b.defaultSort = sort
	return b
}

// SetSearchFields configura campos de búsqueda
func (b *Builder[T]) SetSearchFields(fields []string) *Builder[T] {
	b.searchFields = fields
	return b
}

// ApplyOptions aplica las opciones de consulta
func (b *Builder[T]) ApplyOptions(opts common.QueryOptions) *Builder[T] {
	b.query = b.db.Model(b.entityType)

	// Aplicar filtros
	b.applyFilters(opts)

	// Aplicar búsqueda
	if opts.Search != "" {
		b.applySearch(opts.Search)
	}

	// Aplicar preloads
	for _, preload := range opts.Preloads {
		b.query = b.query.Preload(preload)
	}

	return b
}

// applyFilters aplica filtros de la consulta
func (b *Builder[T]) applyFilters(opts common.QueryOptions) {
	for field, value := range opts.Filters {
		if operator, allowed := b.allowedFilters[field]; allowed {
			switch operator {
			case "=":
				b.query = b.query.Where(fmt.Sprintf("%s = ?", field), value)
			case "LIKE":
				b.query = b.query.Where(fmt.Sprintf("%s ILIKE ?", field), "%"+fmt.Sprintf("%v", value)+"%")
			case "IN":
				b.query = b.query.Where(fmt.Sprintf("%s IN ?", field), value)
			case ">=":
				b.query = b.query.Where(fmt.Sprintf("%s >= ?", field), value)
			case "<=":
				b.query = b.query.Where(fmt.Sprintf("%s <= ?", field), value)
			case ">":
				b.query = b.query.Where(fmt.Sprintf("%s > ?", field), value)
			case "<":
				b.query = b.query.Where(fmt.Sprintf("%s < ?", field), value)
			}
		}
	}
}

// applySearch aplica búsqueda en múltiples campos
func (b *Builder[T]) applySearch(search string) {
	if len(b.searchFields) == 0 {
		return
	}

	var conditions []string
	var args []interface{}

	for _, field := range b.searchFields {
		conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
		args = append(args, "%"+search+"%")
	}

	if len(conditions) > 0 {
		b.query = b.query.Where(strings.Join(conditions, " OR "), args...)
	}
}

// AddWhere agrega una condición WHERE personalizada
func (b *Builder[T]) AddWhere(condition string, args ...interface{}) *Builder[T] {
	b.query = b.query.Where(condition, args...)
	return b
}

// Count cuenta los registros
func (b *Builder[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := b.query.WithContext(ctx).Count(&count).Error
	return count, err
}

// Find ejecuta la consulta y retorna los resultados
func (b *Builder[T]) Find(ctx context.Context, opts common.QueryOptions) ([]*T, error) {
	// Aplicar ordenamiento
	orderField := opts.OrderBy
	if !b.isAllowedSort(orderField) {
		orderField = b.defaultSort
	}

	orderClause := fmt.Sprintf("%s %s", orderField, strings.ToUpper(opts.OrderDir))
	b.query = b.query.Order(orderClause)

	// Aplicar paginación
	b.query = b.query.Offset(opts.Offset).Limit(opts.Limit)

	var results []*T
	err := b.query.WithContext(ctx).Find(&results).Error
	return results, err
}

// FindWithPagination ejecuta la consulta con paginación
func (b *Builder[T]) FindWithPagination(ctx context.Context, opts common.QueryOptions) ([]*T, *common.PaginationMeta, error) {
	// Contar total (antes de aplicar paginación)
	total, err := b.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Obtener resultados paginados
	results, err := b.Find(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	// Crear metadatos de paginación
	pagination := common.NewPaginationMeta(opts.Page, opts.Limit, total)

	return results, pagination, nil
}

// First obtiene el primer resultado
func (b *Builder[T]) First(ctx context.Context) (*T, error) {
	var result T
	err := b.query.WithContext(ctx).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// isAllowedSort verifica si el campo de ordenamiento está permitido
func (b *Builder[T]) isAllowedSort(field string) bool {
	for _, allowed := range b.allowedSorts {
		if allowed == field {
			return true
		}
	}
	return false
}
