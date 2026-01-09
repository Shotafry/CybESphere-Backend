package dto

import (
	"time"
)

// AuthResponse DTO de respuesta para autenticación
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"` // segundos
}

// TokenResponse DTO de respuesta para tokens
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RegisterResponse DTO de respuesta para registro
type RegisterResponse struct {
	Message              string       `json:"message"`
	User                 UserResponse `json:"user"`
	AccessToken          string       `json:"access_token,omitempty"`
	RefreshToken         string       `json:"refresh_token,omitempty"`
	TokenType            string       `json:"token_type,omitempty"`
	ExpiresIn            int          `json:"expires_in,omitempty"`
	RequiresVerification bool         `json:"requires_verification"`
}

// PasswordResetResponse DTO de respuesta para reset de contraseña
type PasswordResetResponse struct {
	Message   string `json:"message"`
	Success   bool   `json:"success"`
	ExpiresIn int    `json:"expires_in,omitempty"` // minutos
}

// EmailVerificationResponse DTO de respuesta para verificación de email
type EmailVerificationResponse struct {
	Message  string `json:"message"`
	Success  bool   `json:"success"`
	Verified bool   `json:"verified"`
}

// LogoutResponse DTO de respuesta para logout
type LogoutResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// SessionListResponse lista de sesiones del usuario
type SessionListResponse struct {
	Sessions       []SessionResponse `json:"sessions"`
	TotalSessions  int               `json:"total_sessions"`
	ActiveSessions int               `json:"active_sessions"`
}
