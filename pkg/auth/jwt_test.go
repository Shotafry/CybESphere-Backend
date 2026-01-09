package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Constantes para testing
const (
	testSecretKey       = "test-secret-key-32-characters-long-minimum"
	testShortSecretKey  = "too-short"
	testUserID          = "123e4567-e89b-12d3-a456-426614174000"
	testEmail           = "test@example.com"
	testRole            = "user"
	testIssuer          = "cybesphere-test"
	testAccessDuration  = 15 * time.Minute
	testRefreshDuration = 7 * 24 * time.Hour
)

// createTestJWTManager crea un JWTManager para testing
func createTestJWTManager() *JWTManager {
	manager, _ := NewJWTManager(testSecretKey, testAccessDuration, testRefreshDuration, testIssuer)
	return manager
}

// TestNewJWTManager tests para crear instancia de JWTManager
func TestNewJWTManager(t *testing.T) {
	tests := []struct {
		name            string
		secretKey       string
		accessDuration  time.Duration
		refreshDuration time.Duration
		issuer          string
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "configuración válida",
			secretKey:       testSecretKey,
			accessDuration:  testAccessDuration,
			refreshDuration: testRefreshDuration,
			issuer:          testIssuer,
			wantErr:         false,
		},
		{
			name:            "secret key muy corta",
			secretKey:       testShortSecretKey,
			accessDuration:  testAccessDuration,
			refreshDuration: testRefreshDuration,
			issuer:          testIssuer,
			wantErr:         true,
			errMsg:          "JWT secret key must be at least 32 characters long",
		},
		{
			name:            "configuración mínima",
			secretKey:       "12345678901234567890123456789012", // Exactamente 32 caracteres
			accessDuration:  time.Minute,
			refreshDuration: time.Hour,
			issuer:          "test",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewJWTManager(tt.secretKey, tt.accessDuration, tt.refreshDuration, tt.issuer)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, manager)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, manager)
				assert.Equal(t, tt.secretKey, manager.secretKey)
				assert.Equal(t, tt.accessDuration, manager.accessTokenDuration)
				assert.Equal(t, tt.refreshDuration, manager.refreshTokenDuration)
				assert.Equal(t, tt.issuer, manager.issuer)
			}
		})
	}
}

// TestGenerateTokenPair tests para generación de pares de tokens
func TestGenerateTokenPair(t *testing.T) {
	manager := createTestJWTManager()

	tests := []struct {
		name    string
		userID  string
		email   string
		role    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "datos válidos",
			userID:  testUserID,
			email:   testEmail,
			role:    testRole,
			wantErr: false,
		},
		{
			name:    "userID vacío",
			userID:  "",
			email:   testEmail,
			role:    testRole,
			wantErr: true,
			errMsg:  "userID, email and role are required",
		},
		{
			name:    "email vacío",
			userID:  testUserID,
			email:   "",
			role:    testRole,
			wantErr: true,
			errMsg:  "userID, email and role are required",
		},
		{
			name:    "role vacío",
			userID:  testUserID,
			email:   testEmail,
			role:    "",
			wantErr: true,
			errMsg:  "userID, email and role are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenPair, err := manager.GenerateTokenPair(tt.userID, tt.email, tt.role)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, tokenPair)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokenPair)

				// Verificar estructura del TokenPair
				assert.NotEmpty(t, tokenPair.AccessToken)
				assert.NotEmpty(t, tokenPair.RefreshToken)
				assert.Equal(t, "Bearer", tokenPair.TokenType)
				assert.True(t, tokenPair.AccessTokenExpiresAt.After(time.Now()))
				assert.True(t, tokenPair.RefreshTokenExpiresAt.After(time.Now()))
				assert.True(t, tokenPair.RefreshTokenExpiresAt.After(tokenPair.AccessTokenExpiresAt))

				// Verificar que los tokens son diferentes
				assert.NotEqual(t, tokenPair.AccessToken, tokenPair.RefreshToken)
			}
		})
	}
}

// TestValidateAccessToken tests para validación de access tokens
func TestValidateAccessToken(t *testing.T) {
	manager := createTestJWTManager()

	// Generar un token válido para testing
	tokenPair, err := manager.GenerateTokenPair(testUserID, testEmail, testRole)
	require.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "access token válido",
			token:   tokenPair.AccessToken,
			wantErr: false,
		},
		{
			name:        "token vacío",
			token:       "",
			wantErr:     true,
			expectedErr: ErrInvalidToken,
		},
		{
			name:        "token malformado",
			token:       "invalid.token.format",
			wantErr:     true,
			expectedErr: ErrInvalidToken,
		},
		{
			name:        "refresh token como access token",
			token:       tokenPair.RefreshToken,
			wantErr:     true,
			expectedErr: ErrInvalidTokenType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := manager.ValidateAccessToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, testUserID, claims.UserID)
				assert.Equal(t, testEmail, claims.Email)
				assert.Equal(t, testRole, claims.Role)
				assert.Equal(t, AccessToken, claims.Type)
				assert.NotEmpty(t, claims.TokenID)
			}
		})
	}
}

// TestValidateRefreshToken tests para validación de refresh tokens
func TestValidateRefreshToken(t *testing.T) {
	manager := createTestJWTManager()

	// Generar un token válido para testing
	tokenPair, err := manager.GenerateTokenPair(testUserID, testEmail, testRole)
	require.NoError(t, err)

	tests := []struct {
		name        string
		token       string
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "refresh token válido",
			token:   tokenPair.RefreshToken,
			wantErr: false,
		},
		{
			name:        "access token como refresh token",
			token:       tokenPair.AccessToken,
			wantErr:     true,
			expectedErr: ErrInvalidTokenType,
		},
		{
			name:        "token inválido",
			token:       "invalid-token",
			wantErr:     true,
			expectedErr: ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := manager.ValidateRefreshToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, testUserID, claims.UserID)
				assert.Equal(t, testEmail, claims.Email)
				assert.Equal(t, testRole, claims.Role)
				assert.Equal(t, RefreshToken, claims.Type)
			}
		})
	}
}

// TestRefreshTokens tests para refrescar tokens
func TestRefreshTokens(t *testing.T) {
	manager := createTestJWTManager()

	// Generar tokens iniciales
	originalTokenPair, err := manager.GenerateTokenPair(testUserID, testEmail, testRole)
	require.NoError(t, err)

	t.Run("refresh con token válido", func(t *testing.T) {
		newTokenPair, err := manager.RefreshTokens(originalTokenPair.RefreshToken)

		assert.NoError(t, err)
		assert.NotNil(t, newTokenPair)

		// Verificar que se generaron nuevos tokens
		assert.NotEqual(t, originalTokenPair.AccessToken, newTokenPair.AccessToken)
		assert.NotEqual(t, originalTokenPair.RefreshToken, newTokenPair.RefreshToken)

		// Verificar que los nuevos tokens son válidos
		claims, err := manager.ValidateAccessToken(newTokenPair.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, testUserID, claims.UserID)
	})

	t.Run("refresh con access token", func(t *testing.T) {
		_, err := manager.RefreshTokens(originalTokenPair.AccessToken)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid refresh token")
	})

	t.Run("refresh con token inválido", func(t *testing.T) {
		_, err := manager.RefreshTokens("invalid-token")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid refresh token")
	})
}

// TestExtractTokenFromHeader tests para extracción de tokens de headers
func TestExtractTokenFromHeader(t *testing.T) {
	manager := createTestJWTManager()

	tests := []struct {
		name       string
		authHeader string
		wantToken  string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "header válido",
			authHeader: "Bearer abc123def456",
			wantToken:  "abc123def456",
			wantErr:    false,
		},
		{
			name:       "header vacío",
			authHeader: "",
			wantErr:    true,
			errMsg:     "authorization header is required",
		},
		{
			name:       "sin Bearer prefix",
			authHeader: "abc123def456",
			wantErr:    true,
			errMsg:     "authorization header must start with 'Bearer '",
		},
		{
			name:       "Bearer sin token",
			authHeader: "Bearer ",
			wantErr:    true,
			errMsg:     "token is required",
		},
		{
			name:       "Bearer con espacios extra",
			authHeader: "Bearer   token123",
			wantToken:  "  token123",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.ExtractTokenFromHeader(tt.authHeader)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}

// TestTokenExpiration tests para manejo de tokens expirados
func TestTokenExpiration(t *testing.T) {
	// Crear manager con tokens de muy corta duración
	shortManager, err := NewJWTManager(testSecretKey, 1*time.Millisecond, 2*time.Millisecond, testIssuer)
	require.NoError(t, err)

	// Generar token
	tokenPair, err := shortManager.GenerateTokenPair(testUserID, testEmail, testRole)
	require.NoError(t, err)

	// Esperar a que expire
	time.Sleep(10 * time.Millisecond)

	// Intentar validar token expirado
	_, err = shortManager.ValidateAccessToken(tokenPair.AccessToken)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrExpiredToken)

	// Verificar IsTokenExpired
	assert.True(t, shortManager.IsTokenExpired(tokenPair.AccessToken))
}

// TestGetTokenClaims tests para extracción de claims sin validación
func TestGetTokenClaims(t *testing.T) {
	manager := createTestJWTManager()

	// Generar token
	tokenPair, err := manager.GenerateTokenPair(testUserID, testEmail, testRole)
	require.NoError(t, err)

	t.Run("extraer claims de token válido", func(t *testing.T) {
		claims, err := manager.GetTokenClaims(tokenPair.AccessToken)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, testUserID, claims.UserID)
		assert.Equal(t, testEmail, claims.Email)
		assert.Equal(t, testRole, claims.Role)
		assert.Equal(t, AccessToken, claims.Type)
	})

	t.Run("extraer claims de token inválido", func(t *testing.T) {
		_, err := manager.GetTokenClaims("invalid-token")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})
}

// TestTokenIDUniqueness tests para verificar unicidad de TokenID
func TestTokenIDUniqueness(t *testing.T) {
	manager := createTestJWTManager()

	// Generar múltiples pares de tokens
	tokenPairs := make([]*TokenPair, 5)
	for i := 0; i < 5; i++ {
		pair, err := manager.GenerateTokenPair(testUserID, testEmail, testRole)
		require.NoError(t, err)
		tokenPairs[i] = pair
	}

	// Verificar que todos los TokenID son únicos
	tokenIDs := make(map[string]bool)
	for _, pair := range tokenPairs {
		accessClaims, err := manager.GetTokenClaims(pair.AccessToken)
		require.NoError(t, err)

		refreshClaims, err := manager.GetTokenClaims(pair.RefreshToken)
		require.NoError(t, err)

		// Verificar unicidad de access token ID
		assert.False(t, tokenIDs[accessClaims.TokenID], "TokenID duplicado para access token")
		tokenIDs[accessClaims.TokenID] = true

		// Verificar unicidad de refresh token ID
		assert.False(t, tokenIDs[refreshClaims.TokenID], "TokenID duplicado para refresh token")
		tokenIDs[refreshClaims.TokenID] = true

		// Verificar que access y refresh token tienen IDs diferentes
		assert.NotEqual(t, accessClaims.TokenID, refreshClaims.TokenID)
	}
}

// TestGenerateSecureRandomString tests para generación de strings aleatorios
func TestGenerateSecureRandomString(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "longitud normal",
			length:  32,
			wantErr: false,
		},
		{
			name:    "longitud mínima",
			length:  1,
			wantErr: false,
		},
		{
			name:    "longitud cero",
			length:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateSecureRandomString(tt.length)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// El resultado debe ser hexadecimal, así que será el doble de longitud
				assert.Equal(t, tt.length*2, len(result))

				// Verificar que contiene solo caracteres hexadecimales
				for _, char := range result {
					assert.True(t, (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f'))
				}
			}
		})
	}

	// Verificar unicidad generando múltiples strings
	strings := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result, err := GenerateSecureRandomString(16)
		require.NoError(t, err)
		assert.False(t, strings[result], "String aleatorio duplicado")
		strings[result] = true
	}
}

// TestHashRefreshToken tests para hashing de refresh tokens
func TestHashRefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "token válido",
			token:   "valid-refresh-token-123",
			wantErr: false,
		},
		{
			name:    "token vacío",
			token:   "",
			wantErr: true,
			errMsg:  "token cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashRefreshToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tt.token, hash)

				// Verificar que el mismo token produce el mismo hash
				hash2, err := HashRefreshToken(tt.token)
				assert.NoError(t, err)
				assert.Equal(t, hash, hash2)
			}
		})
	}
}

// TestIntegrationFlow test de integración del flujo completo
func TestIntegrationFlow(t *testing.T) {
	manager := createTestJWTManager()

	// 1. Generar tokens iniciales
	originalPair, err := manager.GenerateTokenPair(testUserID, testEmail, testRole)
	require.NoError(t, err)

	// 2. Validar access token
	accessClaims, err := manager.ValidateAccessToken(originalPair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, testUserID, accessClaims.UserID)

	// 3. Validar refresh token
	refreshClaims, err := manager.ValidateRefreshToken(originalPair.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, testUserID, refreshClaims.UserID)

	// 4. Refrescar tokens
	newPair, err := manager.RefreshTokens(originalPair.RefreshToken)
	require.NoError(t, err)

	// 5. Verificar que los nuevos tokens son diferentes
	assert.NotEqual(t, originalPair.AccessToken, newPair.AccessToken)
	assert.NotEqual(t, originalPair.RefreshToken, newPair.RefreshToken)

	// 6. Validar nuevos tokens
	newAccessClaims, err := manager.ValidateAccessToken(newPair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, testUserID, newAccessClaims.UserID)

	newRefreshClaims, err := manager.ValidateRefreshToken(newPair.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, testUserID, newRefreshClaims.UserID)

	// 7. Verificar que los TokenID son diferentes
	assert.NotEqual(t, accessClaims.TokenID, newAccessClaims.TokenID)
	assert.NotEqual(t, refreshClaims.TokenID, newRefreshClaims.TokenID)
}
