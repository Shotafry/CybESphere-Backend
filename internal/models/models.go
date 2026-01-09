package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AllModels slice con todos los modelos para auto-migración
var AllModels = []interface{}{
	&User{},
	&Organization{},
	&Event{},
	&RefreshToken{}, // Agregado el nuevo modelo
	&AuditLog{},
}

// AutoMigrate ejecuta la auto-migración de todos los modelos
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels...)
}

// CreateIndexes crea índices adicionales que no se pueden definir con tags
func CreateIndexes(db *gorm.DB) error {
	// Índice compuesto para búsquedas geoespaciales de usuarios
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_location 
		ON users (latitude, longitude) 
		WHERE latitude IS NOT NULL AND longitude IS NOT NULL
	`).Error; err != nil {
		return err
	}

	// Índice compuesto para búsquedas geoespaciales de organizaciones
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_organizations_location 
		ON organizations (latitude, longitude) 
		WHERE latitude IS NOT NULL AND longitude IS NOT NULL
	`).Error; err != nil {
		return err
	}

	// Índice compuesto para búsquedas geoespaciales de eventos
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_events_location 
		ON events (latitude, longitude) 
		WHERE latitude IS NOT NULL AND longitude IS NOT NULL
	`).Error; err != nil {
		return err
	}

	// Índice para eventos por fecha y estado
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_events_status_dates 
		ON events (status, start_date, end_date) 
		WHERE status = 'published'
	`).Error; err != nil {
		return err
	}

	// Índice para búsqueda full-text en eventos
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_events_fulltext 
		ON events USING gin(to_tsvector('spanish_unaccent', title || ' ' || COALESCE(description, '')))
	`).Error; err != nil {
		return err
	}

	// Índice GIN para tags de eventos
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_events_tags 
		ON events USING gin(tags)
	`).Error; err != nil {
		return err
	}

	// Índices específicos para refresh tokens
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_active 
		ON refresh_tokens (user_id, is_revoked, expires_at) 
		WHERE is_revoked = false
	`).Error; err != nil {
		return err
	}

	// Índice para cleanup de tokens expirados
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_cleanup 
		ON refresh_tokens (expires_at, is_revoked)
	`).Error; err != nil {
		return err
	}

	// Índice para auditoría por IP
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_ip_created 
		ON refresh_tokens (ip_address, created_at)
	`).Error; err != nil {
		return err
	}

	return nil
}

// SeedData crea datos de ejemplo para desarrollo
func SeedData(db *gorm.DB) error {
	// Verificar si ya hay datos
	var userCount int64
	if err := db.Model(&User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		return nil // Ya hay datos, no hacer seed
	}

	// Crear usuario admin
	adminUser := &User{
		Email:      "admin@cybesphere.com",
		Password:   "admin123456", // Se hasheará automáticamente
		FirstName:  "Administrator",
		LastName:   "CybESphere",
		Role:       RoleAdmin,
		IsActive:   true,
		IsVerified: true,
		City:       "Madrid",
		Country:    "Spain",
		Timezone:   "Europe/Madrid",
	}

	if err := db.Create(adminUser).Error; err != nil {
		return err
	}

	// Crear organización de ejemplo
	org := &Organization{
		Name:        "CyberSecurity Spain",
		Slug:        "cybersecurity-spain",
		Description: "Comunidad española de profesionales de ciberseguridad",
		Email:       "info@cybersecurityspain.com",
		Website:     "https://cybersecurityspain.com",
		City:        "Madrid",
		Country:     "Spain",
		Status:      OrgStatusActive,
		IsVerified:  true,
	}

	if err := db.Create(org).Error; err != nil {
		return err
	}

	// Crear usuario organizador
	orgIDStr := org.ID.String()
	organizerUser := &User{
		Email:          "organizer@cybersecurityspain.com",
		Password:       "organizer123456",
		FirstName:      "María",
		LastName:       "García",
		Role:           RoleOrganizer,
		IsActive:       true,
		IsVerified:     true,
		OrganizationID: &orgIDStr,
		Company:        "CyberSecurity Spain",
		Position:       "Event Manager",
		City:           "Madrid",
		Country:        "Spain",
		Timezone:       "Europe/Madrid",
	}

	if err := db.Create(organizerUser).Error; err != nil {
		return err
	}

	// Crear usuario normal
	normalUser := &User{
		Email:      "user@example.com",
		Password:   "user123456",
		FirstName:  "Juan",
		LastName:   "Pérez",
		Role:       RoleUser,
		IsActive:   true,
		IsVerified: true,
		Company:    "Tech Corp",
		Position:   "Security Analyst",
		City:       "Barcelona",
		Country:    "Spain",
		Timezone:   "Europe/Madrid",
	}

	if err := db.Create(normalUser).Error; err != nil {
		return err
	}

	// Crear evento de ejemplo
	event := &Event{
		Title:           "Introducción a la Ciberseguridad",
		Description:     "Workshop práctico sobre los fundamentos de la ciberseguridad para principiantes. Cubriremos conceptos básicos, herramientas esenciales y mejores prácticas.",
		ShortDesc:       "Workshop de ciberseguridad para principiantes",
		Type:            EventTypeWorkshop,
		Category:        "Security Basics",
		Level:           "beginner",
		OrganizationID:  org.ID.String(),
		StartDate:       db.NowFunc().Add(24 * 7 * time.Hour),           // En una semana
		EndDate:         db.NowFunc().Add(24*7*time.Hour + 4*time.Hour), // 4 horas de duración
		IsOnline:        true,
		OnlineURL:       "https://meet.google.com/abc-defg-hij",
		MaxAttendees:    func() *int { i := 50; return &i }(),
		IsFree:          true,
		Status:          EventStatusPublished,
		IsPublic:        true,
		ContactEmail:    "events@cybersecurityspain.com",
		MetaTitle:       "Workshop: Introducción a la Ciberseguridad",
		MetaDescription: "Aprende los fundamentos de la ciberseguridad en este workshop práctico para principiantes",
		Tags:            datatypes.JSON([]byte(`["cybersecurity", "workshop", "beginners", "online"]`)),
	}

	now := db.NowFunc()
	event.PublishedAt = &now

	if err := db.Create(event).Error; err != nil {
		return err
	}

	return nil
}
