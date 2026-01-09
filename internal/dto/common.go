package dto

import "time"

// AppliedFilters filtros aplicados (para mostrar al usuario)
type AppliedFilters struct {
	Search   string                 `json:"search,omitempty"`
	OrderBy  string                 `json:"order_by,omitempty"`
	OrderDir string                 `json:"order_dir,omitempty"`
	Filters  map[string]interface{} `json:"filters,omitempty"`
}

// ErrorResponse respuesta de error estándar
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Field   string                 `json:"field,omitempty"`
}

// ValidationErrorResponse errores de validación
type ValidationErrorResponse struct {
	Error   string              `json:"error"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

// SuccessResponse respuesta de éxito genérica
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// DeleteResponse respuesta para eliminación
type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

// BulkOperationResponse respuesta para operaciones en lote
type BulkOperationResponse struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	Processed  int      `json:"processed"`
	Failed     int      `json:"failed"`
	Errors     []string `json:"errors,omitempty"`
	SuccessIDs []string `json:"success_ids,omitempty"`
	FailedIDs  []string `json:"failed_ids,omitempty"`
}

// FileUploadResponse respuesta para carga de archivos
type FileUploadResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}

// HealthCheckResponse respuesta para health check
type HealthCheckResponse struct {
	Status    string                   `json:"status"` // healthy, degraded, unhealthy
	Version   string                   `json:"version"`
	Timestamp time.Time                `json:"timestamp"`
	Services  map[string]ServiceHealth `json:"services,omitempty"`
}

// ServiceHealth salud de un servicio
type ServiceHealth struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency int64  `json:"latency_ms,omitempty"`
}

// SearchResponse respuesta de búsqueda genérica
type SearchResponse struct {
	Query   string      `json:"query"`
	Results interface{} `json:"results"`
	Total   int         `json:"total"`
	Took    int64       `json:"took_ms"` // tiempo de búsqueda en ms
}

// StatsResponse respuesta de estadísticas
type StatsResponse struct {
	Period string                 `json:"period"` // daily, weekly, monthly, yearly
	From   time.Time              `json:"from"`
	To     time.Time              `json:"to"`
	Data   map[string]interface{} `json:"data"`
}

// ExportResponse respuesta de exportación
type ExportResponse struct {
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
	FileURL     string    `json:"file_url,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	RecordCount int       `json:"record_count"`
	Format      string    `json:"format"` // csv, excel, pdf, json
}

// ImportResponse respuesta de importación
type ImportResponse struct {
	Success  bool          `json:"success"`
	Message  string        `json:"message"`
	Imported int           `json:"imported"`
	Skipped  int           `json:"skipped"`
	Failed   int           `json:"failed"`
	Errors   []ImportError `json:"errors,omitempty"`
	Warnings []string      `json:"warnings,omitempty"`
}

// ImportError error de importación
type ImportError struct {
	Row     int    `json:"row"`
	Column  string `json:"column,omitempty"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

// NotificationResponse respuesta de notificación
type NotificationResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Read      bool                   `json:"read"`
	CreatedAt time.Time              `json:"created_at"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// WebSocketMessage mensaje para WebSocket
type WebSocketMessage struct {
	Type     string                 `json:"type"`
	Event    string                 `json:"event"`
	Payload  interface{}            `json:"payload"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ActivityLogResponse registro de actividad
type ActivityLogResponse struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
}
