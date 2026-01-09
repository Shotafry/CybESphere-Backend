package models

import (
	"cybesphere-backend/pkg/utils"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// EventStatus define los estados de un evento (máquina de estados)
type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"     // Borrador
	EventStatusPublished EventStatus = "published" // Publicado y visible
	EventStatusCanceled  EventStatus = "canceled"  // Cancelado
	EventStatusCompleted EventStatus = "completed" // Finalizado
)

// EventType define los tipos de evento
type EventType string

const (
	EventTypeConference  EventType = "conference"  // Conferencia
	EventTypeWorkshop    EventType = "workshop"    // Taller
	EventTypeMeetup      EventType = "meetup"      // Reunión/Meetup
	EventTypeWebinar     EventType = "webinar"     // Webinar online
	EventTypeTraining    EventType = "training"    // Formación
	EventTypeCompetition EventType = "competition" // Competición/CTF
	EventTypeOther       EventType = "other"       // Otro tipo
)

// Event modelo para eventos de ciberseguridad
type Event struct {
	BaseModel

	// Información básica
	Title       string `json:"title" gorm:"not null;size:300;index"`
	Slug        string `json:"slug" gorm:"uniqueIndex;not null;size:150"`
	Description string `json:"description" gorm:"type:text"`
	ShortDesc   string `json:"short_desc" gorm:"size:500"` // Descripción corta para listados

	// Clasificación
	Type     EventType `json:"type" gorm:"not null;size:20;index"`
	Category string    `json:"category" gorm:"size:100;index"` // Categoría específica (ej: "Red Team", "Forensics")
	Level    string    `json:"level" gorm:"size:20;index"`     // beginner, intermediate, advanced

	// Fechas y horarios
	StartDate time.Time `json:"start_date" gorm:"not null;index"`
	EndDate   time.Time `json:"end_date" gorm:"not null;index"`
	Timezone  string    `json:"timezone" gorm:"size:50;default:'Europe/Madrid'"`
	Duration  int       `json:"duration"` // Duración en minutos

	// Ubicación
	IsOnline     bool     `json:"is_online" gorm:"not null;default:false;index"`
	VenueAddress string   `json:"venue_address" gorm:"size:500"`
	VenueName    string   `json:"venue_name" gorm:"size:200"`
	VenueCity    string   `json:"venue_city" gorm:"size:100;index"`
	VenueCountry string   `json:"venue_country" gorm:"size:100;index"`
	Latitude     *float64 `json:"latitude" gorm:"index"`
	Longitude    *float64 `json:"longitude" gorm:"index"`
	OnlineURL    string   `json:"online_url" gorm:"size:500"`    // URL para eventos online
	StreamingURL string   `json:"streaming_url" gorm:"size:500"` // URL de streaming

	// Capacidad y registro
	MaxAttendees     *int   `json:"max_attendees"` // null = ilimitado
	CurrentAttendees int    `json:"current_attendees" gorm:"default:0"`
	IsFree           bool   `json:"is_free" gorm:"not null;default:true;index"`
	Price            *int   `json:"price" gorm:"default:0"` // Precio en céntimos
	Currency         string `json:"currency" gorm:"size:3;default:'EUR'"`
	RegistrationURL  string `json:"registration_url" gorm:"size:500"`

	// Estado y visibilidad
	Status     EventStatus `json:"status" gorm:"not null;default:'draft';size:20;index"`
	IsPublic   bool        `json:"is_public" gorm:"not null;default:true;index"`
	IsFeatured bool        `json:"is_featured" gorm:"not null;default:false;index"`
	ViewsCount int         `json:"views_count" gorm:"default:0"`

	// Contenido y recursos
	ImageURL     string         `json:"image_url" gorm:"size:500"`
	BannerURL    string         `json:"banner_url" gorm:"size:500"`
	Tags         datatypes.JSON `json:"tags" gorm:"type:jsonb"`        // Tags para categorización flexible
	Requirements string         `json:"requirements" gorm:"type:text"` // Requisitos técnicos/conocimientos
	Agenda       string         `json:"agenda" gorm:"type:text"`       // Agenda detallada

	// Fechas importantes
	RegistrationStartDate *time.Time `json:"registration_start_date"`
	RegistrationEndDate   *time.Time `json:"registration_end_date"`
	PublishedAt           *time.Time `json:"published_at"`
	CanceledAt            *time.Time `json:"canceled_at"`
	CompletedAt           *time.Time `json:"completed_at"`

	// Información de contacto
	ContactEmail string `json:"contact_email" gorm:"size:255"`
	ContactPhone string `json:"contact_phone" gorm:"size:20"`

	// Metadatos SEO
	MetaTitle       string `json:"meta_title" gorm:"size:200"`
	MetaDescription string `json:"meta_description" gorm:"size:500"`

	// Relaciones
	OrganizationID string        `json:"organization_id" gorm:"not null;size:36;index"`
	Organization   *Organization `json:"organization" gorm:"foreignKey:OrganizationID;references:ID"`
	FavoritedBy    []User        `json:"favorited_by,omitempty" gorm:"many2many:user_favorite_events;"`
}

// TableName especifica el nombre de tabla
func (Event) TableName() string {
	return "events"
}

// BeforeCreate hook para validación y configuración automática
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if err := e.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Validar datos requeridos
	if err := e.ValidateEvent(); err != nil {
		return err
	}

	// Generar slug si no existe
	if e.Slug == "" {
		e.Slug = e.GenerateSlug()
	}

	// Configurar campos automáticos
	e.normalizeFields()
	e.calculateDuration()

	return nil
}

// BeforeUpdate hook para validaciones
func (e *Event) BeforeUpdate(tx *gorm.DB) error {
	if err := e.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// Validar datos
	if err := e.ValidateEvent(); err != nil {
		return err
	}

	// Normalizar campos
	e.normalizeFields()
	e.calculateDuration()

	return nil
}

// ValidateEvent valida los datos del evento
func (e *Event) ValidateEvent() error {
	if strings.TrimSpace(e.Title) == "" {
		return errors.New("event title is required")
	}

	if len(e.Title) < 5 {
		return errors.New("event title must be at least 5 characters")
	}

	if strings.TrimSpace(e.OrganizationID) == "" {
		return errors.New("organization ID is required")
	}

	// Validar fechas
	if e.EndDate.Before(e.StartDate) {
		return errors.New("end date must be after start date")
	}

	// Validar tipo de evento
	if !e.IsValidType() {
		return errors.New("invalid event type")
	}

	// Validar estado
	if !e.IsValidStatus() {
		return errors.New("invalid event status")
	}

	// Validar ubicación
	if !e.IsOnline && e.VenueAddress == "" {
		return errors.New("venue address is required for non-online events")
	}

	if e.IsOnline && e.OnlineURL == "" {
		return errors.New("online URL is required for online events")
	}

	// Validar capacidad
	if e.MaxAttendees != nil && *e.MaxAttendees < 1 {
		return errors.New("max attendees must be at least 1")
	}

	return nil
}

// IsValidType verifica si el tipo de evento es válido
func (e *Event) IsValidType() bool {
	validTypes := []EventType{
		EventTypeConference, EventTypeWorkshop, EventTypeMeetup,
		EventTypeWebinar, EventTypeTraining, EventTypeCompetition, EventTypeOther,
	}

	for _, validType := range validTypes {
		if e.Type == validType {
			return true
		}
	}
	return false
}

// IsValidStatus verifica si el estado es válido
func (e *Event) IsValidStatus() bool {
	return e.Status == EventStatusDraft || e.Status == EventStatusPublished ||
		e.Status == EventStatusCanceled || e.Status == EventStatusCompleted
}

// GenerateSlug genera un slug único basado en el título
func (e *Event) GenerateSlug() string {
	return utils.GenerateSlug(e.Title, 100)
}

// normalizeFields normaliza campos de texto
func (e *Event) normalizeFields() {
	e.Title = strings.TrimSpace(e.Title)
	e.Description = strings.TrimSpace(e.Description)
	e.ShortDesc = strings.TrimSpace(e.ShortDesc)
	e.ContactEmail = strings.ToLower(strings.TrimSpace(e.ContactEmail))
}

// calculateDuration calcula la duración automáticamente
func (e *Event) calculateDuration() {
	if !e.EndDate.IsZero() && !e.StartDate.IsZero() {
		e.Duration = int(e.EndDate.Sub(e.StartDate).Minutes())
	}
}

// Publish publica el evento
func (e *Event) Publish() error {
	if e.Status != EventStatusDraft {
		return errors.New("only draft events can be published")
	}

	now := time.Now()
	e.Status = EventStatusPublished
	e.PublishedAt = &now

	return nil
}

// Cancel cancela el evento
func (e *Event) Cancel() error {
	if e.Status == EventStatusCanceled || e.Status == EventStatusCompleted {
		return errors.New("event is already canceled or completed")
	}

	now := time.Now()
	e.Status = EventStatusCanceled
	e.CanceledAt = &now

	return nil
}

// Complete marca el evento como completado
func (e *Event) Complete() error {
	if e.Status != EventStatusPublished {
		return errors.New("only published events can be completed")
	}

	now := time.Now()
	e.Status = EventStatusCompleted
	e.CompletedAt = &now

	return nil
}

// IsActive verifica si el evento está activo (publicado y no cancelado)
func (e *Event) IsActive() bool {
	return e.Status == EventStatusPublished
}

// IsUpcoming verifica si el evento es futuro
func (e *Event) IsUpcoming() bool {
	return e.StartDate.After(time.Now()) && e.IsActive()
}

// IsPast verifica si el evento ya pasó
func (e *Event) IsPast() bool {
	return e.EndDate.Before(time.Now())
}

// IsRegistrationOpen verifica si el registro está abierto
func (e *Event) IsRegistrationOpen() bool {
	if !e.IsActive() {
		return false
	}

	now := time.Now()

	// Verificar fechas de registro si están configuradas
	if e.RegistrationStartDate != nil && now.Before(*e.RegistrationStartDate) {
		return false
	}

	if e.RegistrationEndDate != nil && now.After(*e.RegistrationEndDate) {
		return false
	}

	// Verificar capacidad
	if e.MaxAttendees != nil && e.CurrentAttendees >= *e.MaxAttendees {
		return false
	}

	return true
}

// HasAvailableSpots verifica si hay cupos disponibles
func (e *Event) HasAvailableSpots() bool {
	if e.MaxAttendees == nil {
		return true // Capacidad ilimitada
	}

	return e.CurrentAttendees < *e.MaxAttendees
}

// GetAvailableSpots retorna el número de cupos disponibles
func (e *Event) GetAvailableSpots() *int {
	if e.MaxAttendees == nil {
		return nil // Ilimitado
	}

	available := *e.MaxAttendees - e.CurrentAttendees
	if available < 0 {
		available = 0
	}

	return &available
}

// SetLocation establece la geolocalización
func (e *Event) SetLocation(latitude, longitude float64) {
	e.Latitude = &latitude
	e.Longitude = &longitude
}

// HasLocation verifica si tiene coordenadas
func (e *Event) HasLocation() bool {
	return e.Latitude != nil && e.Longitude != nil
}

// AddTag agrega un tag al evento
func (e *Event) AddTag(tag string) error {
	tag = strings.TrimSpace(strings.ToLower(tag))
	if tag == "" {
		return nil
	}

	var tags []string
	if e.Tags != nil {
		if err := json.Unmarshal(e.Tags, &tags); err != nil {
			tags = []string{}
		}
	}

	// Verificar si ya existe
	for _, existingTag := range tags {
		if existingTag == tag {
			return nil
		}
	}

	tags = append(tags, tag)
	data, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	e.Tags = datatypes.JSON(data)
	return nil
}

// RemoveTag remueve un tag del evento
func (e *Event) RemoveTag(tag string) error {
	tag = strings.TrimSpace(strings.ToLower(tag))

	var tags []string
	if e.Tags != nil {
		if err := json.Unmarshal(e.Tags, &tags); err != nil {
			return err
		}
	}

	for i, existingTag := range tags {
		if existingTag == tag {
			tags = append(tags[:i], tags[i+1:]...)
			data, err := json.Marshal(tags)
			if err != nil {
				return err
			}
			e.Tags = datatypes.JSON(data)
			return nil
		}
	}

	return nil
}

// SetTags establece todos los tags del evento
func (e *Event) SetTags(tags []string) error {
	var cleanTags []string
	for _, tag := range tags {
		cleanTag := strings.TrimSpace(strings.ToLower(tag))
		if cleanTag != "" {
			cleanTags = append(cleanTags, cleanTag)
		}
	}

	data, err := json.Marshal(cleanTags)
	if err != nil {
		return err
	}
	e.Tags = datatypes.JSON(data)
	return nil
}

// GetTags retorna los tags del evento
func (e *Event) GetTags() []string {
	var tags []string
	if err := json.Unmarshal(e.Tags, &tags); err != nil {
		return []string{}
	}
	return tags
}

// IncrementViews incrementa el contador de visualizaciones
func (e *Event) IncrementViews(tx *gorm.DB) error {
	return tx.Model(e).UpdateColumn("views_count", gorm.Expr("views_count + ?", 1)).Error
}

// GetAuditData implementa AuditableModel
func (e *Event) GetAuditData() map[string]interface{} {
	return map[string]interface{}{
		"id":              e.ID,
		"title":           e.Title,
		"slug":            e.Slug,
		"type":            e.Type,
		"status":          e.Status,
		"organization_id": e.OrganizationID,
		"start_date":      e.StartDate,
		"end_date":        e.EndDate,
	}
}

func (e Event) GetID() string           { return e.ID.String() }
func (e Event) GetCreatedAt() time.Time { return e.CreatedAt }
func (e Event) GetUpdatedAt() time.Time { return e.UpdatedAt }
