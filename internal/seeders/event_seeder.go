package seeders

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"cybesphere-backend/internal/models"
	"cybesphere-backend/pkg/logger"
	"cybesphere-backend/pkg/utils"

	"gorm.io/gorm"
)

type EventSeeder struct{}

func NewEventSeeder() *EventSeeder {
	return &EventSeeder{}
}

func (es *EventSeeder) Name() string {
	return "EventSeeder"
}

func (es *EventSeeder) Description() string {
	return "Crea eventos de ejemplo: conferencias, workshops, meetups y webinars"
}

func (es *EventSeeder) Priority() int {
	return 3 // Ejecutar después de organizaciones
}

func (es *EventSeeder) CanRun(db *gorm.DB) bool {
	var count int64
	db.Model(&models.Event{}).Count(&count)
	return count == 0
}

func (es *EventSeeder) Seed(db *gorm.DB) error {
	logger.Info("📅 Creando eventos de ejemplo...")

	// 1. Obtener organizaciones activas para asignar eventos
	var organizations []models.Organization
	if err := db.Where("status = ? AND is_verified = ?", models.OrgStatusActive, true).Find(&organizations).Error; err != nil {
		return err
	}

	if len(organizations) == 0 {
		logger.Warn("No hay organizaciones verificadas disponibles para crear eventos")
		return nil
	}

	// 2. Crear eventos destacados principales
	if err := es.createFeaturedEvents(db, organizations); err != nil {
		return err
	}

	// 3. Crear eventos próximos
	if err := es.createUpcomingEvents(db, organizations); err != nil {
		return err
	}

	// 4. Crear eventos pasados
	if err := es.createPastEvents(db, organizations); err != nil {
		return err
	}

	// 5. Crear eventos en diferentes estados
	if err := es.createEventsInDifferentStates(db, organizations); err != nil {
		return err
	}

	// 6. Generar eventos adicionales aleatorios
	if err := es.generateRandomEvents(db, organizations, 25); err != nil {
		return err
	}

	logger.Info("✅ Eventos creados exitosamente")
	return nil
}

// createFeaturedEvents crea eventos destacados principales
func (es *EventSeeder) createFeaturedEvents(db *gorm.DB, organizations []models.Organization) error {
	// Buscar la organización "CyberSecurity Spain" para eventos principales
	var mainOrg models.Organization
	for _, org := range organizations {
		if strings.Contains(org.Name, "CyberSecurity Spain") {
			mainOrg = org
			break
		}
	}

	if mainOrg.ID.String() == "" {
		mainOrg = organizations[0] // Usar la primera si no se encuentra
	}

	featuredEvents := []*models.Event{
		{
			Title:           "CyberSec Summit Madrid 2024",
			Description:     "La mayor conferencia de ciberseguridad de España. Dos días llenos de charlas magistrales, talleres prácticos y networking con los mejores profesionales del sector. Speakers internacionales, casos de estudio reales y las últimas tendencias en ciberseguridad.",
			ShortDesc:       "La mayor conferencia de ciberseguridad de España con speakers internacionales",
			Type:            models.EventTypeConference,
			Category:        "Security Conference",
			Level:           "intermediate",
			OrganizationID:  mainOrg.ID.String(),
			StartDate:       time.Now().AddDate(0, 2, 15), // En 2 meses y 15 días
			EndDate:         time.Now().AddDate(0, 2, 16), // 2 días de duración
			Timezone:        "Europe/Madrid",
			IsOnline:        false,
			VenueAddress:    "Palacio de Congresos de Madrid, Paseo de la Castellana, 99",
			VenueName:       "Palacio de Congresos",
			VenueCity:       "Madrid",
			VenueCountry:    "Spain",
			MaxAttendees:    func() *int { i := 500; return &i }(),
			IsFree:          false,
			Price:           func() *int { i := 25000; return &i }(), // 250€
			Currency:        "EUR",
			RegistrationURL: "https://cybersec-summit.es/registro",
			ImageURL:        "https://images.example.com/cybersec-summit-2024.jpg",
			Status:          models.EventStatusPublished,
			IsPublic:        true,
			IsFeatured:      true,
			ContactEmail:    "summit@cybersecurityspain.org",
			ContactPhone:    "+34 91 123 45 67",
			Requirements:    "Conocimientos básicos en ciberseguridad. Se recomienda traer portátil para talleres prácticos.",
			Agenda: `09:00 - Registro y bienvenida
10:00 - Keynote: El futuro de la ciberseguridad
11:00 - Track 1: Threat Hunting / Track 2: Zero Trust Architecture
12:30 - Networking break
13:00 - Panel: Incidentes más críticos de 2024
14:00 - Almuerzo
15:30 - Talleres prácticos (4 tracks paralelos)
17:00 - Clausura y sorteos`,
			MetaTitle:       "CyberSec Summit Madrid 2024 - Conferencia de Ciberseguridad",
			MetaDescription: "Únete a la mayor conferencia de ciberseguridad de España. 2 días, +30 speakers, talleres prácticos y networking.",
		},
		{
			Title:          "Ethical Hacking Bootcamp Intensivo",
			Description:    "Bootcamp intensivo de 3 días para aprender ethical hacking desde cero hasta nivel avanzado. Incluye laboratorios prácticos, certificación y acceso a plataforma de práctica durante 6 meses.",
			ShortDesc:      "Bootcamp intensivo de ethical hacking con laboratorios prácticos",
			Type:           models.EventTypeTraining,
			Category:       "Ethical Hacking",
			Level:          "beginner",
			OrganizationID: mainOrg.ID.String(),
			StartDate:      time.Now().AddDate(0, 1, 20), // En 1 mes y 20 días
			EndDate:        time.Now().AddDate(0, 1, 22), // 3 días
			Timezone:       "Europe/Madrid",
			IsOnline:       true,
			OnlineURL:      "https://training.cybersecurityspain.org/bootcamp",
			MaxAttendees:   func() *int { i := 30; return &i }(),
			IsFree:         false,
			Price:          func() *int { i := 49900; return &i }(), // 499€
			Currency:       "EUR",
			Status:         models.EventStatusPublished,
			IsPublic:       true,
			IsFeatured:     true,
			ContactEmail:   "bootcamp@cybersecurityspain.org",
			Requirements:   "Conocimientos básicos de redes y sistemas. Conexión a internet estable para laboratorios remotos.",
			Agenda: `Día 1: Fundamentos y Reconocimiento
- Metodologías de pentesting
- Information gathering
- Network scanning
- Vulnerability assessment

Día 2: Explotación y Post-explotación
- Web application hacking
- System exploitation
- Privilege escalation
- Persistence techniques

Día 3: Reporting y Certificación
- Report writing
- Remediation advice
- Examen de certificación
- Q&A y recursos adicionales`,
		},
	}

	// Configurar fechas de registro y tags para eventos destacados
	for i, event := range featuredEvents {
		// Fechas de registro
		regStart := time.Now().AddDate(0, 0, -30)   // Comenzó hace 30 días
		regEnd := event.StartDate.AddDate(0, 0, -7) // Termina 7 días antes del evento
		event.RegistrationStartDate = &regStart
		event.RegistrationEndDate = &regEnd

		// Publicado hace un mes
		published := time.Now().AddDate(0, 0, -30)
		event.PublishedAt = &published

		// Configurar ubicación
		if !event.IsOnline {
			event.SetLocation(40.4168, -3.7038) // Madrid
		}

		// Tags específicos
		if i == 0 {
			if err := event.SetTags([]string{"conference", "networking", "speakers", "madrid", "cybersec", "summit"}); err != nil {
				logger.Warnf("Error asignando tags al evento %s: %v", event.Title, err)
			}
		} else {
			if err := event.SetTags([]string{"bootcamp", "training", "ethical-hacking", "certification", "hands-on"}); err != nil {
				logger.Warnf("Error asignando tags al evento %s: %v", event.Title, err)
			}
		}

		if err := db.Create(event).Error; err != nil {
			return err
		}
	}

	return nil
}

// createUpcomingEvents crea eventos próximos
func (es *EventSeeder) createUpcomingEvents(db *gorm.DB, organizations []models.Organization) error {
	upcomingEvents := []*models.Event{
		{
			Title:          "Introducción a la Ciberseguridad para Principiantes",
			Description:    "Workshop gratuito perfecto para quienes quieren comenzar en el mundo de la ciberseguridad. Cubriremos conceptos básicos, herramientas esenciales y rutas de aprendizaje.",
			ShortDesc:      "Workshop gratuito para principiantes en ciberseguridad",
			Type:           models.EventTypeWorkshop,
			Category:       "Security Basics",
			Level:          "beginner",
			OrganizationID: organizations[0].ID.String(),
			StartDate:      time.Now().AddDate(0, 0, 7), // En una semana
			EndDate:        time.Now().AddDate(0, 0, 7).Add(3 * time.Hour),
			Timezone:       "Europe/Madrid",
			IsOnline:       true,
			OnlineURL:      "https://meet.google.com/abc-defg-hij",
			MaxAttendees:   func() *int { i := 100; return &i }(),
			IsFree:         true,
			Status:         models.EventStatusPublished,
			IsPublic:       true,
			ContactEmail:   "eventos@cybersecurityspain.org",
		},
		{
			Title:          "Red Team vs Blue Team: Simulacro en Vivo",
			Description:    "Evento presencial donde dos equipos competirán en tiempo real: Red Team intentando comprometer la infraestructura mientras Blue Team la defiende. Público puede seguir la acción en pantallas gigantes.",
			ShortDesc:      "Competición Red Team vs Blue Team en tiempo real",
			Type:           models.EventTypeCompetition,
			Category:       "Red Team",
			Level:          "advanced",
			OrganizationID: organizations[1].ID.String(),
			StartDate:      time.Now().AddDate(0, 0, 14), // En 2 semanas
			EndDate:        time.Now().AddDate(0, 0, 14).Add(6 * time.Hour),
			IsOnline:       false,
			VenueAddress:   "Campus Universitario, Aula Magna",
			VenueName:      "Universidad Politécnica",
			VenueCity:      "Barcelona",
			VenueCountry:   "Spain",
			MaxAttendees:   func() *int { i := 200; return &i }(),
			IsFree:         true,
			Status:         models.EventStatusPublished,
			IsPublic:       true,
		},
		{
			Title:          "Webinar: Nuevas Amenazas en Cloud Security",
			Description:    "Análisis de las últimas amenazas de seguridad en entornos cloud y las mejores prácticas para proteger infraestructuras AWS, Azure y GCP.",
			ShortDesc:      "Webinar sobre amenazas y protección en cloud",
			Type:           models.EventTypeWebinar,
			Category:       "Cloud Security",
			Level:          "intermediate",
			OrganizationID: organizations[0].ID.String(),
			StartDate:      time.Now().AddDate(0, 0, 21), // En 3 semanas
			EndDate:        time.Now().AddDate(0, 0, 21).Add(90 * time.Minute),
			IsOnline:       true,
			OnlineURL:      "https://zoom.us/webinar/123",
			MaxAttendees:   func() *int { i := 300; return &i }(),
			IsFree:         true,
			Status:         models.EventStatusPublished,
			IsPublic:       true,
		},
	}

	// Crear eventos próximos
	for i, event := range upcomingEvents {
		// Configurar fechas de registro
		regStart := time.Now().AddDate(0, 0, -15)
		regEnd := event.StartDate.AddDate(0, 0, -1)
		event.RegistrationStartDate = &regStart
		event.RegistrationEndDate = &regEnd

		// Publicado recientemente
		published := time.Now().AddDate(0, 0, -7)
		event.PublishedAt = &published

		// Configurar ubicación para eventos presenciales
		if !event.IsOnline {
			if i == 1 { // Barcelona
				event.SetLocation(41.3851, 2.1734)
			}
		}

		// Tags
		switch i {
		case 0:
			if err := event.SetTags([]string{"workshop", "beginners", "security", "basics", "free"}); err != nil {
				logger.Warnf("Error asignando tags al evento %s: %v", event.Title, err)
			}
		case 1:
			if err := event.SetTags([]string{"competition", "red-team", "blue-team", "live", "advanced"}); err != nil {
				logger.Warnf("Error asignando tags al evento %s: %v", event.Title, err)
			}
		case 2:
			if err := event.SetTags([]string{"webinar", "cloud", "threats", "aws", "azure", "gcp"}); err != nil {
				logger.Warnf("Error asignando tags al evento %s: %v", event.Title, err)
			}
		}

		if err := db.Create(event).Error; err != nil {
			return err
		}
	}

	return nil
}

// createPastEvents crea eventos que ya ocurrieron
// createPastEvents crea eventos que ya ocurrieron
func (es *EventSeeder) createPastEvents(db *gorm.DB, organizations []models.Organization) error {
	pastEvents := []*models.Event{
		{
			Title:            "DevSecOps: Integrando Seguridad en CI/CD",
			Description:      "Taller práctico sobre cómo integrar herramientas y prácticas de seguridad en pipelines de DevOps.",
			ShortDesc:        "Taller práctico de DevSecOps e integración de seguridad",
			Type:             models.EventTypeWorkshop,
			Category:         "DevSecOps",
			Level:            "intermediate",
			OrganizationID:   organizations[0].ID.String(),
			StartDate:        time.Now().AddDate(0, 0, -30), // Hace 30 días
			EndDate:          time.Now().AddDate(0, 0, -30).Add(4 * time.Hour),
			Status:           models.EventStatusCompleted,
			IsOnline:         false,
			VenueAddress:     "Calle de Alcalá, 123, Madrid",   // ADD THIS LINE
			VenueName:        "Centro de Formación TechMadrid", // ADD THIS LINE
			VenueCity:        "Madrid",
			VenueCountry:     "Spain",
			MaxAttendees:     func() *int { i := 50; return &i }(),
			CurrentAttendees: 45, // 45 personas asistieron
			IsFree:           false,
			Price:            func() *int { i := 5000; return &i }(), // 50€
			IsPublic:         true,
			ViewsCount:       234, // Muchas visualizaciones
		},
		{
			Title:            "Análisis Forense Digital: Casos Reales",
			Description:      "Conferencia magistral analizando casos reales de análisis forense digital con herramientas profesionales.",
			ShortDesc:        "Conferencia de análisis forense con casos reales",
			Type:             models.EventTypeConference,
			Category:         "Digital Forensics",
			Level:            "advanced",
			OrganizationID:   organizations[1].ID.String(),
			StartDate:        time.Now().AddDate(0, 0, -45), // Hace 45 días
			EndDate:          time.Now().AddDate(0, 0, -45).Add(2 * time.Hour),
			Status:           models.EventStatusCompleted,
			IsOnline:         true,
			OnlineURL:        "https://zoom.us/j/completed-session", // ADD THIS LINE FOR ONLINE EVENT
			MaxAttendees:     func() *int { i := 150; return &i }(),
			CurrentAttendees: 132,
			IsFree:           true,
			IsPublic:         true,
			ViewsCount:       456,
		},
	}

	// Configurar eventos pasados
	for i, event := range pastEvents {
		// Fechas de registro en el pasado
		regStart := event.StartDate.AddDate(0, 0, -21)
		regEnd := event.StartDate.AddDate(0, 0, -1)
		event.RegistrationStartDate = &regStart
		event.RegistrationEndDate = &regEnd

		// Publicado antes del evento
		published := event.StartDate.AddDate(0, 0, -14)
		event.PublishedAt = &published

		// Completado después del evento
		completed := event.EndDate.Add(1 * time.Hour)
		event.CompletedAt = &completed

		// Configurar ubicación para eventos presenciales
		if !event.IsOnline {
			if i == 0 { // Madrid
				event.SetLocation(40.4168, -3.7038)
			}
		}

		// Tags
		if i == 0 {
			if err := event.SetTags([]string{"devsecops", "workshop", "cicd", "automation", "security"}); err != nil {
				logger.Warnf("Error asignando tags al evento %s: %v", event.Title, err)
			}
		} else {
			if err := event.SetTags([]string{"forensics", "conference", "analysis", "tools", "cases"}); err != nil {
				logger.Warnf("Error asignando tags al evento %s: %v", event.Title, err)
			}
		}

		if err := db.Create(event).Error; err != nil {
			return err
		}
	}

	return nil
}

// createEventsInDifferentStates crea eventos en diferentes estados para testing
// createEventsInDifferentStates crea eventos en diferentes estados para testing
func (es *EventSeeder) createEventsInDifferentStates(db *gorm.DB, organizations []models.Organization) error {
	testEvents := []*models.Event{
		{
			Title:          "Evento en Borrador - Test",
			Description:    "Este evento está en estado borrador para testing de la aplicación.",
			Type:           models.EventTypeWorkshop,
			Category:       "Testing",
			Level:          "beginner",
			OrganizationID: organizations[0].ID.String(),
			StartDate:      time.Now().AddDate(0, 1, 0), // En un mes
			EndDate:        time.Now().AddDate(0, 1, 0).Add(2 * time.Hour),
			Status:         models.EventStatusDraft, // Estado borrador
			IsOnline:       true,
			OnlineURL:      "https://meet.google.com/draft-event-test", // ADD THIS LINE
			IsFree:         true,
			IsPublic:       false, // No público hasta publicar
		},
		{
			Title:          "Evento Cancelado - Test",
			Description:    "Este evento fue cancelado para testing de estados.",
			Type:           models.EventTypeMeetup,
			Category:       "Testing",
			Level:          "intermediate",
			OrganizationID: organizations[1].ID.String(),
			StartDate:      time.Now().AddDate(0, 0, -10), // Hace 10 días
			EndDate:        time.Now().AddDate(0, 0, -10).Add(2 * time.Hour),
			Status:         models.EventStatusCanceled, // Cancelado
			IsOnline:       false,
			VenueAddress:   "Plaza Universidad, 1, Barcelona", // ADD THIS LINE
			VenueName:      "Aula Magna Universidad",          // ADD THIS LINE
			VenueCity:      "Barcelona",
			VenueCountry:   "Spain",
			IsFree:         true,
			IsPublic:       true,
		},
	}

	// Crear eventos de testing
	for i, event := range testEvents {
		if i == 1 { // Evento cancelado
			canceled := event.StartDate.AddDate(0, 0, -5) // Cancelado 5 días antes
			event.CanceledAt = &canceled
			// Set location for the cancelled event
			event.SetLocation(41.3851, 2.1734) // Barcelona coordinates
		}

		if err := db.Create(event).Error; err != nil {
			return err
		}
	}

	return nil
}

// generateRandomEvents crea eventos adicionales aleatorios
// generateRandomEvents crea eventos adicionales aleatorios
func (es *EventSeeder) generateRandomEvents(db *gorm.DB, organizations []models.Organization, count int) error {
	eventTitles := [][]string{
		// Workshops
		{"Workshop", "Introducción a", "Taller de", "Curso práctico de", "Masterclass de"},
		// Conferences
		{"Conferencia", "Summit", "Congreso", "Jornadas de", "Simposio de"},
		// Meetups
		{"Meetup", "Encuentro", "Reunión", "Networking", "Charla informal sobre"},
		// Webinars
		{"Webinar", "Charla online", "Sesión virtual", "Presentación", "Demo de"},
		// Training
		{"Bootcamp", "Curso", "Formación", "Entrenamiento", "Certificación en"},
		// Competition
		{"Competición", "CTF", "Hackathon", "Desafío", "Torneo de"},
	}

	topics := []string{
		"Pentesting Avanzado", "Análisis de Malware", "Incident Response", "OSINT",
		"Vulnerability Assessment", "Network Security", "Web Application Security",
		"Mobile Security", "IoT Security", "Cloud Security", "Blockchain Security",
		"AI Security", "Social Engineering", "Cryptography", "Digital Forensics",
		"Red Team Operations", "Blue Team Defense", "Threat Hunting", "SIEM",
		"Zero Trust Architecture", "DevSecOps", "Compliance", "Risk Assessment",
		"Secure Coding", "API Security", "Container Security", "Kubernetes Security",
	}

	eventTypes := []models.EventType{
		models.EventTypeWorkshop,
		models.EventTypeConference,
		models.EventTypeMeetup,
		models.EventTypeWebinar,
		models.EventTypeTraining,
		models.EventTypeCompetition,
	}

	levels := []string{"beginner", "intermediate", "advanced"}

	cities := []struct {
		name     string
		lat, lng float64
	}{
		{"Madrid", 40.4168, -3.7038},
		{"Barcelona", 41.3851, 2.1734},
		{"Valencia", 39.4699, -0.3763},
		{"Sevilla", 37.3886, -5.9823},
		{"Bilbao", 43.2627, -2.9253},
		{"Málaga", 36.7213, -4.4214},
	}

	// Sufijos adicionales para hacer títulos únicos
	timeSuffixes := []string{
		"Primavera", "Verano", "Otoño", "Invierno",
		"Q1", "Q2", "Q3", "Q4",
		"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
	}

	editionSuffixes := []string{
		"Edición Especial", "Nueva Generación", "Versión Pro", "Masterclass",
		"Intensivo", "Express", "Premium", "Avanzado", "Básico", "Completo",
	}

	for i := 0; i < count; i++ {
		eventType := eventTypes[utils.SecureRandInt(len(eventTypes))]

		// Buscar el índice de eventType en eventTypes
		typeIndex := 0
		for idx, et := range eventTypes {
			if et == eventType {
				typeIndex = idx
				break
			}
		}

		titlePrefix := eventTitles[typeIndex][utils.SecureRandInt(len(eventTitles[typeIndex]))]
		topic := topics[utils.SecureRandInt(len(topics))]

		// Crear título único con variaciones
		var title string
		variation := utils.SecureRandInt(5)
		switch variation {
		case 0:
			// Título básico con contador
			title = fmt.Sprintf("%s %s #%d", titlePrefix, topic, i+1)
		case 1:
			// Título con ciudad
			city := cities[utils.SecureRandInt(len(cities))]
			title = fmt.Sprintf("%s %s - %s", titlePrefix, topic, city.name)
		case 2:
			// Título con año
			year := time.Now().Year()
			title = fmt.Sprintf("%s %s %d", titlePrefix, topic, year)
		case 3:
			// Título con sufijo temporal
			timeSuffix := timeSuffixes[utils.SecureRandInt(len(timeSuffixes))]
			title = fmt.Sprintf("%s %s (%s)", titlePrefix, topic, timeSuffix)
		case 4:
			// Título con edición especial
			editionSuffix := editionSuffixes[utils.SecureRandInt(len(editionSuffixes))]
			title = fmt.Sprintf("%s %s - %s", titlePrefix, topic, editionSuffix)
		}

		org := organizations[utils.SecureRandInt(len(organizations))]
		isOnline := utils.SecureRandFloat32() > 0.6 // 40% online
		level := levels[utils.SecureRandInt(len(levels))]

		// Generar fecha aleatoria (pasado, presente, futuro)
		daysOffset := utils.SecureRandInt(120) - 60 // Entre -60 y +60 días
		startDate := time.Now().AddDate(0, 0, daysOffset)

		// Duración del evento
		duration := time.Duration(utils.SecureRandInt(6)+1) * time.Hour
		if eventType == models.EventTypeConference {
			duration = time.Duration(utils.SecureRandInt(3)+1) * 24 * time.Hour // 1-3 días
		}
		endDate := startDate.Add(duration)

		// Determinar estado basado en la fecha
		var status models.EventStatus
		if startDate.Before(time.Now()) {
			if endDate.Before(time.Now()) {
				status = models.EventStatusCompleted
			} else {
				status = models.EventStatusPublished
			}
		} else {
			if utils.SecureRandFloat32() > 0.8 {
				status = models.EventStatusDraft // 20% en borrador
			} else {
				status = models.EventStatusPublished
			}
		}

		isFree := func() bool {
			n, err := rand.Int(rand.Reader, big.NewInt(100))
			if err != nil {
				// Si hay error, por defecto no es gratis
				return false
			}
			return n.Int64() > 40 // 60% gratuitos
		}()

		event := &models.Event{
			Title:          title,
			Description:    fmt.Sprintf("Evento sobre %s organizado por %s. Dirigido a profesionales de nivel %s.", topic, org.Name, level),
			ShortDesc:      fmt.Sprintf("%s sobre %s", titlePrefix, topic),
			Type:           eventType,
			Category:       topic,
			Level:          level,
			OrganizationID: org.ID.String(),
			StartDate:      startDate,
			EndDate:        endDate,
			Timezone:       "Europe/Madrid",
			IsOnline:       isOnline,
			Status:         status,
			IsPublic:       status == models.EventStatusPublished,
			IsFree:         isFree,
		}

		// Configurar ubicación
		if !isOnline {
			city := cities[utils.SecureRandInt(len(cities))]
			event.VenueCity = city.name
			event.VenueCountry = "Spain"
			event.VenueName = fmt.Sprintf("Centro de Eventos %s", city.name)
			event.VenueAddress = fmt.Sprintf("Avenida Principal, %d, %s", utils.SecureRandInt(999)+1, city.name)
			event.SetLocation(city.lat, city.lng)
		} else {
			event.OnlineURL = "https://meet.google.com/generated-link"
		}

		// Ensure venue address is not empty for non-online events
		if !event.IsOnline && strings.TrimSpace(event.VenueAddress) == "" {
			event.VenueAddress = "Dirección por confirmar"
		}

		// Ensure online URL is not empty for online events
		if event.IsOnline && strings.TrimSpace(event.OnlineURL) == "" {
			event.OnlineURL = "https://meet.google.com/generated-link"
		}

		// Configurar precio y capacidad
		if !event.IsFree {
			price := (utils.SecureRandInt(20) + 1) * 500 // Entre 5€ y 100€
			event.Price = &price
		}

		maxAttendees := (utils.SecureRandInt(20) + 1) * 10 // Entre 10 y 200
		event.MaxAttendees = &maxAttendees

		// Para eventos pasados, configurar asistentes
		if status == models.EventStatusCompleted {
			event.CurrentAttendees = utils.SecureRandInt(maxAttendees)
			event.ViewsCount = utils.SecureRandInt(500) + 50 // Entre 50 y 550 vistas
		}

		// Fechas de registro
		regStart := startDate.AddDate(0, 0, -utils.SecureRandInt(30)-7) // 7-37 días antes
		regEnd := startDate.AddDate(0, 0, -utils.SecureRandInt(7)-1)    // 1-7 días antes
		event.RegistrationStartDate = &regStart
		event.RegistrationEndDate = &regEnd

		// Si está publicado, configurar fecha de publicación
		if status == models.EventStatusPublished || status == models.EventStatusCompleted {
			published := regStart.AddDate(0, 0, utils.SecureRandInt(7)+1) // Publicado después de abrir registro
			event.PublishedAt = &published
		}

		// Si está completado, configurar fecha de completado
		if status == models.EventStatusCompleted {
			completed := endDate.Add(time.Duration(utils.SecureRandInt(120)) * time.Minute)
			event.CompletedAt = &completed
		}

		// Configurar tags aleatorios
		tags := []string{strings.ToLower(topic), level}
		if event.IsOnline {
			tags = append(tags, "online")
		} else {
			tags = append(tags, "presencial")
		}
		if event.IsFree {
			tags = append(tags, "free")
		}
		if err := event.SetTags(tags); err != nil {
			logger.Warnf("Error asignando tags al evento %s: %v", event.Title, err)
		}

		// Email de contacto
		event.ContactEmail = fmt.Sprintf("eventos@%s.com", strings.ToLower(strings.ReplaceAll(org.Name, " ", "")))

		if err := db.Create(event).Error; err != nil {
			return fmt.Errorf("error creando evento '%s': %w", event.Title, err)
		}
	}

	return nil
}
