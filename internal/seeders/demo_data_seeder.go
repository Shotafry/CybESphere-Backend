package seeders

import (
	"fmt"
	"time"

	"cybesphere-backend/internal/models"
	"cybesphere-backend/pkg/logger"
	"cybesphere-backend/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DemoDataSeeder struct{}

func NewDemoDataSeeder() *DemoDataSeeder {
	return &DemoDataSeeder{}
}

func (ds *DemoDataSeeder) Name() string {
	return "DemoDataSeeder"
}

func (ds *DemoDataSeeder) Description() string {
	return "Crea datos adicionales como favoritos, relaciones y estadísticas"
}

func (ds *DemoDataSeeder) Priority() int {
	return 4 // Ejecutar al final
}

func (ds *DemoDataSeeder) CanRun(db *gorm.DB) bool {
	// Verificar si hay usuarios y eventos para crear relaciones
	var userCount, eventCount int64
	db.Model(&models.User{}).Count(&userCount)
	db.Model(&models.Event{}).Count(&eventCount)

	// Solo ejecutar si hay datos base y no hay favoritos ya creados
	if userCount == 0 || eventCount == 0 {
		return false
	}

	// Verificar si ya hay favoritos (indicaría que el seeder ya se ejecutó)
	var favCount int64
	db.Table("user_favorite_events").Count(&favCount)
	return favCount == 0
}

func (ds *DemoDataSeeder) Seed(db *gorm.DB) error {
	logger.Info("🎯 Creando datos demo adicionales...")

	// 1. Crear relaciones de favoritos
	if err := ds.createFavoriteRelations(db); err != nil {
		return err
	}

	// 2. Actualizar contadores de eventos en organizaciones
	if err := ds.updateEventCounters(db); err != nil {
		return err
	}

	// 3. Simular algunas visualizaciones de eventos
	if err := ds.simulateEventViews(db); err != nil {
		return err
	}

	// 4. Crear algunos refresh tokens de ejemplo (para testing de sesiones)
	if err := ds.createSampleRefreshTokens(db); err != nil {
		return err
	}

	logger.Info("✅ Datos demo adicionales creados exitosamente")
	return nil
}

// createFavoriteRelations crea relaciones de favoritos entre usuarios y eventos
func (ds *DemoDataSeeder) createFavoriteRelations(db *gorm.DB) error {
	logger.Debug("Creando relaciones de favoritos...")

	// Obtener usuarios activos
	var users []models.User
	if err := db.Where("is_active = ? AND role != ?", true, models.RoleAdmin).Find(&users).Error; err != nil {
		return err
	}

	// Obtener eventos publicados
	var events []models.Event
	if err := db.Where("status = ? AND is_public = ?", models.EventStatusPublished, true).Find(&events).Error; err != nil {
		return err
	}

	if len(users) == 0 || len(events) == 0 {
		logger.Debug("No hay usuarios o eventos suficientes para crear favoritos")
		return nil
	}

	// Crear favoritos para cada usuario (entre 0 y 5 eventos favoritos por usuario)
	for _, user := range users {
		numFavorites := utils.SecureRandInt(6) // 0 a 5 favoritos

		if numFavorites == 0 {
			continue // Este usuario no tiene favoritos
		}

		// Seleccionar eventos aleatorios
		selectedEvents := make([]models.Event, 0, numFavorites)
		usedIndices := make(map[int]bool)

		for len(selectedEvents) < numFavorites && len(selectedEvents) < len(events) {
			index := utils.SecureRandInt(len(events))
			if !usedIndices[index] {
				selectedEvents = append(selectedEvents, events[index])
				usedIndices[index] = true
			}
		}

		// Crear relaciones de favoritos
		if len(selectedEvents) > 0 {
			if err := db.Model(&user).Association("FavoriteEvents").Append(selectedEvents); err != nil {
				logger.Warnf("Error creando favoritos para usuario %s: %v", user.Email, err)
			} else {
				logger.Debugf("Creados %d favoritos para usuario %s", len(selectedEvents), user.Email)
			}
		}
	}

	return nil
}

// updateEventCounters actualiza los contadores de eventos en las organizaciones
func (ds *DemoDataSeeder) updateEventCounters(db *gorm.DB) error {
	logger.Debug("Actualizando contadores de eventos...")

	// Obtener todas las organizaciones
	var organizations []models.Organization
	if err := db.Find(&organizations).Error; err != nil {
		return err
	}

	// Actualizar contador para cada organización
	for _, org := range organizations {
		var eventCount int64
		if err := db.Model(&models.Event{}).Where("organization_id = ?", org.ID).Count(&eventCount).Error; err != nil {
			logger.Warnf("Error contando eventos para organización %s: %v", org.Name, err)
			continue
		}

		// Actualizar contador
		org.EventsCount = int(eventCount)
		if err := db.Save(&org).Error; err != nil {
			logger.Warnf("Error actualizando contador para organización %s: %v", org.Name, err)
		} else {
			logger.Debugf("Organización %s: %d eventos", org.Name, eventCount)
		}
	}

	return nil
}

// simulateEventViews simula visualizaciones adicionales para eventos
func (ds *DemoDataSeeder) simulateEventViews(db *gorm.DB) error {
	logger.Debug("Simulando visualizaciones de eventos...")

	// Obtener eventos que necesitan más vistas
	var events []models.Event
	if err := db.Where("views_count < ?", 50).Find(&events).Error; err != nil {
		return err
	}

	for _, event := range events {
		// Simular vistas basadas en el tipo de evento y tiempo transcurrido
		baseViews := 10

		switch event.Type {
		case models.EventTypeConference:
			baseViews = 100
		case models.EventTypeTraining, models.EventTypeWorkshop:
			baseViews = 50
		case models.EventTypeWebinar:
			baseViews = 75
		case models.EventTypeMeetup:
			baseViews = 25
		case models.EventTypeCompetition:
			baseViews = 80
		}

		// Añadir factor de tiempo (eventos más antiguos tienen más vistas)
		daysSinceCreated := int(time.Since(event.CreatedAt).Hours() / 24)
		timeBonus := daysSinceCreated * 2

		// Factor aleatorio
		randomBonus := utils.SecureRandInt(50)
		// Si el evento está destacado, más vistas
		if event.IsFeatured {
			baseViews *= 3
		}

		totalViews := baseViews + timeBonus + randomBonus

		// Actualizar vistas
		event.ViewsCount = totalViews
		if err := db.Model(&event).UpdateColumn("views_count", totalViews).Error; err != nil {
			logger.Warnf("Error actualizando vistas para evento %s: %v", event.Title, err)
		}
	}

	return nil
}

// createSampleRefreshTokens crea algunos refresh tokens de ejemplo para testing
func (ds *DemoDataSeeder) createSampleRefreshTokens(db *gorm.DB) error {
	logger.Debug("Creando refresh tokens de ejemplo...")

	// Obtener algunos usuarios para crear tokens
	var users []models.User
	if err := db.Where("is_active = ?", true).Limit(5).Find(&users).Error; err != nil {
		return err
	}

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
	}

	deviceInfos := []string{
		"Desktop - Chrome", "Desktop - macOS", "Desktop - Linux", "Mobile - iPhone", "Tablet - iPad",
	}

	ipAddresses := []string{
		"192.168.1.10", "10.0.0.15", "172.16.0.5", "192.168.100.25", "10.10.10.10",
	}

	for i, user := range users {
		// Crear 1-3 tokens por usuario para simular múltiples sesiones
		numTokens := utils.SecureRandInt(3) + 1

		for j := 0; j < numTokens; j++ {
			tokenIndex := (i + j) % len(userAgents)

			// Generar UUIDs reales para TokenHash y TokenID que respeten las limitaciones
			tokenHash := fmt.Sprintf("hash_%s", uuid.New().String()[:8]) // 13 caracteres
			tokenID := uuid.New().String()                               // 36 caracteres exactos

			refreshToken := &models.RefreshToken{
				UserID:     user.ID.String(),
				TokenHash:  tokenHash,
				TokenID:    tokenID,
				ExpiresAt:  time.Now().AddDate(0, 1, 0),     // Expira en 1 mes
				IsRevoked:  utils.SecureRandFloat32() < 0.1, // 10% revocados
				UserAgent:  userAgents[tokenIndex],
				IPAddress:  ipAddresses[tokenIndex],
				DeviceInfo: deviceInfos[tokenIndex],
			}

			// Algunos tokens con último uso reciente
			if utils.SecureRandFloat32() > 0.3 {
				lastUsed := time.Now().AddDate(0, 0, -utils.SecureRandInt(7)) // Usado en los últimos 7 días
				refreshToken.LastUsedAt = &lastUsed
			}

			// Si está revocado, configurar fecha de revocación
			if refreshToken.IsRevoked {
				revoked := time.Now().AddDate(0, 0, -utils.SecureRandInt(30)) // Revocado en los últimos 30 días
				refreshToken.RevokedAt = &revoked
			}

			if err := db.Create(refreshToken).Error; err != nil {
				logger.Warnf("Error creando refresh token para usuario %s: %v", user.Email, err)
			}
		}
	}

	return nil
}
