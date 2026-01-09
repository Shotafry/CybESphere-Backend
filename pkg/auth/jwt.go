package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrInvalidClaims    = errors.New("invalid token claims")
)

// TokenType define los tipos de token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims estructura personalizada para JWT claims
type Claims struct {
	UserID   string    `json:"user_id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	TokenID  string    `json:"token_id"` // Único para cada token
	Type     TokenType `json:"type"`     // access o refresh
	IssuedAt time.Time `json:"issued_at"`
	jwt.RegisteredClaims
}

// JWTManager maneja la generación y validación de tokens JWT
type JWTManager struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

// NewJWTManager crea una nueva instancia del manager JWT
func NewJWTManager(secretKey string, accessDuration, refreshDuration time.Duration, issuer string) (*JWTManager, error) {
	if len(secretKey) < 32 {
		return nil, errors.New("JWT secret key must be at least 32 characters long")
	}

	return &JWTManager{
		secretKey:            secretKey,
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
		issuer:               issuer,
	}, nil
}

// GenerateTokenPair genera un par de tokens (access y refresh) para un usuario
func (m *JWTManager) GenerateTokenPair(userID, email, role string) (*TokenPair, error) {
	if userID == "" || email == "" || role == "" {
		return nil, errors.New("userID, email and role are required")
	}

	now := time.Now()

	// Generar access token
	accessTokenID := generateTokenID()
	accessToken, err := m.generateToken(userID, email, role, accessTokenID, AccessToken, now, m.accessTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generar refresh token
	refreshTokenID := generateTokenID()
	refreshToken, err := m.generateToken(userID, email, role, refreshTokenID, RefreshToken, now, m.refreshTokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  now.Add(m.accessTokenDuration),
		RefreshTokenExpiresAt: now.Add(m.refreshTokenDuration),
		TokenType:             "Bearer",
	}, nil
}

// generateToken genera un token individual
func (m *JWTManager) generateToken(userID, email, role, tokenID string, tokenType TokenType, issuedAt time.Time, duration time.Duration) (string, error) {
	expiresAt := issuedAt.Add(duration)

	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		TokenID:  tokenID,
		Type:     tokenType,
		IssuedAt: issuedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
			Issuer:    m.issuer,
			Subject:   userID,
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// ValidateAccessToken valida un access token y retorna las claims
func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, AccessToken)
}

// ValidateRefreshToken valida un refresh token y retorna las claims
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, RefreshToken)
}

// validateToken valida un token y verifica su tipo
func (m *JWTManager) validateToken(tokenString string, expectedType TokenType) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firmado sea correcto
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Verificar que el tipo de token sea correcto
	if claims.Type != expectedType {
		return nil, ErrInvalidTokenType
	}

	// Verificar claims básicas
	if claims.UserID == "" || claims.Email == "" || claims.Role == "" {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// RefreshTokens genera un nuevo par de tokens usando un refresh token válido
func (m *JWTManager) RefreshTokens(refreshTokenString string) (*TokenPair, error) {
	// Validar el refresh token
	claims, err := m.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generar nuevos tokens
	return m.GenerateTokenPair(claims.UserID, claims.Email, claims.Role)
}

// ExtractTokenFromHeader extrae el token del header Authorization
func (m *JWTManager) ExtractTokenFromHeader(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "

	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", errors.New("token is required")
	}

	return token, nil
}

// GetTokenClaims extrae claims de un token sin validar expiración (útil para refresh)
func (m *JWTManager) GetTokenClaims(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// IsTokenExpired verifica si un token ha expirado sin validar la firma
func (m *JWTManager) IsTokenExpired(tokenString string) bool {
	claims, err := m.GetTokenClaims(tokenString)
	if err != nil {
		return true
	}

	return claims.ExpiresAt.Before(time.Now())
}

// TokenPair representa un par de tokens
type TokenPair struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	TokenType             string    `json:"token_type"`
}

// generateTokenID genera un ID único para un token
func generateTokenID() string {
	// Generar UUID para máxima unicidad
	id := uuid.New()
	return id.String()
}

// GenerateSecureRandomString genera una cadena aleatoria segura
func GenerateSecureRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashRefreshToken genera un hash del refresh token para almacenar en BD
func HashRefreshToken(token string) (string, error) {
	if token == "" {
		return "", errors.New("token cannot be empty")
	}
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:]), nil
}
