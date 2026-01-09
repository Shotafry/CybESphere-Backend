package models

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole define los roles disponibles en el sistema
type UserRole string

const (
	RoleAdmin     UserRole = "admin"     // Control total del sistema
	RoleOrganizer UserRole = "organizer" // Control de su organización y eventos
	RoleUser      UserRole = "user"      // Solo consulta y favoritos
)

// User modelo de usuario con autenticación y geolocalización
type User struct {
	BaseModel

	// Información personal
	Email     string `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Password  string `json:"-" gorm:"not null;size:255"` // Nunca exponer en JSON
	FirstName string `json:"first_name" gorm:"not null;size:100"`
	LastName  string `json:"last_name" gorm:"not null;size:100"`

	// Perfil profesional
	Company  string `json:"company" gorm:"size:200"`
	Position string `json:"position" gorm:"size:200"`
	Bio      string `json:"bio" gorm:"type:text"`
	Website  string `json:"website" gorm:"size:255"`
	LinkedIn string `json:"linkedin" gorm:"size:255"`
	Twitter  string `json:"twitter" gorm:"size:255"`

	// Sistema de roles
	Role UserRole `json:"role" gorm:"not null;default:'user';size:20;index"`

	// Estado de la cuenta
	IsActive    bool       `json:"is_active" gorm:"not null;default:true;index"`
	IsVerified  bool       `json:"is_verified" gorm:"not null;default:false"`
	LastLoginAt *time.Time `json:"last_login_at"`

	// Geolocalización (para búsquedas espaciales)
	Latitude  *float64 `json:"latitude" gorm:"index"`
	Longitude *float64 `json:"longitude" gorm:"index"`
	City      string   `json:"city" gorm:"size:100"`
	Country   string   `json:"country" gorm:"size:100"`

	// Configuraciones de usuario
	Timezone          string `json:"timezone" gorm:"size:50;default:'Europe/Madrid'"`
	Language          string `json:"language" gorm:"size:10;default:'es'"`
	NewsletterEnabled bool   `json:"newsletter_enabled" gorm:"default:true"`

	// Relaciones
	OrganizationID *string       `json:"organization_id,omitempty" gorm:"size:36;index"`
	Organization   *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID;references:ID"`
	FavoriteEvents []Event       `json:"favorite_events,omitempty" gorm:"many2many:user_favorite_events;"`
}

// TableName especifica el nombre de tabla
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook para validación y hash de password
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.BaseModel.BeforeCreate(tx); err != nil {
		return err
	}

	// Validar datos requeridos
	if err := u.ValidateUser(); err != nil {
		return err
	}

	// Hash de la password si no está hasheada
	if !u.IsPasswordHashed() {
		if err := u.HashPassword(); err != nil {
			return err
		}
	}

	// Normalizar email
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	return nil
}

// BeforeUpdate hook para validaciones
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if err := u.BaseModel.BeforeUpdate(tx); err != nil {
		return err
	}

	// Solo validar campos que no sean password si está siendo actualizada
	if err := u.ValidateUserUpdate(); err != nil {
		return err
	}

	// Normalizar email si ha cambiado
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	return nil
}

// ValidateUser valida los datos del usuario en creación
func (u *User) ValidateUser() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(u.FirstName) == "" {
		return errors.New("first name is required")
	}

	if strings.TrimSpace(u.LastName) == "" {
		return errors.New("last name is required")
	}

	if len(strings.TrimSpace(u.Password)) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Validar role
	if !u.IsValidRole() {
		return errors.New("invalid role")
	}

	return nil
}

// ValidateUserUpdate valida datos en actualización (sin password)
func (u *User) ValidateUserUpdate() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(u.FirstName) == "" {
		return errors.New("first name is required")
	}

	if strings.TrimSpace(u.LastName) == "" {
		return errors.New("last name is required")
	}

	// Validar role
	if !u.IsValidRole() {
		return errors.New("invalid role")
	}

	return nil
}

// IsValidRole verifica si el role es válido
func (u *User) IsValidRole() bool {
	return u.Role == RoleAdmin || u.Role == RoleOrganizer || u.Role == RoleUser
}

// HashPassword hashea la password del usuario
func (u *User) HashPassword() error {
	if u.Password == "" {
		return errors.New("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedBytes)
	return nil
}

// CheckPassword verifica la password del usuario
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsPasswordHashed verifica si la password ya está hasheada (por el prefijo bcrypt)
func (u *User) IsPasswordHashed() bool {
	return strings.HasPrefix(u.Password, "$2a$") || strings.HasPrefix(u.Password, "$2b$") || strings.HasPrefix(u.Password, "$2y$")
}

// GetFullName retorna el nombre completo del usuario
func (u *User) GetFullName() string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

// HasPermission verifica si el usuario tiene un permiso específico
func (u *User) HasPermission(action string, resource string, resourceID string) bool {
	switch u.Role {
	case RoleAdmin:
		// Admin tiene todos los permisos
		return true

	case RoleOrganizer:
		// Organizer tiene permisos sobre su organización
		switch action {
		case "read":
			return true // Puede leer todo lo público
		case "create", "update", "delete":
			// Solo puede modificar contenido de su organización
			if resource == "event" || resource == "organization" {
				return u.OrganizationID != nil && *u.OrganizationID == resourceID
			}
			return false
		}

	case RoleUser:
		// Usuario solo puede leer y gestionar sus favoritos
		switch action {
		case "read":
			return true
		case "create", "update", "delete":
			return resource == "favorite" // Solo gestionar favoritos
		}
	}

	return false
}

// IsOrganizer verifica si el usuario es organizador
func (u *User) IsOrganizer() bool {
	return u.Role == RoleOrganizer
}

// IsAdmin verifica si el usuario es administrador
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// CanManageOrganization verifica si puede gestionar una organización específica
func (u *User) CanManageOrganization(organizationID string) bool {
	if u.IsAdmin() {
		return true
	}

	if u.IsOrganizer() && u.OrganizationID != nil {
		return *u.OrganizationID == organizationID
	}

	return false
}

// UpdateLastLogin actualiza el timestamp del último login
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// SetLocation establece la geolocalización del usuario
func (u *User) SetLocation(latitude, longitude float64, city, country string) {
	u.Latitude = &latitude
	u.Longitude = &longitude
	u.City = city
	u.Country = country
}

// HasLocation verifica si el usuario tiene coordenadas de geolocalización
func (u *User) HasLocation() bool {
	return u.Latitude != nil && u.Longitude != nil
}

// GetAuditData implementa AuditableModel
func (u *User) GetAuditData() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"email":      u.Email,
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"role":       u.Role,
		"is_active":  u.IsActive,
	}
}

// Métodos de base model implementados
func (u User) GetID() string           { return u.ID.String() }
func (u User) GetCreatedAt() time.Time { return u.CreatedAt }
func (u User) GetUpdatedAt() time.Time { return u.UpdatedAt }
