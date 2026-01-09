package seeders

import (
	"fmt"
	"sort"
	"time"

	"cybesphere-backend/pkg/logger"

	"gorm.io/gorm"
)

// Seeder interface que deben implementar todos los seeders
type Seeder interface {
	Seed(db *gorm.DB) error
	Name() string
	Description() string
	Priority() int           // Orden de ejecución (menor = primero)
	CanRun(db *gorm.DB) bool // Verificar si puede ejecutarse
}

// SeederResult representa el resultado de ejecutar un seeder
type SeederResult struct {
	Name           string
	Success        bool
	Error          error
	RecordsCreated int
	Duration       time.Duration
}

// SeederManager maneja la ejecución de todos los seeders
type SeederManager struct {
	seeders []Seeder
	results []SeederResult
}

// NewSeederManager crea una nueva instancia del manager
func NewSeederManager() *SeederManager {
	return &SeederManager{
		seeders: make([]Seeder, 0),
		results: make([]SeederResult, 0),
	}
}

// RegisterSeeder registra un seeder en el manager
func (sm *SeederManager) RegisterSeeder(seeder Seeder) {
	sm.seeders = append(sm.seeders, seeder)
	logger.Debugf("Seeder registrado: %s (prioridad: %d)", seeder.Name(), seeder.Priority())
}

// RunAll ejecuta todos los seeders registrados en orden de prioridad
func (sm *SeederManager) RunAll(db *gorm.DB, force bool) error {
	logger.Info("🌱 Iniciando proceso de seeders...")

	// Ordenar seeders por prioridad
	sm.sortSeedersByPriority()

	successCount := 0
	skipCount := 0

	for _, seeder := range sm.seeders {
		// Verificar si puede ejecutarse
		if !force && !seeder.CanRun(db) {
			logger.Infof("⏭️  Saltando seeder %s - ya tiene datos", seeder.Name())
			skipCount++
			continue
		}

		if err := sm.runSingleSeeder(db, seeder); err != nil {
			return fmt.Errorf("error ejecutando seeder %s: %w", seeder.Name(), err)
		}

		successCount++
	}

	logger.Infof("✅ Proceso de seeders completado: %d ejecutados, %d saltados", successCount, skipCount)
	return nil
}

// RunSpecific ejecuta un seeder específico por nombre
func (sm *SeederManager) RunSpecific(db *gorm.DB, name string, force bool) error {
	for _, seeder := range sm.seeders {
		if seeder.Name() == name {
			logger.Infof("🎯 Ejecutando seeder específico: %s", name)

			if !force && !seeder.CanRun(db) {
				return fmt.Errorf("seeder %s no puede ejecutarse - ya tiene datos (usa --force para ejecutar de todos modos)", name)
			}

			return sm.runSingleSeeder(db, seeder)
		}
	}
	return fmt.Errorf("seeder %s no encontrado", name)
}

// RunByPriority ejecuta seeders hasta una prioridad específica
func (sm *SeederManager) RunByPriority(db *gorm.DB, maxPriority int, force bool) error {
	logger.Infof("🌱 Ejecutando seeders hasta prioridad %d...", maxPriority)

	sm.sortSeedersByPriority()

	for _, seeder := range sm.seeders {
		if seeder.Priority() > maxPriority {
			break
		}

		if !force && !seeder.CanRun(db) {
			logger.Infof("⏭️  Saltando seeder %s", seeder.Name())
			continue
		}

		if err := sm.runSingleSeeder(db, seeder); err != nil {
			return err
		}
	}

	return nil
}

// ListSeeders lista todos los seeders registrados
func (sm *SeederManager) ListSeeders(db *gorm.DB) {
	logger.Info("📋 Seeders registrados:")

	sm.sortSeedersByPriority()

	for _, seeder := range sm.seeders {
		status := "✅ Puede ejecutarse"
		if !seeder.CanRun(db) {
			status = "⏭️  Ya tiene datos"
		}

		logger.Infof("  %d. %s - %s (%s)",
			seeder.Priority(),
			seeder.Name(),
			seeder.Description(),
			status,
		)
	}
}

// GetResults retorna los resultados de la última ejecución
func (sm *SeederManager) GetResults() []SeederResult {
	return sm.results
}

// runSingleSeeder ejecuta un seeder individual
func (sm *SeederManager) runSingleSeeder(db *gorm.DB, seeder Seeder) error {
	logger.Infof("🌱 Ejecutando: %s", seeder.Name())
	logger.Debugf("   Descripción: %s", seeder.Description())

	// Contar registros antes
	var beforeCount int64
	tableName := getTableNameFromSeeder(seeder)
	if tableName != "" {
		db.Table(tableName).Count(&beforeCount)
	}

	// Ejecutar seeder
	start := time.Now()
	err := seeder.Seed(db)
	duration := time.Since(start)

	// Contar registros después
	var afterCount int64
	if tableName != "" {
		db.Table(tableName).Count(&afterCount)
	}
	recordsCreated := int(afterCount - beforeCount)

	// Registrar resultado
	result := SeederResult{
		Name:           seeder.Name(),
		Success:        err == nil,
		Error:          err,
		RecordsCreated: recordsCreated,
		Duration:       duration,
	}
	sm.results = append(sm.results, result)

	if err != nil {
		logger.Errorf("❌ Error en seeder %s: %v", seeder.Name(), err)
		return err
	}

	logger.Infof("✅ Seeder %s completado (%d registros, %s)",
		seeder.Name(), recordsCreated, duration.String())

	return nil
}

// sortSeedersByPriority ordena los seeders por prioridad usando sort.Slice
func (sm *SeederManager) sortSeedersByPriority() {
	sort.Slice(sm.seeders, func(i, j int) bool {
		return sm.seeders[i].Priority() < sm.seeders[j].Priority()
	})
}

// getTableNameFromSeeder intenta inferir el nombre de tabla del seeder
func getTableNameFromSeeder(seeder Seeder) string {
	name := seeder.Name()
	switch name {
	case "UserSeeder":
		return "users"
	case "OrganizationSeeder":
		return "organizations"
	case "EventSeeder":
		return "events"
	default:
		return ""
	}
}

// RegisterDefaultSeeders registra todos los seeders por defecto
func (sm *SeederManager) RegisterDefaultSeeders() {
	sm.RegisterSeeder(NewUserSeeder())
	sm.RegisterSeeder(NewOrganizationSeeder())
	sm.RegisterSeeder(NewEventSeeder())
	sm.RegisterSeeder(NewDemoDataSeeder())

	logger.Info("🔧 Seeders por defecto registrados")
}

// PrintSummary imprime un resumen de los resultados
func (sm *SeederManager) PrintSummary() {
	if len(sm.results) == 0 {
		return
	}

	logger.Info("\n📊 Resumen de seeders ejecutados:")

	totalRecords := 0
	successCount := 0

	for _, result := range sm.results {
		status := "✅"
		if !result.Success {
			status = "❌"
		} else {
			successCount++
		}

		logger.Infof("  %s %s: %d registros (%s)",
			status, result.Name, result.RecordsCreated, result.Duration.String())

		if result.Error != nil {
			logger.Errorf("     Error: %v", result.Error)
		}

		totalRecords += result.RecordsCreated
	}

	logger.Infof("\n🎯 Total: %d seeders ejecutados, %d registros creados",
		successCount, totalRecords)
}
