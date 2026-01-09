package models

import (
	"cybesphere-backend/pkg/utils"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// OrganizationStatus define los estados de una organización
type OrganizationStatus string

const (
	OrgStatusPending   OrganizationStatus = "pending"   // Pendiente de verificación
	OrgStatusActive    OrganizationStatus = "active"    // Verificada y activa
	OrgStatusSuspended OrganizationStatus = "suspended" // Suspendida temporalmente
	OrgStatusInactive  OrganizationStatus = "inactive"  // Inactiva por decisión propia
)

// Organization modelo para organizaciones que crean eventos
type Organization struct {
	BaseModel

	// Información básica
	Name        string `json:"name" gorm:"not null;size:200;index"`
	Slug        string `json:"slug" gorm:"uniqueIndex;not null;size:100"`
	Description string `json:"description" gorm:"type:text"`
	Website     string `json:"website" gorm:"size:255"`

	// Información de contacto
	Email      string `json:"email" gorm:"not null;size:255;index"`
	Phone      string `json:"phone" gorm:"size:20"`
	Address    string `json:"address" gorm:"size:500"`
	City       string `json:"city" gorm:"size:100"`
	Country    string `json:"country" gorm:"size:100;index"`
	PostalCode string `json:"postal_code" gorm:"size:20"`

	// Geolocalización
	Latitude  *float64 `json:"latitude" gorm:"index"`
	Longitude *float64 `json:"longitude" gorm:"index"`

	// Branding y medios
	LogoURL        string `json:"logo_url" gorm:"size:500"`
	BannerURL      string `json:"banner_url" gorm:"size:500"`
	PrimaryColor   string `json:"primary_color" gorm:"size:7"`   // Hex color
	SecondaryColor string `json:"secondary_color" gorm:"size:7"` // Hex color

	// Redes sociales
	LinkedIn  string `json:"linkedin" gorm:"size:255"`
	Twitter   string `json:"twitter" gorm:"size:255"`
	Facebook  string `json:"facebook" gorm:"size:255"`
	Instagram string `json:"instagram" gorm:"size:255"`
	YouTube   string `json:"youtube" gorm:"size:255"`

	// Estado y verificación
	Status     OrganizationStatus `json:"status" gorm:"not null;default:'pending';size:20;index"`
	IsVerified bool               `json:"is_verified" gorm:"not null;default:false;index"`
	VerifiedAt *string            `json:"verified_at,omitempty"`
	VerifiedBy *string            `json:"verified_by,omitempty" gorm:"size:36"`

	// Metadatos de verificación
	TaxID            string `json:"tax_id,omitempty" gorm:"size:50"` // NIF, CIF, etc.
	LegalName        string `json:"legal_name,omitempty" gorm:"size:300"`
	RegistrationDocs string `json:"registration_docs,omitempty" gorm:"size:500"` // URLs de documentos

	// Configuraciones
	EventsCount     int  `json:"events_count" gorm:"default:0"`
	MaxEvents       *int `json:"max_events,omitempty"` // Límite de eventos (null = ilimitado)
	CanCreateEvents bool `json:"can_create_events" gorm:"not null;default:true"`

	// Relaciones
	Users  []User  `json:"users,omitempty" gorm:"foreignKey:OrganizationID"`
	Events []Event `json:"events,omitempty" gorm:"foreignKey:OrganizationID"`
}

// TableName especifica el nombre de tabla
func (Organization) TableName() string {
	return "organizations"
}

// BeforeCreate hook para validación y slug generation
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	if err := o.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Validar datos requeridos
	if err := o.ValidateOrganization(); err != nil {
		return err
	}

	// Generar slug si no existe
	if o.Slug == "" {
		o.Slug = o.GenerateSlug()
	}

	// Normalizar campos
	o.normalizeFields()

	return nil
}

// BeforeUpdate hook para validaciones
func (o *Organization) BeforeUpdate(tx *gorm.DB) error {
	if err := o.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// Validar datos
	if err := o.ValidateOrganization(); err != nil {
		return err
	}

	// Normalizar campos
	o.normalizeFields()

	return nil
}

// ValidateOrganization valida los datos de la organización
func (o *Organization) ValidateOrganization() error {
	if strings.TrimSpace(o.Name) == "" {
		return errors.New("organization name is required")
	}

	if len(o.Name) < 3 {
		return errors.New("organization name must be at least 3 characters")
	}

	if strings.TrimSpace(o.Email) == "" {
		return errors.New("organization email is required")
	}

	// Validar status
	if !o.IsValidStatus() {
		return errors.New("invalid organization status")
	}

	// Validar colores si están presentes
	if o.PrimaryColor != "" && !isValidHexColor(o.PrimaryColor) {
		return errors.New("invalid primary color format")
	}

	if o.SecondaryColor != "" && !isValidHexColor(o.SecondaryColor) {
		return errors.New("invalid secondary color format")
	}

	return nil
}

// IsValidStatus verifica si el status es válido
func (o *Organization) IsValidStatus() bool {
	return o.Status == OrgStatusPending || o.Status == OrgStatusActive ||
		o.Status == OrgStatusSuspended || o.Status == OrgStatusInactive
}

// GenerateSlug genera un slug único basado en el nombre
func (o *Organization) GenerateSlug() string {
	return utils.GenerateSlug(o.Name, 100)
}

// normalizeFields normaliza campos de texto
func (o *Organization) normalizeFields() {
	o.Name = strings.TrimSpace(o.Name)
	o.Email = strings.ToLower(strings.TrimSpace(o.Email))
	o.Website = strings.TrimSpace(o.Website)
	o.Description = strings.TrimSpace(o.Description)
}

// IsActive verifica si la organización está activa
func (o *Organization) IsActive() bool {
	return o.Status == OrgStatusActive
}

// CanCreateEvent verifica si puede crear eventos
func (o *Organization) CanCreateEvent() bool {
	if !o.IsActive() || !o.CanCreateEvents {
		return false
	}

	// Verificar límite de eventos si está configurado
	if o.MaxEvents != nil && o.EventsCount >= *o.MaxEvents {
		return false
	}

	return true
}

// Verify marca la organización como verificada
func (o *Organization) Verify(verifiedBy string) {
	o.IsVerified = true
	o.VerifiedBy = &verifiedBy
	now := strings.TrimSpace(string(rune(1))) // Placeholder para timestamp
	o.VerifiedAt = &now

	if o.Status == OrgStatusPending {
		o.Status = OrgStatusActive
	}
}

// Suspend suspende la organización
func (o *Organization) Suspend() {
	o.Status = OrgStatusSuspended
	o.CanCreateEvents = false
}

// Activate activa la organización
func (o *Organization) Activate() {
	o.Status = OrgStatusActive
	o.CanCreateEvents = true
}

// Deactivate desactiva la organización
func (o *Organization) Deactivate() {
	o.Status = OrgStatusInactive
	o.CanCreateEvents = false
}

// IncrementEventsCount incrementa el contador de eventos
func (o *Organization) IncrementEventsCount(tx *gorm.DB) error {
	return tx.Model(o).UpdateColumn("events_count", gorm.Expr("events_count + ?", 1)).Error
}

// DecrementEventsCount decrementa el contador de eventos
func (o *Organization) DecrementEventsCount(tx *gorm.DB) error {
	return tx.Model(o).UpdateColumn("events_count", gorm.Expr("GREATEST(events_count - ?, 0)", 1)).Error
}

// SetLocation establece la geolocalización
func (o *Organization) SetLocation(latitude, longitude float64) {
	o.Latitude = &latitude
	o.Longitude = &longitude
}

// HasLocation verifica si tiene coordenadas de geolocalización
func (o *Organization) HasLocation() bool {
	return o.Latitude != nil && o.Longitude != nil
}

// SetBranding establece los colores de branding
func (o *Organization) SetBranding(primaryColor, secondaryColor string) error {
	if primaryColor != "" && !isValidHexColor(primaryColor) {
		return errors.New("invalid primary color format")
	}

	if secondaryColor != "" && !isValidHexColor(secondaryColor) {
		return errors.New("invalid secondary color format")
	}

	o.PrimaryColor = primaryColor
	o.SecondaryColor = secondaryColor
	return nil
}

// GetAuditData implementa AuditableModel
func (o *Organization) GetAuditData() map[string]interface{} {
	return map[string]interface{}{
		"id":          o.ID,
		"name":        o.Name,
		"slug":        o.Slug,
		"email":       o.Email,
		"status":      o.Status,
		"is_verified": o.IsVerified,
	}
}

// isValidHexColor valida formato de color hexadecimal
func isValidHexColor(color string) bool {
	if len(color) != 7 || !strings.HasPrefix(color, "#") {
		return false
	}

	for _, char := range color[1:] {
		if !((char >= '0' && char <= '9') || (char >= 'A' && char <= 'F') || (char >= 'a' && char <= 'f')) {
			return false
		}
	}

	return true
}

// Métodos de base model implementados
func (o Organization) GetID() string           { return o.ID.String() }
func (o Organization) GetCreatedAt() time.Time { return o.CreatedAt }
func (o Organization) GetUpdatedAt() time.Time { return o.UpdatedAt }
