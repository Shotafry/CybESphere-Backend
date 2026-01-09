package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"cybesphere-backend/internal/config"
)

// DB instancia global de la base de datos
var DB *gorm.DB

// Connect establece la conexión con PostgreSQL
func Connect(cfg *config.DatabaseConfig) error {
	// Configurar logger de GORM basado en la configuración de logging
	var gormLogger logger.Interface

	// Leer LOG_LEVEL del ambiente para determinar el nivel de GORM
	logLevel := os.Getenv("LOG_LEVEL")
	serverMode := os.Getenv("SERVER_MODE")

	// Si está en modo release O el log level es warn/error/fatal, silenciar GORM
	if serverMode == "release" || logLevel == "warn" || logLevel == "error" || logLevel == "fatal" || logLevel == "panic" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else if logLevel == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		// Para "info" y otros, usar Warn (menos verboso que Info)
		gormLogger = logger.Default.LogMode(logger.Warn)
	}

	// Configuración de GORM
	gormConfig := &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			// Usar timezone de Madrid
			loc, _ := time.LoadLocation("Europe/Madrid")
			return time.Now().In(loc)
		},
	}

	// Conectar a PostgreSQL
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configurar connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configurar pool de conexiones
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verificar conexión
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close cierra la conexión a la base de datos
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}

// GetDB retorna la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}

// Ping verifica que la conexión esté activa
func Ping() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Ping()
}

// GetStats retorna estadísticas de la conexión
func GetStats() (map[string]interface{}, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	stats := sqlDB.Stats()

	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}

// AutoMigrate ejecuta las migraciones automáticas
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	return DB.AutoMigrate(models...)
}

// IsConnected verifica si hay conexión activa
func IsConnected() bool {
	if DB == nil {
		return false
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}
