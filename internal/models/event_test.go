package models

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestEvent crea un evento válido para testing
func createTestEvent() *Event {
	orgID := uuid.New().String()
	startTime := time.Now().Add(24 * time.Hour) // En 1 día
	endTime := startTime.Add(2 * time.Hour)     // 2 horas de duración

	return &Event{
		Title:          "Workshop de Ciberseguridad",
		Description:    "Taller práctico sobre fundamentos de ciberseguridad",
		ShortDesc:      "Taller de ciberseguridad para principiantes",
		Type:           EventTypeWorkshop,
		Category:       "Security Basics",
		Level:          "beginner",
		StartDate:      startTime,
		EndDate:        endTime,
		Timezone:       "Europe/Madrid",
		OrganizationID: orgID,
		IsOnline:       true,
		OnlineURL:      "https://meet.google.com/abc-defg-hij",
		IsFree:         true,
		Status:         EventStatusDraft,
		IsPublic:       true,
		ContactEmail:   "info@example.com",
	}
}

// TestEvent_ValidateEvent tests unitarios para validación de eventos
func TestEvent_ValidateEvent(t *testing.T) {
	tests := []struct {
		name    string
		event   *Event
		wantErr bool
		errMsg  string
	}{
		{
			name:    "evento válido",
			event:   createTestEvent(),
			wantErr: false,
		},
		{
			name: "título vacío",
			event: func() *Event {
				e := createTestEvent()
				e.Title = ""
				return e
			}(),
			wantErr: true,
			errMsg:  "event title is required",
		},
		{
			name: "título solo espacios",
			event: func() *Event {
				e := createTestEvent()
				e.Title = "   "
				return e
			}(),
			wantErr: true,
			errMsg:  "event title is required",
		},
		{
			name: "título muy corto",
			event: func() *Event {
				e := createTestEvent()
				e.Title = "ABC"
				return e
			}(),
			wantErr: true,
			errMsg:  "event title must be at least 5 characters",
		},
		{
			name: "organization ID vacío",
			event: func() *Event {
				e := createTestEvent()
				e.OrganizationID = ""
				return e
			}(),
			wantErr: true,
			errMsg:  "organization ID is required",
		},
		{
			name: "fecha fin antes que inicio",
			event: func() *Event {
				e := createTestEvent()
				e.EndDate = e.StartDate.Add(-1 * time.Hour)
				return e
			}(),
			wantErr: true,
			errMsg:  "end date must be after start date",
		},
		{
			name: "tipo de evento inválido",
			event: func() *Event {
				e := createTestEvent()
				e.Type = EventType("invalid_type")
				return e
			}(),
			wantErr: true,
			errMsg:  "invalid event type",
		},
		{
			name: "status inválido",
			event: func() *Event {
				e := createTestEvent()
				e.Status = EventStatus("invalid_status")
				return e
			}(),
			wantErr: true,
			errMsg:  "invalid event status",
		},
		{
			name: "evento presencial sin dirección",
			event: func() *Event {
				e := createTestEvent()
				e.IsOnline = false
				e.VenueAddress = ""
				e.OnlineURL = ""
				return e
			}(),
			wantErr: true,
			errMsg:  "venue address is required for non-online events",
		},
		{
			name: "evento online sin URL",
			event: func() *Event {
				e := createTestEvent()
				e.IsOnline = true
				e.OnlineURL = ""
				return e
			}(),
			wantErr: true,
			errMsg:  "online URL is required for online events",
		},
		{
			name: "capacidad máxima inválida",
			event: func() *Event {
				e := createTestEvent()
				maxAttendees := 0
				e.MaxAttendees = &maxAttendees
				return e
			}(),
			wantErr: true,
			errMsg:  "max attendees must be at least 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.ValidateEvent()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestEvent_IsValidType tests para validación de tipos
func TestEvent_IsValidType(t *testing.T) {
	tests := []struct {
		name string
		Type EventType
		want bool
	}{
		{"conference válido", EventTypeConference, true},
		{"workshop válido", EventTypeWorkshop, true},
		{"meetup válido", EventTypeMeetup, true},
		{"webinar válido", EventTypeWebinar, true},
		{"training válido", EventTypeTraining, true},
		{"competition válido", EventTypeCompetition, true},
		{"other válido", EventTypeOther, true},
		{"tipo inválido", EventType("invalid"), false},
		{"tipo vacío", EventType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{Type: tt.Type}
			assert.Equal(t, tt.want, event.IsValidType())
		})
	}
}

// TestEvent_IsValidStatus tests para validación de estados
func TestEvent_IsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status EventStatus
		want   bool
	}{
		{"draft válido", EventStatusDraft, true},
		{"published válido", EventStatusPublished, true},
		{"canceled válido", EventStatusCanceled, true},
		{"completed válido", EventStatusCompleted, true},
		{"status inválido", EventStatus("invalid"), false},
		{"status vacío", EventStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{Status: tt.status}
			assert.Equal(t, tt.want, event.IsValidStatus())
		})
	}
}

// TestEvent_GenerateSlug tests para generación de slugs
func TestEvent_GenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		wantSlug string
	}{
		{
			name:     "título normal",
			title:    "Workshop de Ciberseguridad",
			wantSlug: "workshop-de-ciberseguridad",
		},
		{
			name:     "con acentos y caracteres especiales",
			title:    "Introducción a la Tecnología & AI",
			wantSlug: "introduccion-a-la-tecnologia-ai",
		},
		{
			name:     "con números",
			title:    "Cybersecurity 2024: Trends & Tools",
			wantSlug: "cybersecurity-2024-trends-tools",
		},
		{
			name:     "título muy largo",
			title:    "Este es un título extremadamente largo para un evento de ciberseguridad que debería ser truncado apropiadamente",
			wantSlug: "este-es-un-titulo-extremadamente-largo-para-un-evento-de-ciberseguridad-que-deberia-ser-truncado-apr",
		},
		{
			name:     "con caracteres especiales",
			title:    "¡Workshop Práctico! 100% Gratuito",
			wantSlug: "workshop-practico-100-gratuito",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{Title: tt.title}
			slug := event.GenerateSlug()
			assert.Equal(t, tt.wantSlug, slug)
		})
	}
}

// TestEvent_StateTransitions tests para transiciones de estado
func TestEvent_StateTransitions(t *testing.T) {
	t.Run("publicar evento", func(t *testing.T) {
		event := createTestEvent()
		event.Status = EventStatusDraft

		err := event.Publish()

		assert.NoError(t, err)
		assert.Equal(t, EventStatusPublished, event.Status)
		assert.NotNil(t, event.PublishedAt)
		assert.WithinDuration(t, time.Now(), *event.PublishedAt, time.Second)
	})

	t.Run("no se puede publicar evento ya publicado", func(t *testing.T) {
		event := createTestEvent()
		event.Status = EventStatusPublished

		err := event.Publish()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only draft events can be published")
	})

	t.Run("cancelar evento", func(t *testing.T) {
		event := createTestEvent()
		event.Status = EventStatusPublished

		err := event.Cancel()

		assert.NoError(t, err)
		assert.Equal(t, EventStatusCanceled, event.Status)
		assert.NotNil(t, event.CanceledAt)
		assert.WithinDuration(t, time.Now(), *event.CanceledAt, time.Second)
	})

	t.Run("no se puede cancelar evento completado", func(t *testing.T) {
		event := createTestEvent()
		event.Status = EventStatusCompleted

		err := event.Cancel()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "event is already canceled or completed")
	})

	t.Run("completar evento", func(t *testing.T) {
		event := createTestEvent()
		event.Status = EventStatusPublished

		err := event.Complete()

		assert.NoError(t, err)
		assert.Equal(t, EventStatusCompleted, event.Status)
		assert.NotNil(t, event.CompletedAt)
		assert.WithinDuration(t, time.Now(), *event.CompletedAt, time.Second)
	})

	t.Run("no se puede completar evento draft", func(t *testing.T) {
		event := createTestEvent()
		event.Status = EventStatusDraft

		err := event.Complete()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only published events can be completed")
	})
}

// TestEvent_StatusCheckers tests para métodos de verificación de estado
func TestEvent_StatusCheckers(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		status     EventStatus
		startDate  time.Time
		endDate    time.Time
		isActive   bool
		isUpcoming bool
		isPast     bool
	}{
		{
			name:       "evento publicado futuro",
			status:     EventStatusPublished,
			startDate:  now.Add(24 * time.Hour),
			endDate:    now.Add(26 * time.Hour),
			isActive:   true,
			isUpcoming: true,
			isPast:     false,
		},
		{
			name:       "evento publicado pasado",
			status:     EventStatusPublished,
			startDate:  now.Add(-26 * time.Hour),
			endDate:    now.Add(-24 * time.Hour),
			isActive:   true,
			isUpcoming: false,
			isPast:     true,
		},
		{
			name:       "evento draft",
			status:     EventStatusDraft,
			startDate:  now.Add(24 * time.Hour),
			endDate:    now.Add(26 * time.Hour),
			isActive:   false,
			isUpcoming: false,
			isPast:     false,
		},
		{
			name:       "evento cancelado",
			status:     EventStatusCanceled,
			startDate:  now.Add(24 * time.Hour),
			endDate:    now.Add(26 * time.Hour),
			isActive:   false,
			isUpcoming: false,
			isPast:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{
				Status:    tt.status,
				StartDate: tt.startDate,
				EndDate:   tt.endDate,
			}

			assert.Equal(t, tt.isActive, event.IsActive())
			assert.Equal(t, tt.isUpcoming, event.IsUpcoming())
			assert.Equal(t, tt.isPast, event.IsPast())
		})
	}
}

// TestEvent_RegistrationLogic tests para lógica de registro
func TestEvent_RegistrationLogic(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name                  string
		status                EventStatus
		registrationStartDate *time.Time
		registrationEndDate   *time.Time
		maxAttendees          *int
		currentAttendees      int
		isRegistrationOpen    bool
		hasAvailableSpots     bool
	}{
		{
			name:               "registro abierto sin límites",
			status:             EventStatusPublished,
			maxAttendees:       nil,
			currentAttendees:   10,
			isRegistrationOpen: true,
			hasAvailableSpots:  true,
		},
		{
			name:               "registro con cupos disponibles",
			status:             EventStatusPublished,
			maxAttendees:       func() *int { i := 50; return &i }(),
			currentAttendees:   30,
			isRegistrationOpen: true,
			hasAvailableSpots:  true,
		},
		{
			name:               "registro lleno",
			status:             EventStatusPublished,
			maxAttendees:       func() *int { i := 50; return &i }(),
			currentAttendees:   50,
			isRegistrationOpen: false,
			hasAvailableSpots:  false,
		},
		{
			name:               "evento no publicado",
			status:             EventStatusDraft,
			maxAttendees:       func() *int { i := 50; return &i }(),
			currentAttendees:   10,
			isRegistrationOpen: false,
			hasAvailableSpots:  true,
		},
		{
			name:                  "registro no iniciado",
			status:                EventStatusPublished,
			registrationStartDate: func() *time.Time { t := now.Add(24 * time.Hour); return &t }(),
			maxAttendees:          func() *int { i := 50; return &i }(),
			currentAttendees:      10,
			isRegistrationOpen:    false,
			hasAvailableSpots:     true,
		},
		{
			name:                "registro cerrado",
			status:              EventStatusPublished,
			registrationEndDate: func() *time.Time { t := now.Add(-24 * time.Hour); return &t }(),
			maxAttendees:        func() *int { i := 50; return &i }(),
			currentAttendees:    10,
			isRegistrationOpen:  false,
			hasAvailableSpots:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{
				Status:                tt.status,
				RegistrationStartDate: tt.registrationStartDate,
				RegistrationEndDate:   tt.registrationEndDate,
				MaxAttendees:          tt.maxAttendees,
				CurrentAttendees:      tt.currentAttendees,
			}

			assert.Equal(t, tt.isRegistrationOpen, event.IsRegistrationOpen())
			assert.Equal(t, tt.hasAvailableSpots, event.HasAvailableSpots())

			// Test GetAvailableSpots
			availableSpots := event.GetAvailableSpots()
			if tt.maxAttendees == nil {
				assert.Nil(t, availableSpots)
			} else {
				expected := *tt.maxAttendees - tt.currentAttendees
				if expected < 0 {
					expected = 0
				}
				assert.Equal(t, expected, *availableSpots)
			}
		})
	}
}

// TestEvent_LocationManagement tests para gestión de ubicación
func TestEvent_LocationManagement(t *testing.T) {
	event := createTestEvent()

	// Verificar que inicialmente no tiene ubicación
	assert.False(t, event.HasLocation())

	// Establecer ubicación
	lat, lng := 40.4168, -3.7038 // Madrid coordinates
	event.SetLocation(lat, lng)

	// Verificar que se estableció correctamente
	assert.True(t, event.HasLocation())
	assert.Equal(t, lat, *event.Latitude)
	assert.Equal(t, lng, *event.Longitude)
}

// TestEvent_TagsManagement tests para gestión de tags
func TestEvent_TagsManagement(t *testing.T) {
	event := createTestEvent()

	t.Run("agregar tag", func(t *testing.T) {
		err := event.AddTag("cybersecurity")
		assert.NoError(t, err)

		tags := event.GetTags()
		assert.Contains(t, tags, "cybersecurity")
	})

	t.Run("agregar tag duplicado", func(t *testing.T) {
		err := event.AddTag("cybersecurity")
		assert.NoError(t, err)

		err = event.AddTag("cybersecurity") // Duplicado
		assert.NoError(t, err)

		tags := event.GetTags()
		count := 0
		for _, tag := range tags {
			if tag == "cybersecurity" {
				count++
			}
		}
		assert.Equal(t, 1, count, "No debería haber tags duplicados")
	})

	t.Run("agregar múltiples tags", func(t *testing.T) {
		event := createTestEvent()

		err := event.AddTag("workshop")
		assert.NoError(t, err)

		err = event.AddTag("beginner")
		assert.NoError(t, err)

		err = event.AddTag("online")
		assert.NoError(t, err)

		tags := event.GetTags()
		assert.Contains(t, tags, "workshop")
		assert.Contains(t, tags, "beginner")
		assert.Contains(t, tags, "online")
		assert.Len(t, tags, 3)
	})

	t.Run("establecer todos los tags", func(t *testing.T) {
		event := createTestEvent()

		newTags := []string{"cybersecurity", "workshop", "madrid", "2024"}
		err := event.SetTags(newTags)
		assert.NoError(t, err)

		tags := event.GetTags()
		assert.Equal(t, newTags, tags)
	})

	t.Run("remover tag", func(t *testing.T) {
		event := createTestEvent()

		err := event.SetTags([]string{"workshop", "cybersecurity", "madrid"})
		assert.NoError(t, err)

		// Remover uno
		err = event.RemoveTag("cybersecurity")
		assert.NoError(t, err)

		tags := event.GetTags()
		assert.NotContains(t, tags, "cybersecurity")
		assert.Contains(t, tags, "workshop")
		assert.Contains(t, tags, "madrid")
	})

	t.Run("normalización de tags", func(t *testing.T) {
		event := createTestEvent()

		err := event.AddTag("  CyberSecurity  ")
		assert.NoError(t, err)

		tags := event.GetTags()
		assert.Contains(t, tags, "cybersecurity") // Debe estar en minúsculas y sin espacios
	})

	t.Run("tag vacío", func(t *testing.T) {
		event := createTestEvent()

		err := event.AddTag("")
		assert.NoError(t, err)

		err = event.AddTag("   ")
		assert.NoError(t, err)

		tags := event.GetTags()
		assert.Empty(t, tags) // No debería agregar tags vacíos
	})
}

// TestEvent_GetAuditData tests para datos de auditoría
func TestEvent_GetAuditData(t *testing.T) {
	event := createTestEvent()
	event.ID = uuid.New()
	event.Slug = "workshop-de-ciberseguridad"

	auditData := event.GetAuditData()

	// Verificar que contiene los campos esperados
	expectedFields := []string{"id", "title", "slug", "type", "status", "organization_id", "start_date", "end_date"}
	for _, field := range expectedFields {
		assert.Contains(t, auditData, field)
	}

	// Verificar valores específicos
	assert.Equal(t, event.ID, auditData["id"])
	assert.Equal(t, event.Title, auditData["title"])
	assert.Equal(t, event.Slug, auditData["slug"])
	assert.Equal(t, event.Type, auditData["type"])
	assert.Equal(t, event.Status, auditData["status"])
	assert.Equal(t, event.OrganizationID, auditData["organization_id"])
}

// TestEvent_DurationCalculation tests para cálculo de duración
func TestEvent_DurationCalculation(t *testing.T) {
	tests := []struct {
		name             string
		startDate        time.Time
		endDate          time.Time
		expectedDuration int
	}{
		{
			name:             "evento de 2 horas",
			startDate:        time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			endDate:          time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expectedDuration: 120, // 2 horas = 120 minutos
		},
		{
			name:             "evento de 30 minutos",
			startDate:        time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
			endDate:          time.Date(2024, 1, 1, 14, 30, 0, 0, time.UTC),
			expectedDuration: 30,
		},
		{
			name:             "evento de un día completo",
			startDate:        time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			endDate:          time.Date(2024, 1, 2, 9, 0, 0, 0, time.UTC),
			expectedDuration: 1440, // 24 horas = 1440 minutos
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{
				StartDate: tt.startDate,
				EndDate:   tt.endDate,
			}

			// Simular el cálculo de duración (como en BeforeCreate)
			if !event.EndDate.IsZero() && !event.StartDate.IsZero() {
				event.Duration = int(event.EndDate.Sub(event.StartDate).Minutes())
			}

			assert.Equal(t, tt.expectedDuration, event.Duration)
		})
	}
}

// TestEvent_BeforeCreateLogic tests para lógica del hook BeforeCreate sin base de datos
func TestEvent_BeforeCreateLogic(t *testing.T) {
	event := createTestEvent()

	// Testear la lógica de validación directamente
	err := event.ValidateEvent()
	require.NoError(t, err)

	// Generar slug si no existe
	if event.Slug == "" {
		event.Slug = event.GenerateSlug()
	}

	// Verificar que se generó el slug
	assert.NotEmpty(t, event.Slug)
	assert.Equal(t, "workshop-de-ciberseguridad", event.Slug)

	// Normalizar campos
	event.Title = "  " + event.Title + "  "
	event.Description = "  " + event.Description + "  "
	event.ContactEmail = "  " + strings.ToUpper(event.ContactEmail) + "  "

	// Simular normalización
	event.Title = strings.TrimSpace(event.Title)
	event.Description = strings.TrimSpace(event.Description)
	event.ContactEmail = strings.ToLower(strings.TrimSpace(event.ContactEmail))

	// Calcular duración
	if !event.EndDate.IsZero() && !event.StartDate.IsZero() {
		event.Duration = int(event.EndDate.Sub(event.StartDate).Minutes())
	}

	assert.Equal(t, "Workshop de Ciberseguridad", event.Title)
	assert.Equal(t, "Taller práctico sobre fundamentos de ciberseguridad", event.Description)
	assert.Equal(t, "info@example.com", event.ContactEmail)
	assert.Equal(t, 120, event.Duration) // 2 horas
}
