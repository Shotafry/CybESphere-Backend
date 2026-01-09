package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"cybesphere-backend/internal/config"

	"github.com/sirupsen/logrus"
)

// Logger instancia global del logger
var Logger *logrus.Logger

// Init inicializa el sistema de logging
func Init(cfg *config.LoggingConfig) error {
	Logger = logrus.New()

	// Configurar nivel de log
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// Configurar formato
	if strings.ToLower(cfg.Format) == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Configurar output
	var output io.Writer
	switch strings.ToLower(cfg.Output) {
	case "stderr":
		output = os.Stderr
	case "file":
		if cfg.FilePath == "" {
			cfg.FilePath = "logs/app.log"
		}

		// Crear directorio si no existe
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0750); err != nil {
			return err
		}

		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
		output = file
	default:
		output = os.Stdout
	}

	Logger.SetOutput(output)

	return nil
}

// GetLogger retorna la instancia del logger
func GetLogger() *logrus.Logger {
	if Logger == nil {
		// Logger por defecto si no se ha inicializado
		Logger = logrus.New()
		Logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return Logger
}

// WithFields crea un logger con campos adicionales
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// WithField crea un logger con un campo adicional
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// Debug log a nivel debug
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf log a nivel debug con formato
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info log a nivel info
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof log a nivel info con formato
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn log a nivel warning
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf log a nivel warning con formato
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error log a nivel error
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf log a nivel error con formato
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal log a nivel fatal (termina la aplicación)
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf log a nivel fatal con formato (termina la aplicación)
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// Panic log a nivel panic
func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

// Panicf log a nivel panic con formato
func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}

// LogRequest log para requests HTTP
func LogRequest(method, uri, userAgent, ip string, statusCode int, duration int64) {
	WithFields(logrus.Fields{
		"method":      method,
		"uri":         uri,
		"status_code": statusCode,
		"duration_ms": duration,
		"user_agent":  userAgent,
		"ip":          ip,
		"type":        "http_request",
	}).Info("HTTP Request")
}

// LogDBQuery log para queries de base de datos
func LogDBQuery(query string, duration int64, rows int64) {
	WithFields(logrus.Fields{
		"query":       query,
		"duration_ms": duration,
		"rows":        rows,
		"type":        "db_query",
	}).Debug("Database Query")
}

// LogAuth log para eventos de autenticación
func LogAuth(userID string, action string, success bool, reason string) {
	fields := logrus.Fields{
		"user_id": userID,
		"action":  action,
		"success": success,
		"type":    "auth",
	}

	if reason != "" {
		fields["reason"] = reason
	}

	WithFields(fields).Info("Authentication Event")
}

// LogAudit log para eventos de auditoría
func LogAudit(userID, action, resource, resourceID string, changes map[string]interface{}) {
	fields := logrus.Fields{
		"user_id":     userID,
		"action":      action,
		"resource":    resource,
		"resource_id": resourceID,
		"type":        "audit",
	}

	if changes != nil {
		fields["changes"] = changes
	}

	WithFields(fields).Info("Audit Event")
}
