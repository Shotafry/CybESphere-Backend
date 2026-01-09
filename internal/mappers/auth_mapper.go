package mappers

import (
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/models"
	"fmt"
	"strings"
	"time"
)

// AuthMapperImpl implementación del mapper de autenticación
type AuthMapperImpl struct{}

// NewAuthMapper crea nueva instancia del mapper
func NewAuthMapper() AuthMapperImpl {
	return AuthMapperImpl{}
}

// RegisterRequestToUser convierte RegisterRequest a modelo User
func (m AuthMapperImpl) RegisterRequestToUser(req *dto.RegisterRequest) (*models.User, error) {
	user := &models.User{
		// Información básica requerida
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Password:  req.Password, // Se hasheará automáticamente en BeforeCreate
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),

		// Información profesional opcional
		Company:  strings.TrimSpace(req.Company),
		Position: strings.TrimSpace(req.Position),
		City:     strings.TrimSpace(req.City),
		Country:  strings.TrimSpace(req.Country),

		// Configuraciones por defecto para España
		Role:              models.RoleUser,
		IsActive:          true,
		IsVerified:        false, // Requiere verificación por email
		Timezone:          "Europe/Madrid",
		Language:          "es",
		NewsletterEnabled: true, // Opt-in por defecto
	}

	return user, nil
}

// UserToAuthResponse convierte User y tokens a AuthResponse
func (m AuthMapperImpl) UserToAuthResponse(user *models.User, accessToken, refreshToken string, expiresIn int) dto.AuthResponse {
	return dto.AuthResponse{
		User:         m.userToResponseForAuth(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}
}

// UserToRegisterResponse convierte User y tokens a RegisterResponse
func (m AuthMapperImpl) UserToRegisterResponse(user *models.User, accessToken, refreshToken string, expiresIn int) dto.RegisterResponse {
	response := dto.RegisterResponse{
		Message:              "Usuario registrado exitosamente",
		User:                 m.userToResponseForAuth(user),
		RequiresVerification: !user.IsVerified,
	}

	// Solo incluir tokens si el usuario está activo
	if user.IsActive {
		response.AccessToken = accessToken
		response.RefreshToken = refreshToken
		response.TokenType = "Bearer"
		response.ExpiresIn = expiresIn
	}

	// Ajustar mensaje según estado
	if !user.IsVerified {
		response.Message = "Usuario registrado exitosamente. Se ha enviado un email de verificación."
	}

	return response
}

// BuildPasswordResetResponse crea respuesta para reset de contraseña
func (m AuthMapperImpl) BuildPasswordResetResponse(success bool, expiresInMinutes int) dto.PasswordResetResponse {
	response := dto.PasswordResetResponse{
		Success: success,
	}

	if success {
		response.Message = "Se ha enviado un enlace de recuperación a tu email"
		response.ExpiresIn = expiresInMinutes
	} else {
		response.Message = "No se pudo procesar la solicitud de recuperación"
	}

	return response
}

// BuildEmailVerificationResponse crea respuesta para verificación de email
func (m AuthMapperImpl) BuildEmailVerificationResponse(success bool, isVerified bool) dto.EmailVerificationResponse {
	response := dto.EmailVerificationResponse{
		Success:  success,
		Verified: isVerified,
	}

	if success && isVerified {
		response.Message = "Email verificado exitosamente"
	} else if success && !isVerified {
		response.Message = "Enlace de verificación enviado"
	} else {
		response.Message = "No se pudo verificar el email"
	}

	return response
}

// BuildLogoutResponse crea respuesta para logout
func (m AuthMapperImpl) BuildLogoutResponse(success bool) dto.LogoutResponse {
	response := dto.LogoutResponse{
		Success: success,
	}

	if success {
		response.Message = "Sesión cerrada exitosamente"
	} else {
		response.Message = "Error al cerrar sesión"
	}

	return response
}

// BuildTokenResponse crea respuesta para renovación de tokens
func (m AuthMapperImpl) BuildTokenResponse(accessToken, refreshToken string, expiresAt time.Time) dto.TokenResponse {
	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(time.Until(expiresAt).Seconds()),
		ExpiresAt:    expiresAt,
	}
}

// RefreshTokensToSessionList convierte refresh tokens a lista de sesiones - IMPLEMENTACIÓN COMPLETA
func (m AuthMapperImpl) RefreshTokensToSessionList(tokens []*models.RefreshToken, currentTokenID string) dto.SessionListResponse {
	sessions := make([]dto.SessionResponse, 0, len(tokens))
	activeCount := 0

	for _, token := range tokens {
		if token.IsActive() {
			activeCount++
		}

		sessionResponse := dto.SessionResponse{
			ID:         token.ID.String(),
			TokenID:    token.TokenID,
			DeviceInfo: m.determineDeviceInfo(token.UserAgent, token.DeviceInfo),
			IPAddress:  token.IPAddress,
			UserAgent:  m.sanitizeUserAgent(token.UserAgent),
			CreatedAt:  token.CreatedAt,
			LastUsedAt: token.LastUsedAt,
			ExpiresAt:  token.ExpiresAt,
			IsCurrent:  token.TokenID == currentTokenID,
		}

		sessions = append(sessions, sessionResponse)
	}

	return dto.SessionListResponse{
		Sessions:       sessions,
		TotalSessions:  len(sessions),
		ActiveSessions: activeCount,
	}
}

// =============================================================================
// MÉTODOS HELPER PRIVADOS - IMPLEMENTACIONES COMPLETAS
// =============================================================================

// userToResponseForAuth convierte User a UserResponse específico para auth
func (m AuthMapperImpl) userToResponseForAuth(user *models.User) dto.UserResponse {
	response := dto.UserResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		FullName:    user.GetFullName(),
		Role:        string(user.Role),
		Company:     user.Company,
		Position:    user.Position,
		City:        user.City,
		Country:     user.Country,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		Timezone:    user.Timezone,
		Language:    user.Language,
		CreatedAt:   user.CreatedAt,
		LastLoginAt: user.LastLoginAt,
	}

	// Incluir organización si existe
	if user.Organization != nil {
		response.Organization = &dto.OrganizationSummaryResponse{
			ID:          user.Organization.ID.String(),
			Slug:        user.Organization.Slug,
			Name:        user.Organization.Name,
			LogoURL:     user.Organization.LogoURL,
			IsVerified:  user.Organization.IsVerified,
			EventsCount: user.Organization.EventsCount,
			City:        user.Organization.City,
			Country:     user.Organization.Country,
		}
	}

	return response
}

// determineDeviceInfo determina información del dispositivo desde User-Agent y DeviceInfo
func (m AuthMapperImpl) determineDeviceInfo(userAgent, deviceInfo string) string {
	// Si ya tenemos deviceInfo específico, usarlo
	if deviceInfo != "" {
		return deviceInfo
	}

	// Si no, extraer del User-Agent
	return m.extractDeviceInfo(userAgent)
}

// extractDeviceInfo extrae información del dispositivo del User-Agent - IMPLEMENTACIÓN COMPLETA
func (m AuthMapperImpl) extractDeviceInfo(userAgent string) string {
	if userAgent == "" {
		return "Dispositivo Desconocido"
	}

	ua := strings.ToLower(userAgent)

	// Detectar sistemas operativos móviles
	if strings.Contains(ua, "android") {
		if strings.Contains(ua, "mobile") {
			return "Android Móvil"
		}
		return "Android Tablet"
	}

	if strings.Contains(ua, "iphone") {
		return "iPhone"
	}

	if strings.Contains(ua, "ipad") {
		return "iPad"
	}

	// Detectar navegadores desktop
	if strings.Contains(ua, "windows") {
		browser := m.extractBrowser(ua)
		return fmt.Sprintf("Windows - %s", browser)
	}

	if strings.Contains(ua, "macintosh") || strings.Contains(ua, "mac os") {
		browser := m.extractBrowser(ua)
		return fmt.Sprintf("macOS - %s", browser)
	}

	if strings.Contains(ua, "linux") {
		browser := m.extractBrowser(ua)
		return fmt.Sprintf("Linux - %s", browser)
	}

	// Fallback genérico
	return m.extractBrowser(ua)
}

// extractBrowser extrae el navegador del User-Agent - NUEVA FUNCIÓN
func (m AuthMapperImpl) extractBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)

	if strings.Contains(ua, "edg/") {
		return "Microsoft Edge"
	}
	if strings.Contains(ua, "chrome/") && !strings.Contains(ua, "edg/") {
		return "Google Chrome"
	}
	if strings.Contains(ua, "firefox/") {
		return "Mozilla Firefox"
	}
	if strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/") {
		return "Safari"
	}
	if strings.Contains(ua, "opera/") || strings.Contains(ua, "opr/") {
		return "Opera"
	}

	return "Navegador Desconocido"
}

// sanitizeUserAgent limpia y trunca User-Agent para almacenamiento
func (m AuthMapperImpl) sanitizeUserAgent(userAgent string) string {
	if userAgent == "" {
		return ""
	}

	// Truncar User-Agent muy largos
	maxLength := 500
	if len(userAgent) > maxLength {
		return userAgent[:maxLength] + "..."
	}
	return userAgent
}
