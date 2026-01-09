package main

import (
	"flag"
	"fmt"
	"os"

	"cybesphere-backend/internal/config"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/internal/seeders"
	"cybesphere-backend/pkg/database"
	"cybesphere-backend/pkg/logger"

	"gorm.io/gorm"
)

func main() {
	// Configurar flags de comandos
	var (
		seederName = flag.String("seeder", "", "Nombre del seeder específico a ejecutar")
		priority   = flag.Int("priority", -1, "Ejecutar seeders hasta esta prioridad (inclusive)")
		force      = flag.Bool("force", false, "Forzar ejecución incluso si ya hay datos")
		fresh      = flag.Bool("fresh", false, "Recrear todas las tablas antes de sembrar")
		list       = flag.Bool("list", false, "Listar todos los seeders disponibles")
		help       = flag.Bool("help", false, "Mostrar ayuda")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error cargando configuración: %v\n", err)
		os.Exit(1)
	}

	// 2. Inicializar logger
	if err := logger.Init(&cfg.Logging); err != nil {
		fmt.Printf("Error inicializando logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("🌱 Iniciando herramienta de seeders de CybESphere")

	// 3. Conectar a la base de datos
	if err := database.Connect(&cfg.Database); err != nil {
		logger.Fatalf("Error conectando a la base de datos: %v", err)
	}
	defer database.Close()

	db := database.GetDB()

	// 4. Recrear tablas si se especifica fresh
	if *fresh {
		if err := recreateAllTables(db); err != nil {
			logger.Fatalf("Error recreando tablas: %v", err)
		}
		logger.Info("✅ Tablas recreadas exitosamente")
	}

	// 5. Configurar manager de seeders
	sm := seeders.NewSeederManager()
	sm.RegisterDefaultSeeders()

	// 6. Listar seeders si se solicita
	if *list {
		sm.ListSeeders(db)
		return
	}

	// 7. Ejecutar seeders según parámetros
	if err := executeSeeders(sm, db, *seederName, *priority, *force); err != nil {
		logger.Fatalf("Error ejecutando seeders: %v", err)
	}

	// 8. Mostrar resumen
	sm.PrintSummary()

	// 9. Mostrar estadísticas finales
	printDatabaseStats(db)

	logger.Info("🎉 Proceso de seeders completado exitosamente")
}

// executeSeeders ejecuta los seeders según los parámetros especificados
func executeSeeders(sm *seeders.SeederManager, db *gorm.DB, seederName string, priority int, force bool) error {
	if seederName != "" {
		// Ejecutar seeder específico
		return sm.RunSpecific(db, seederName, force)
	} else if priority >= 0 {
		// Ejecutar seeders hasta prioridad específica
		return sm.RunByPriority(db, priority, force)
	} else {
		// Ejecutar todos los seeders
		return sm.RunAll(db, force)
	}
}

// recreateAllTables elimina y recrea todas las tablas
func recreateAllTables(db *gorm.DB) error {
	logger.Info("🔄 Recreando todas las tablas...")

	// Lista de modelos a recrear (orden importante para relaciones)
	models := []interface{}{
		&models.RefreshToken{}, // Primero las tablas dependientes
		&models.Event{},
		&models.User{},
		&models.Organization{},
		&models.AuditLog{},
	}

	// Eliminar tablas en orden reverso
	for i := len(models) - 1; i >= 0; i-- {
		if err := db.Migrator().DropTable(models[i]); err != nil {
			logger.Warnf("Error eliminando tabla: %v", err)
		}
	}

	// Eliminar tabla de relaciones many-to-many
	if err := db.Exec("DROP TABLE IF EXISTS user_favorite_events").Error; err != nil {
		logger.Warnf("Error eliminando tabla user_favorite_events: %v", err)
	}

	// Recrear tablas
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("error migrando modelo %T: %w", model, err)
		}
	}

	return nil
}

// printDatabaseStats muestra estadísticas de la base de datos después del seeding
func printDatabaseStats(db *gorm.DB) {
	logger.Info("\n📊 Estadísticas de la base de datos:")

	stats := []struct {
		table string
		model interface{}
	}{
		{"Usuarios", &models.User{}},
		{"Organizaciones", &models.Organization{}},
		{"Eventos", &models.Event{}},
		{"Refresh Tokens", &models.RefreshToken{}},
	}

	for _, stat := range stats {
		var count int64
		if err := db.Model(stat.model).Count(&count).Error; err != nil {
			logger.Warnf("Error contando %s: %v", stat.table, err)
			continue
		}
		logger.Infof("  📋 %s: %d registros", stat.table, count)
	}

	// Estadísticas adicionales
	var favCount int64
	if err := db.Table("user_favorite_events").Count(&favCount).Error; err == nil {
		logger.Infof("  ❤️  Favoritos: %d relaciones", favCount)
	}

	// Eventos por estado
	eventStats := []struct {
		status string
		count  int64
	}{}

	statuses := []models.EventStatus{
		models.EventStatusDraft,
		models.EventStatusPublished,
		models.EventStatusCanceled,
		models.EventStatusCompleted,
	}

	for _, status := range statuses {
		var count int64
		db.Model(&models.Event{}).Where("status = ?", status).Count(&count)
		eventStats = append(eventStats, struct {
			status string
			count  int64
		}{string(status), count})
	}

	logger.Info("  📅 Eventos por estado:")
	for _, stat := range eventStats {
		if stat.count > 0 {
			logger.Infof("    - %s: %d", stat.status, stat.count)
		}
	}

	// Usuarios por rol
	var adminCount, organizerCount, userCount int64
	db.Model(&models.User{}).Where("role = ?", models.RoleAdmin).Count(&adminCount)
	db.Model(&models.User{}).Where("role = ?", models.RoleOrganizer).Count(&organizerCount)
	db.Model(&models.User{}).Where("role = ?", models.RoleUser).Count(&userCount)

	logger.Info("  👥 Usuarios por rol:")
	logger.Infof("    - Administradores: %d", adminCount)
	logger.Infof("    - Organizadores: %d", organizerCount)
	logger.Infof("    - Usuarios: %d", userCount)
}

// printHelp muestra la ayuda del comando
func printHelp() {
	fmt.Println(`
🌱 CybESphere Database Seeder

Este comando permite poblar la base de datos con datos de ejemplo para desarrollo y testing.

USO:
  go run cmd/seed/main.go [opciones]

OPCIONES:
  -seeder string    Ejecutar un seeder específico por nombre
  -priority int     Ejecutar seeders hasta esta prioridad (inclusive)
  -force           Forzar ejecución incluso si ya hay datos
  -fresh           Recrear todas las tablas antes de sembrar
  -list            Listar todos los seeders disponibles
  -help            Mostrar esta ayuda

EJEMPLOS:
  # Ejecutar todos los seeders
  go run cmd/seed/main.go

  # Recrear tablas y ejecutar todos los seeders
  go run cmd/seed/main.go -fresh

  # Forzar ejecución de todos los seeders
  go run cmd/seed/main.go -force

  # Ejecutar solo el seeder de usuarios
  go run cmd/seed/main.go -seeder UserSeeder

  # Ejecutar seeders hasta prioridad 2 (usuarios y organizaciones)
  go run cmd/seed/main.go -priority 2

  # Listar seeders disponibles
  go run cmd/seed/main.go -list

SEEDERS DISPONIBLES:
  1. UserSeeder          - Crea usuarios de ejemplo (admin, organizadores, usuarios)
  2. OrganizationSeeder  - Crea organizaciones y asigna organizadores
  3. EventSeeder         - Crea eventos en diferentes estados y fechas
  4. DemoDataSeeder      - Crea relaciones, favoritos y datos adicionales

ORDEN DE EJECUCIÓN:
Los seeders se ejecutan en orden de prioridad:
  1. Usuarios (base para todo)
  2. Organizaciones (necesitan usuarios organizadores)
  3. Eventos (necesitan organizaciones)
  4. Datos demo (necesitan usuarios y eventos)

NOTAS:
- Por defecto, los seeders no se ejecutan si ya hay datos
- Usa -force para ejecutar de todos modos
- Usa -fresh para empezar desde cero (¡CUIDADO: elimina todos los datos!)
- Los datos creados son realistas pero ficticios
- Incluye usuarios de testing con diferentes roles y estados`)
}
