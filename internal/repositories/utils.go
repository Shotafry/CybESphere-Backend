package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Transaction ejecuta una función dentro de una transacción
func Transaction[T any](ctx context.Context, db *gorm.DB, fn func(*gorm.DB) (T, error)) (T, error) {
	var result T

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = fn(tx)
		return err
	})

	return result, err
}

// BatchCreate crea múltiples entidades en batch
func BatchCreate[T any](ctx context.Context, db *gorm.DB, entities []T, batchSize int) error {
	if len(entities) == 0 {
		return nil
	}

	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize
		if end > len(entities) {
			end = len(entities)
		}

		batch := entities[i:end]
		if err := db.WithContext(ctx).CreateInBatches(batch, batchSize).Error; err != nil {
			return fmt.Errorf("failed to create batch %d-%d: %w", i, end-1, err)
		}
	}

	return nil
}

// BulkUpdate actualiza múltiples registros
func BulkUpdate[T any](ctx context.Context, db *gorm.DB, ids []string, updates map[string]interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	var model T
	err := db.WithContext(ctx).Model(&model).
		Where("id IN ?", ids).
		Updates(updates).Error

	return err
}
