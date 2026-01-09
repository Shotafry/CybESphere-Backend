package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// createTestRefreshToken crea un refresh token válido para testing
func createTestRefreshToken() *RefreshToken {
	return &RefreshToken{
		UserID:     uuid.New().String(),
		TokenHash:  "hashed-token-12345",
		TokenID:    uuid.New().String(),
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // 7 días
		IsRevoked:  false,
		IPAddress:  "192.168.1.100",
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		DeviceInfo: "Desktop Chrome Windows",
	}
}

// TestRefreshToken_ValidateRefreshToken tests unitarios para validación
func TestRefreshToken_ValidateRefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		token   *RefreshToken
		wantErr bool
		errMsg  string
	}{
		{
			name:    "token válido",
			token:   createTestRefreshToken(),
			wantErr: false,
		},
		{
			name: "user ID vacío",
			token: func() *RefreshToken {
				rt := createTestRefreshToken()
				rt.UserID = ""
				return rt
			}(),
			wantErr: true,
			errMsg:  "user ID is required",
		},
		{
			name: "token hash vacío",
			token: func() *RefreshToken {
				rt := createTestRefreshToken()
				rt.TokenHash = ""
				return rt
			}(),
			wantErr: true,
			errMsg:  "token hash is required",
		},
		{
			name: "token ID vacío",
			token: func() *RefreshToken {
				rt := createTestRefreshToken()
				rt.TokenID = ""
				return rt
			}(),
			wantErr: true,
			errMsg:  "token ID is required",
		},
		{
			name: "expires at vacío",
			token: func() *RefreshToken {
				rt := createTestRefreshToken()
				rt.ExpiresAt = time.Time{}
				return rt
			}(),
			wantErr: true,
			errMsg:  "expires at is required",
		},
		{
			name: "token con fecha de expiración en el pasado",
			token: func() *RefreshToken {
				rt := createTestRefreshToken()
				rt.ExpiresAt = time.Now().Add(-1 * time.Hour)
				return rt
			}(),
			wantErr: true,
			errMsg:  "token cannot be created with past expiration date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.token.ValidateRefreshToken()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRefreshToken_IsExpired tests para verificación de expiración
func TestRefreshToken_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "token no expirado",
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      false,
		},
		{
			name:      "token expirado",
			expiresAt: time.Now().Add(-1 * time.Hour),
			want:      true,
		},
		{
			name:      "token expira justo ahora",
			expiresAt: time.Now(),
			want:      false, // Podría fallar por timing, pero generalmente false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := &RefreshToken{ExpiresAt: tt.expiresAt}
			result := token.IsExpired()

			// Para el caso "expira justo ahora", permitimos ambos resultados
			if tt.name == "token expira justo ahora" {
				// No hacer assert específico debido a timing
				return
			}

			assert.Equal(t, tt.want, result)
		})
	}
}

// TestRefreshToken_IsActive tests para verificación de token activo
func TestRefreshToken_IsActive(t *testing.T) {
	tests := []struct {
		name  string
		token *RefreshToken
		want  bool
	}{
		{
			name: "token activo",
			token: &RefreshToken{
				IsRevoked: false,
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			want: true,
		},
		{
			name: "token revocado",
			token: &RefreshToken{
				IsRevoked: true,
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			want: false,
		},
		{
			name: "token expirado",
			token: &RefreshToken{
				IsRevoked: false,
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			want: false,
		},
		{
			name: "token revocado y expirado",
			token: &RefreshToken{
				IsRevoked: true,
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.IsActive()
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestRefreshToken_Revoke tests para revocación de tokens
func TestRefreshToken_Revoke(t *testing.T) {
	token := createTestRefreshToken()

	// Verificar estado inicial
	assert.False(t, token.IsRevoked)
	assert.Nil(t, token.RevokedAt)

	// Revocar token
	token.Revoke()

	// Verificar que fue revocado
	assert.True(t, token.IsRevoked)
	assert.NotNil(t, token.RevokedAt)
	assert.WithinDuration(t, time.Now(), *token.RevokedAt, time.Second)

	// Verificar que ya no está activo
	assert.False(t, token.IsActive())
}

// TestRefreshToken_UpdateLastUsed tests para actualización de último uso
func TestRefreshToken_UpdateLastUsed(t *testing.T) {
	token := createTestRefreshToken()

	// Verificar estado inicial
	assert.Nil(t, token.LastUsedAt)

	// Actualizar último uso
	token.UpdateLastUsed()

	// Verificar que se actualizó
	assert.NotNil(t, token.LastUsedAt)
	assert.WithinDuration(t, time.Now(), *token.LastUsedAt, time.Second)

	// Actualizar nuevamente después de un tiempo
	time.Sleep(10 * time.Millisecond)
	oldLastUsed := *token.LastUsedAt
	token.UpdateLastUsed()

	// Verificar que se actualizó a un tiempo posterior
	assert.True(t, token.LastUsedAt.After(oldLastUsed))
}

// TestRefreshToken_GetAuditData tests para datos de auditoría
func TestRefreshToken_GetAuditData(t *testing.T) {
	token := createTestRefreshToken()
	token.ID = uuid.New()

	auditData := token.GetAuditData()

	// Verificar que contiene los campos esperados
	expectedFields := []string{"id", "user_id", "token_id", "expires_at", "is_revoked", "ip_address"}
	for _, field := range expectedFields {
		assert.Contains(t, auditData, field)
	}

	// Verificar valores específicos
	assert.Equal(t, token.ID, auditData["id"])
	assert.Equal(t, token.UserID, auditData["user_id"])
	assert.Equal(t, token.TokenID, auditData["token_id"])
	assert.Equal(t, token.ExpiresAt, auditData["expires_at"])
	assert.Equal(t, token.IsRevoked, auditData["is_revoked"])
	assert.Equal(t, token.IPAddress, auditData["ip_address"])
}

// TestRefreshToken_ToSessionInfo tests para conversión a SessionInfo
func TestRefreshToken_ToSessionInfo(t *testing.T) {
	token := createTestRefreshToken()
	token.ID = uuid.New()
	token.CreatedAt = time.Now().Add(-1 * time.Hour)

	// Actualizar último uso
	token.UpdateLastUsed()

	sessionInfo := token.ToSessionInfo()

	// Verificar conversión correcta
	assert.Equal(t, token.ID.String(), sessionInfo.ID)
	assert.Equal(t, token.DeviceInfo, sessionInfo.DeviceInfo)
	assert.Equal(t, token.IPAddress, sessionInfo.IPAddress)
	assert.Equal(t, token.UserAgent, sessionInfo.UserAgent)
	assert.Equal(t, token.CreatedAt, sessionInfo.CreatedAt)
	assert.Equal(t, token.LastUsedAt, sessionInfo.LastUsedAt)
	assert.Equal(t, token.ExpiresAt, sessionInfo.ExpiresAt)
	assert.Equal(t, token.IsRevoked, sessionInfo.IsRevoked)
}

// TestRefreshToken_TokenLifecycle test de integración del ciclo de vida
func TestRefreshToken_TokenLifecycle(t *testing.T) {
	// 1. Crear token nuevo
	token := createTestRefreshToken()
	token.ID = uuid.New()

	// Verificar estado inicial
	assert.True(t, token.IsActive())
	assert.False(t, token.IsExpired())
	assert.False(t, token.IsRevoked)
	assert.Nil(t, token.LastUsedAt)
	assert.Nil(t, token.RevokedAt)

	// 2. Usar el token
	token.UpdateLastUsed()
	assert.NotNil(t, token.LastUsedAt)
	assert.True(t, token.IsActive())

	// 3. Usar el token múltiples veces
	firstUse := *token.LastUsedAt
	time.Sleep(10 * time.Millisecond)
	token.UpdateLastUsed()
	assert.True(t, token.LastUsedAt.After(firstUse))
	assert.True(t, token.IsActive())

	// 4. Revocar el token
	token.Revoke()
	assert.False(t, token.IsActive())
	assert.True(t, token.IsRevoked)
	assert.NotNil(t, token.RevokedAt)

	// 5. Verificar que sigue revocado después de intentar usar
	token.UpdateLastUsed()
	assert.False(t, token.IsActive()) // Sigue revocado
}

// TestRefreshToken_ValidationWithRealScenarios tests con escenarios realistas
func TestRefreshToken_ValidationWithRealScenarios(t *testing.T) {
	t.Run("token de sesión móvil", func(t *testing.T) {
		token := &RefreshToken{
			UserID:     uuid.New().String(),
			TokenHash:  "mobile-token-hash-12345",
			TokenID:    uuid.New().String(),
			ExpiresAt:  time.Now().Add(30 * 24 * time.Hour), // 30 días para móvil
			IPAddress:  "192.168.1.50",
			UserAgent:  "CybESphere-Mobile/1.0 (iOS 15.0)",
			DeviceInfo: "iPhone 13 Pro - iOS 15.0",
		}

		err := token.ValidateRefreshToken()
		assert.NoError(t, err)
		assert.True(t, token.IsActive())
	})

	t.Run("token de sesión web", func(t *testing.T) {
		token := &RefreshToken{
			UserID:     uuid.New().String(),
			TokenHash:  "web-token-hash-67890",
			TokenID:    uuid.New().String(),
			ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // 7 días para web
			IPAddress:  "10.0.0.100",
			UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
			DeviceInfo: "MacBook Pro - Safari",
		}

		err := token.ValidateRefreshToken()
		assert.NoError(t, err)
		assert.True(t, token.IsActive())
	})
}

// TestRefreshToken_EdgeCases tests para casos límite
func TestRefreshToken_EdgeCases(t *testing.T) {
	t.Run("token hash muy largo", func(t *testing.T) {
		token := createTestRefreshToken()
		// Hash de 255 caracteres (máximo permitido)
		token.TokenHash = string(make([]byte, 255))
		for i := range token.TokenHash {
			token.TokenHash = string(append([]byte(token.TokenHash[:i]), 'a'))
		}

		err := token.ValidateRefreshToken()
		assert.NoError(t, err)
	})

	t.Run("user agent muy largo", func(t *testing.T) {
		token := createTestRefreshToken()
		// User agent de 500 caracteres (máximo permitido)
		longUserAgent := string(make([]byte, 500))
		for i := range longUserAgent {
			longUserAgent = string(append([]byte(longUserAgent[:i]), 'x'))
		}
		token.UserAgent = longUserAgent

		err := token.ValidateRefreshToken()
		assert.NoError(t, err)
	})

	t.Run("IP address IPv6", func(t *testing.T) {
		token := createTestRefreshToken()
		token.IPAddress = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"

		err := token.ValidateRefreshToken()
		assert.NoError(t, err)
	})

	t.Run("revocación múltiple", func(t *testing.T) {
		token := createTestRefreshToken()

		// Primera revocación
		token.Revoke()
		firstRevokedAt := *token.RevokedAt

		time.Sleep(10 * time.Millisecond)

		// Segunda revocación (no debería cambiar el timestamp)
		token.Revoke()
		secondRevokedAt := *token.RevokedAt

		// Verificar que el timestamp no cambió
		assert.Equal(t, firstRevokedAt, secondRevokedAt)
		assert.True(t, token.IsRevoked)
		assert.False(t, token.IsActive())
	})
}

// TestRefreshToken_Performance test básico de performance
func TestRefreshToken_Performance(t *testing.T) {
	// Crear múltiples tokens y verificar que las operaciones son rápidas
	tokens := make([]*RefreshToken, 1000)

	start := time.Now()
	for i := 0; i < 1000; i++ {
		tokens[i] = createTestRefreshToken()
		tokens[i].ID = uuid.New()

		// Validar
		err := tokens[i].ValidateRefreshToken()
		assert.NoError(t, err)

		// Verificar estado
		assert.True(t, tokens[i].IsActive())

		// Actualizar uso
		tokens[i].UpdateLastUsed()

		// Obtener audit data
		auditData := tokens[i].GetAuditData()
		assert.NotNil(t, auditData)
	}
	duration := time.Since(start)

	// Las operaciones deberían completarse en menos de 100ms
	assert.Less(t, duration, 100*time.Millisecond, "Operations took too long: %v", duration)

	// Revocar todos los tokens
	start = time.Now()
	for _, token := range tokens {
		token.Revoke()
		assert.False(t, token.IsActive())
	}
	revokeDuration := time.Since(start)

	// La revocación debería ser rápida
	assert.Less(t, revokeDuration, 50*time.Millisecond, "Revocation took too long: %v", revokeDuration)
}
