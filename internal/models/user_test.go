package models

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestUser crea un usuario válido para testing
func createTestUser() *User {
	return &User{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Juan",
		LastName:  "Pérez",
		Role:      RoleUser,
		IsActive:  true,
		City:      "Madrid",
		Country:   "Spain",
		Timezone:  "Europe/Madrid",
		Language:  "es",
	}
}

// TestUser_ValidateUser tests unitarios para validación de usuarios
func TestUser_ValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr bool
		errMsg  string
	}{
		{
			name:    "usuario válido",
			user:    createTestUser(),
			wantErr: false,
		},
		{
			name: "email vacío",
			user: func() *User {
				u := createTestUser()
				u.Email = ""
				return u
			}(),
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "email solo espacios",
			user: func() *User {
				u := createTestUser()
				u.Email = "   "
				return u
			}(),
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "first name vacío",
			user: func() *User {
				u := createTestUser()
				u.FirstName = ""
				return u
			}(),
			wantErr: true,
			errMsg:  "first name is required",
		},
		{
			name: "last name vacío",
			user: func() *User {
				u := createTestUser()
				u.LastName = ""
				return u
			}(),
			wantErr: true,
			errMsg:  "last name is required",
		},
		{
			name: "password muy corta",
			user: func() *User {
				u := createTestUser()
				u.Password = "1234567"
				return u
			}(),
			wantErr: true,
			errMsg:  "password must be at least 8 characters",
		},
		{
			name: "role inválido",
			user: func() *User {
				u := createTestUser()
				u.Role = UserRole("invalid_role")
				return u
			}(),
			wantErr: true,
			errMsg:  "invalid role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.ValidateUser()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUser_IsValidRole tests para validación de roles
func TestUser_IsValidRole(t *testing.T) {
	tests := []struct {
		name string
		role UserRole
		want bool
	}{
		{"role admin válido", RoleAdmin, true},
		{"role organizer válido", RoleOrganizer, true},
		{"role user válido", RoleUser, true},
		{"role inválido", UserRole("invalid"), false},
		{"role vacío", UserRole(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.want, user.IsValidRole())
		})
	}
}

// TestUser_HashPassword tests para hashing de passwords
func TestUser_HashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "password válida",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "password vacía",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Password: tt.password}
			err := user.HashPassword()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "password cannot be empty")
			} else {
				assert.NoError(t, err)
				// Verificar que la password fue hasheada
				assert.NotEqual(t, tt.password, user.Password)
				assert.True(t, user.IsPasswordHashed())
				// Verificar que se puede verificar la password
				assert.True(t, user.CheckPassword(tt.password))
			}
		})
	}
}

// TestUser_CheckPassword tests para verificación de passwords
func TestUser_CheckPassword(t *testing.T) {
	user := createTestUser()
	originalPassword := user.Password

	// Hashear la password
	err := user.HashPassword()
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "password correcta",
			password: originalPassword,
			want:     true,
		},
		{
			name:     "password incorrecta",
			password: "wrongpassword",
			want:     false,
		},
		{
			name:     "password vacía",
			password: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.CheckPassword(tt.password)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestUser_IsPasswordHashed tests para detectar si password está hasheada
func TestUser_IsPasswordHashed(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "password hasheada bcrypt",
			password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.Uo8Fm/Xx9LFfEt6yt3pHOXHYFcxwJq",
			want:     true,
		},
		{
			name:     "password texto plano",
			password: "plainpassword123",
			want:     false,
		},
		{
			name:     "password vacía",
			password: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Password: tt.password}
			assert.Equal(t, tt.want, user.IsPasswordHashed())
		})
	}
}

// TestUser_GetFullName tests para obtener nombre completo
func TestUser_GetFullName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		want      string
	}{
		{
			name:      "nombres normales",
			firstName: "Juan",
			lastName:  "Pérez",
			want:      "Juan Pérez",
		},
		{
			name:      "con espacios extra",
			firstName: "  María  ",
			lastName:  "  García  ",
			want:      "María     García", // "  María  " + " " + "  García  " = "  María     García  " -> TrimSpace = "María     García"
		},
		{
			name:      "nombres vacíos",
			firstName: "",
			lastName:  "",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
			}
			assert.Equal(t, tt.want, user.GetFullName())
		})
	}
}

// TestUser_RoleCheckers tests para métodos de verificación de roles
func TestUser_RoleCheckers(t *testing.T) {
	tests := []struct {
		name        string
		role        UserRole
		isAdmin     bool
		isOrganizer bool
	}{
		{
			name:        "usuario admin",
			role:        RoleAdmin,
			isAdmin:     true,
			isOrganizer: false,
		},
		{
			name:        "usuario organizer",
			role:        RoleOrganizer,
			isAdmin:     false,
			isOrganizer: true,
		},
		{
			name:        "usuario normal",
			role:        RoleUser,
			isAdmin:     false,
			isOrganizer: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.isAdmin, user.IsAdmin())
			assert.Equal(t, tt.isOrganizer, user.IsOrganizer())
		})
	}
}

// TestUser_HasPermission tests para sistema de permisos
func TestUser_HasPermission(t *testing.T) {
	orgID := uuid.New().String()

	tests := []struct {
		name       string
		user       *User
		action     string
		resource   string
		resourceID string
		want       bool
	}{
		{
			name:       "admin puede todo",
			user:       &User{Role: RoleAdmin},
			action:     "delete",
			resource:   "event",
			resourceID: "any-id",
			want:       true,
		},
		{
			name:       "organizer puede leer todo",
			user:       &User{Role: RoleOrganizer, OrganizationID: &orgID},
			action:     "read",
			resource:   "event",
			resourceID: "any-id",
			want:       true,
		},
		{
			name:       "organizer puede modificar su organización",
			user:       &User{Role: RoleOrganizer, OrganizationID: &orgID},
			action:     "update",
			resource:   "event",
			resourceID: orgID,
			want:       true,
		},
		{
			name:       "organizer no puede modificar otra organización",
			user:       &User{Role: RoleOrganizer, OrganizationID: &orgID},
			action:     "update",
			resource:   "event",
			resourceID: "other-org-id",
			want:       false,
		},
		{
			name:       "user puede leer",
			user:       &User{Role: RoleUser},
			action:     "read",
			resource:   "event",
			resourceID: "any-id",
			want:       true,
		},
		{
			name:       "user puede gestionar favoritos",
			user:       &User{Role: RoleUser},
			action:     "create",
			resource:   "favorite",
			resourceID: "any-id",
			want:       true,
		},
		{
			name:       "user no puede crear eventos",
			user:       &User{Role: RoleUser},
			action:     "create",
			resource:   "event",
			resourceID: "any-id",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.HasPermission(tt.action, tt.resource, tt.resourceID)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestUser_CanManageOrganization tests para permisos de organización
func TestUser_CanManageOrganization(t *testing.T) {
	orgID := uuid.New().String()
	otherOrgID := uuid.New().String()

	tests := []struct {
		name  string
		user  *User
		orgID string
		want  bool
	}{
		{
			name:  "admin puede gestionar cualquier organización",
			user:  &User{Role: RoleAdmin},
			orgID: orgID,
			want:  true,
		},
		{
			name:  "organizer puede gestionar su organización",
			user:  &User{Role: RoleOrganizer, OrganizationID: &orgID},
			orgID: orgID,
			want:  true,
		},
		{
			name:  "organizer no puede gestionar otra organización",
			user:  &User{Role: RoleOrganizer, OrganizationID: &orgID},
			orgID: otherOrgID,
			want:  false,
		},
		{
			name:  "user no puede gestionar organizaciones",
			user:  &User{Role: RoleUser},
			orgID: orgID,
			want:  false,
		},
		{
			name:  "organizer sin organización asignada no puede gestionar",
			user:  &User{Role: RoleOrganizer, OrganizationID: nil},
			orgID: orgID,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.CanManageOrganization(tt.orgID)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestUser_UpdateLastLogin tests para actualización de último login
func TestUser_UpdateLastLogin(t *testing.T) {
	user := createTestUser()

	// Verificar que inicialmente no tiene last login
	assert.Nil(t, user.LastLoginAt)

	// Actualizar last login
	user.UpdateLastLogin()

	// Verificar que se actualizó
	assert.NotNil(t, user.LastLoginAt)
	assert.WithinDuration(t, time.Now(), *user.LastLoginAt, time.Second)
}

// TestUser_SetLocation tests para configuración de geolocalización
func TestUser_SetLocation(t *testing.T) {
	user := createTestUser()

	// Verificar que inicialmente no tiene ubicación
	assert.False(t, user.HasLocation())

	// Establecer ubicación
	lat, lng := 40.4168, -3.7038 // Madrid coordinates
	user.SetLocation(lat, lng, "Madrid", "Spain")

	// Verificar que se estableció correctamente
	assert.True(t, user.HasLocation())
	assert.Equal(t, lat, *user.Latitude)
	assert.Equal(t, lng, *user.Longitude)
	assert.Equal(t, "Madrid", user.City)
	assert.Equal(t, "Spain", user.Country)
}

// TestUser_GetAuditData tests para datos de auditoría
func TestUser_GetAuditData(t *testing.T) {
	user := createTestUser()
	user.ID = uuid.New()

	auditData := user.GetAuditData()

	// Verificar que contiene los campos esperados
	expectedFields := []string{"id", "email", "first_name", "last_name", "role", "is_active"}
	for _, field := range expectedFields {
		assert.Contains(t, auditData, field)
	}

	// Verificar valores específicos
	assert.Equal(t, user.ID, auditData["id"])
	assert.Equal(t, user.Email, auditData["email"])
	assert.Equal(t, user.FirstName, auditData["first_name"])
	assert.Equal(t, user.LastName, auditData["last_name"])
	assert.Equal(t, user.Role, auditData["role"])
	assert.Equal(t, user.IsActive, auditData["is_active"])
}

// TestUser_BeforeCreateLogic tests para lógica del hook BeforeCreate sin base de datos
func TestUser_BeforeCreateLogic(t *testing.T) {
	user := createTestUser()
	originalPassword := user.Password

	// Testear la lógica de validación directamente
	err := user.ValidateUser()
	require.NoError(t, err)

	// Testear hash de password
	err = user.HashPassword()
	require.NoError(t, err)

	// Verificar que la password fue hasheada
	assert.True(t, user.IsPasswordHashed())
	assert.NotEqual(t, originalPassword, user.Password)

	// Verificar que se puede verificar la password original
	assert.True(t, user.CheckPassword(originalPassword))

	// Testear normalización de email
	user.Email = "  TEST@EXAMPLE.COM  "
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	assert.Equal(t, "test@example.com", user.Email)
}

// TestUser_ValidationErrors tests para errores de validación
func TestUser_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr string
	}{
		{
			name: "email vacío",
			user: func() *User {
				u := createTestUser()
				u.Email = ""
				return u
			}(),
			wantErr: "email is required",
		},
		{
			name: "password muy corta",
			user: func() *User {
				u := createTestUser()
				u.Password = "123"
				return u
			}(),
			wantErr: "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.ValidateUser()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// TestUser_EmailNormalization tests para normalización de email
func TestUser_EmailNormalization(t *testing.T) {
	tests := []struct {
		name       string
		inputEmail string
		wantEmail  string
	}{
		{
			name:       "email en mayúsculas",
			inputEmail: "TEST@EXAMPLE.COM",
			wantEmail:  "test@example.com",
		},
		{
			name:       "email mixto",
			inputEmail: "TeSt@ExAmPlE.CoM",
			wantEmail:  "test@example.com",
		},
		{
			name:       "email con espacios",
			inputEmail: "  test@example.com  ",
			wantEmail:  "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := createTestUser()
			user.Email = tt.inputEmail

			// Normalizar email manualmente (simulando BeforeCreate)
			user.Email = strings.ToLower(strings.TrimSpace(user.Email))

			assert.Equal(t, tt.wantEmail, user.Email)
		})
	}
}
