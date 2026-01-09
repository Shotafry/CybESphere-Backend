package dto

// LoginRequest DTO para login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest DTO para registro
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=128"`
	FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string `json:"last_name" binding:"required,min=1,max=100"`
	Company   string `json:"company" binding:"max=200"`
	Position  string `json:"position" binding:"max=200"`
	City      string `json:"city" binding:"max=100"`
	Country   string `json:"country" binding:"max=100"`
}

// RefreshTokenRequest DTO para renovar token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest DTO para solicitar reset de contraseña
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest DTO para resetear contraseña
type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// VerifyEmailRequest DTO para verificar email
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// ResendVerificationRequest DTO para reenviar verificación
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// LogoutRequest DTO para logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutAllRequest DTO para logout de todas las sesiones
type LogoutAllRequest struct {
	ConfirmLogoutAll bool `json:"confirm_logout_all" binding:"required"`
}
