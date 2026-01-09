package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel implementa el patrón de auditoría automática
// Todos los modelos del dominio deben embeber este struct
type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Campos de auditoría
	CreatedBy string `json:"created_by" gorm:"size:100"`
	UpdatedBy string `json:"updated_by" gorm:"size:100"`
}

// BeforeCreate hook de GORM para población automática
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	// Obtener usuario del contexto si está disponible
	if userID := getUserFromContext(tx); userID != "" {
		b.CreatedBy = userID
		b.UpdatedBy = userID
	}

	return nil
}

// BeforeUpdate hook de GORM para auditoría de actualizaciones
func (b *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	// Obtener usuario del contexto si está disponible
	if userID := getUserFromContext(tx); userID != "" {
		b.UpdatedBy = userID
	}

	return nil
}

// getUserFromContext extrae el ID del usuario del contexto de la transacción
func getUserFromContext(tx *gorm.DB) string {
	if userCtx := tx.Statement.Context.Value("user_id"); userCtx != nil {
		if userID, ok := userCtx.(string); ok {
			return userID
		}
	}
	return ""
}

// AuditableModel interfaz para modelos que requieren auditoría extendida
type AuditableModel interface {
	GetAuditData() map[string]interface{}
	SetAuditContext(userID, action string)
}

// AuditLog modelo para registro detallado de cambios
type AuditLog struct {
	BaseModel
	UserID     string         `json:"user_id" gorm:"not null;size:100;index"`
	Action     string         `json:"action" gorm:"not null;size:50;index"`
	Resource   string         `json:"resource" gorm:"not null;size:100;index"`
	ResourceID string         `json:"resource_id" gorm:"not null;size:100;index"`
	Changes    map[string]any `json:"changes" gorm:"type:jsonb"`
	IPAddress  string         `json:"ip_address" gorm:"size:45"`
	UserAgent  string         `json:"user_agent" gorm:"size:500"`
	Timestamp  time.Time      `json:"timestamp" gorm:"not null;index"`
	Status     int            `json:"status" gorm:"not null;default:0"`
}

// TableName especifica el nombre de tabla para AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// LogAuditEvent registra un evento de auditoría
func LogAuditEvent(db *gorm.DB, userID, action, resource, resourceID string, changes map[string]interface{}, ipAddress, userAgent string) error {
	auditLog := AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Changes:    changes,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Timestamp:  time.Now(),
	}

	return db.Create(&auditLog).Error
}

// GetID implementa BaseEntity
func (b *BaseModel) GetID() string {
	return b.ID.String()
}

// GetCreatedAt implementa BaseEntity
func (b *BaseModel) GetCreatedAt() time.Time {
	return b.CreatedAt
}

// GetUpdatedAt implementa BaseEntity
func (b *BaseModel) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}
