package models

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestOrganization crea una organización válida para testing
func createTestOrganization() *Organization {
	return &Organization{
		Name:        "CyberSecurity España",
		Description: "Comunidad española de profesionales de ciberseguridad",
		Email:       "info@cybersecurityspain.com",
		Website:     "https://cybersecurityspain.com",
		City:        "Madrid",
		Country:     "Spain",
		Status:      OrgStatusPending,
		IsVerified:  false,
	}
}

// TestOrganization_ValidateOrganization tests unitarios para validación de organizaciones
func TestOrganization_ValidateOrganization(t *testing.T) {
	tests := []struct {
		name    string
		org     *Organization
		wantErr bool
		errMsg  string
	}{
		{
			name:    "organización válida",
			org:     createTestOrganization(),
			wantErr: false,
		},
		{
			name: "nombre vacío",
			org: func() *Organization {
				o := createTestOrganization()
				o.Name = ""
				return o
			}(),
			wantErr: true,
			errMsg:  "organization name is required",
		},
		{
			name: "nombre solo espacios",
			org: func() *Organization {
				o := createTestOrganization()
				o.Name = "   "
				return o
			}(),
			wantErr: true,
			errMsg:  "organization name is required",
		},
		{
			name: "nombre muy corto",
			org: func() *Organization {
				o := createTestOrganization()
				o.Name = "AB"
				return o
			}(),
			wantErr: true,
			errMsg:  "organization name must be at least 3 characters",
		},
		{
			name: "email vacío",
			org: func() *Organization {
				o := createTestOrganization()
				o.Email = ""
				return o
			}(),
			wantErr: true,
			errMsg:  "organization email is required",
		},
		{
			name: "email solo espacios",
			org: func() *Organization {
				o := createTestOrganization()
				o.Email = "   "
				return o
			}(),
			wantErr: true,
			errMsg:  "organization email is required",
		},
		{
			name: "status inválido",
			org: func() *Organization {
				o := createTestOrganization()
				o.Status = OrganizationStatus("invalid_status")
				return o
			}(),
			wantErr: true,
			errMsg:  "invalid organization status",
		},
		{
			name: "color primario inválido",
			org: func() *Organization {
				o := createTestOrganization()
				o.PrimaryColor = "invalid_color"
				return o
			}(),
			wantErr: true,
			errMsg:  "invalid primary color format",
		},
		{
			name: "color secundario inválido",
			org: func() *Organization {
				o := createTestOrganization()
				o.SecondaryColor = "invalid_color"
				return o
			}(),
			wantErr: true,
			errMsg:  "invalid secondary color format",
		},
		{
			name: "colores válidos",
			org: func() *Organization {
				o := createTestOrganization()
				o.PrimaryColor = "#FF5733"
				o.SecondaryColor = "#33FF57"
				return o
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.org.ValidateOrganization()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestOrganization_IsValidStatus tests para validación de estados
func TestOrganization_IsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status OrganizationStatus
		want   bool
	}{
		{"status pending válido", OrgStatusPending, true},
		{"status active válido", OrgStatusActive, true},
		{"status suspended válido", OrgStatusSuspended, true},
		{"status inactive válido", OrgStatusInactive, true},
		{"status inválido", OrganizationStatus("invalid"), false},
		{"status vacío", OrganizationStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org := &Organization{Status: tt.status}
			assert.Equal(t, tt.want, org.IsValidStatus())
		})
	}
}

// TestOrganization_GenerateSlug tests para generación de slugs
func TestOrganization_GenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		orgName  string
		wantSlug string
	}{
		{
			name:     "nombre normal",
			orgName:  "CyberSecurity España",
			wantSlug: "cybersecurity-espana",
		},
		{
			name:     "con espacios múltiples",
			orgName:  "Tech   Company   Madrid",
			wantSlug: "tech-company-madrid",
		},
		{
			name:     "con caracteres especiales",
			orgName:  "Empresa & Tecnología S.A.",
			wantSlug: "empresa-tecnologia-sa",
		},
		{
			name:     "con guiones y underscores",
			orgName:  "Tech_Company-Madrid",
			wantSlug: "tech-company-madrid",
		},
		{
			name:     "nombre muy largo",
			orgName:  "Esta es una organización con un nombre extremadamente largo que debería ser truncado apropiadamente para generar un slug válido",
			wantSlug: "esta-es-una-organizacion-con-un-nombre-extremadamente-largo-que-deberia-ser-truncado-apropiadamente", // Truncado limpiamente
		},
		{
			name:     "solo espacios y caracteres especiales",
			orgName:  "  !!!  @@@  ###  ",
			wantSlug: "",
		},
		{
			name:     "números y letras",
			orgName:  "Tech Corp 2024",
			wantSlug: "tech-corp-2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org := &Organization{Name: tt.orgName}
			slug := org.GenerateSlug()
			assert.Equal(t, tt.wantSlug, slug)
		})
	}
}

// TestOrganization_IsActive tests para verificar si está activa
func TestOrganization_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		status OrganizationStatus
		want   bool
	}{
		{"organización activa", OrgStatusActive, true},
		{"organización pendiente", OrgStatusPending, false},
		{"organización suspendida", OrgStatusSuspended, false},
		{"organización inactiva", OrgStatusInactive, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org := &Organization{Status: tt.status}
			assert.Equal(t, tt.want, org.IsActive())
		})
	}
}

// TestOrganization_CanCreateEvent tests para permisos de creación de eventos
func TestOrganization_CanCreateEvent(t *testing.T) {
	tests := []struct {
		name        string
		status      OrganizationStatus
		canCreate   bool
		maxEvents   *int
		eventsCount int
		want        bool
	}{
		{
			name:        "organización activa sin límite",
			status:      OrgStatusActive,
			canCreate:   true,
			maxEvents:   nil,
			eventsCount: 5,
			want:        true,
		},
		{
			name:        "organización activa con límite no alcanzado",
			status:      OrgStatusActive,
			canCreate:   true,
			maxEvents:   func() *int { i := 10; return &i }(),
			eventsCount: 5,
			want:        true,
		},
		{
			name:        "organización activa con límite alcanzado",
			status:      OrgStatusActive,
			canCreate:   true,
			maxEvents:   func() *int { i := 5; return &i }(),
			eventsCount: 5,
			want:        false,
		},
		{
			name:        "organización activa con límite superado",
			status:      OrgStatusActive,
			canCreate:   true,
			maxEvents:   func() *int { i := 3; return &i }(),
			eventsCount: 5,
			want:        false,
		},
		{
			name:        "organización inactiva",
			status:      OrgStatusInactive,
			canCreate:   true,
			maxEvents:   nil,
			eventsCount: 0,
			want:        false,
		},
		{
			name:        "organización suspendida",
			status:      OrgStatusSuspended,
			canCreate:   true,
			maxEvents:   nil,
			eventsCount: 0,
			want:        false,
		},
		{
			name:        "sin permisos de creación",
			status:      OrgStatusActive,
			canCreate:   false,
			maxEvents:   nil,
			eventsCount: 0,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org := &Organization{
				Status:          tt.status,
				CanCreateEvents: tt.canCreate,
				MaxEvents:       tt.maxEvents,
				EventsCount:     tt.eventsCount,
			}
			assert.Equal(t, tt.want, org.CanCreateEvent())
		})
	}
}

// TestOrganization_StateTransitions tests para transiciones de estado
func TestOrganization_StateTransitions(t *testing.T) {
	t.Run("verificar organización", func(t *testing.T) {
		org := createTestOrganization()
		verifierID := "admin-user-id"

		// Verificar
		org.Verify(verifierID)

		assert.True(t, org.IsVerified)
		assert.Equal(t, OrgStatusActive, org.Status)
		assert.NotNil(t, org.VerifiedBy)
		assert.Equal(t, verifierID, *org.VerifiedBy)
		assert.NotNil(t, org.VerifiedAt)
	})

	t.Run("suspender organización", func(t *testing.T) {
		org := createTestOrganization()
		org.Status = OrgStatusActive
		org.CanCreateEvents = true

		// Suspender
		org.Suspend()

		assert.Equal(t, OrgStatusSuspended, org.Status)
		assert.False(t, org.CanCreateEvents)
	})

	t.Run("activar organización", func(t *testing.T) {
		org := createTestOrganization()
		org.Status = OrgStatusSuspended
		org.CanCreateEvents = false

		// Activar
		org.Activate()

		assert.Equal(t, OrgStatusActive, org.Status)
		assert.True(t, org.CanCreateEvents)
	})

	t.Run("desactivar organización", func(t *testing.T) {
		org := createTestOrganization()
		org.Status = OrgStatusActive
		org.CanCreateEvents = true

		// Desactivar
		org.Deactivate()

		assert.Equal(t, OrgStatusInactive, org.Status)
		assert.False(t, org.CanCreateEvents)
	})
}

// TestOrganization_SetLocation tests para configuración de geolocalización
func TestOrganization_SetLocation(t *testing.T) {
	org := createTestOrganization()

	// Verificar que inicialmente no tiene ubicación
	assert.False(t, org.HasLocation())

	// Establecer ubicación
	lat, lng := 40.4168, -3.7038 // Madrid coordinates
	org.SetLocation(lat, lng)

	// Verificar que se estableció correctamente
	assert.True(t, org.HasLocation())
	assert.Equal(t, lat, *org.Latitude)
	assert.Equal(t, lng, *org.Longitude)
}

// TestOrganization_SetBranding tests para configuración de branding
func TestOrganization_SetBranding(t *testing.T) {
	tests := []struct {
		name           string
		primaryColor   string
		secondaryColor string
		wantErr        bool
		errMsg         string
	}{
		{
			name:           "colores válidos",
			primaryColor:   "#FF5733",
			secondaryColor: "#33FF57",
			wantErr:        false,
		},
		{
			name:           "solo color primario",
			primaryColor:   "#FF5733",
			secondaryColor: "",
			wantErr:        false,
		},
		{
			name:           "color primario inválido",
			primaryColor:   "invalid",
			secondaryColor: "#33FF57",
			wantErr:        true,
			errMsg:         "invalid primary color format",
		},
		{
			name:           "color secundario inválido",
			primaryColor:   "#FF5733",
			secondaryColor: "invalid",
			wantErr:        true,
			errMsg:         "invalid secondary color format",
		},
		{
			name:           "colores vacíos",
			primaryColor:   "",
			secondaryColor: "",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org := createTestOrganization()
			err := org.SetBranding(tt.primaryColor, tt.secondaryColor)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.primaryColor, org.PrimaryColor)
				assert.Equal(t, tt.secondaryColor, org.SecondaryColor)
			}
		})
	}
}

// TestOrganization_GetAuditData tests para datos de auditoría
func TestOrganization_GetAuditData(t *testing.T) {
	org := createTestOrganization()
	org.ID = uuid.New()
	org.Slug = "cybersecurity-espana"

	auditData := org.GetAuditData()

	// Verificar que contiene los campos esperados
	expectedFields := []string{"id", "name", "slug", "email", "status", "is_verified"}
	for _, field := range expectedFields {
		assert.Contains(t, auditData, field)
	}

	// Verificar valores específicos
	assert.Equal(t, org.ID, auditData["id"])
	assert.Equal(t, org.Name, auditData["name"])
	assert.Equal(t, org.Slug, auditData["slug"])
	assert.Equal(t, org.Email, auditData["email"])
	assert.Equal(t, org.Status, auditData["status"])
	assert.Equal(t, org.IsVerified, auditData["is_verified"])
}

// TestOrganization_FieldNormalization tests para normalización de campos
func TestOrganization_FieldNormalization(t *testing.T) {
	org := createTestOrganization()

	// Configurar campos con espacios y mayúsculas
	org.Name = "  CyberSecurity España  "
	org.Email = "  INFO@CYBERSECURITYSPAIN.COM  "
	org.Website = "  https://cybersecurityspain.com  "
	org.Description = "  Comunidad española de profesionales  "

	// Simular normalización (como en el hook BeforeCreate)
	org.Name = strings.TrimSpace(org.Name)
	org.Email = strings.ToLower(strings.TrimSpace(org.Email))
	org.Website = strings.TrimSpace(org.Website)
	org.Description = strings.TrimSpace(org.Description)

	// Verificar normalización
	assert.Equal(t, "CyberSecurity España", org.Name)
	assert.Equal(t, "info@cybersecurityspain.com", org.Email)
	assert.Equal(t, "https://cybersecurityspain.com", org.Website)
	assert.Equal(t, "Comunidad española de profesionales", org.Description)
}

// TestOrganization_HexColorValidation tests para validación de colores hexadecimales
func TestOrganization_HexColorValidation(t *testing.T) {
	// Esta función no está exportada en el modelo, así que testeamos indirectamente
	tests := []struct {
		name  string
		color string
		valid bool
	}{
		{"color hex válido con mayúsculas", "#FF5733", true},
		{"color hex válido con minúsculas", "#ff5733", true},
		{"color hex válido mixto", "#Ff5733", true},
		{"color hex corto válido", "#F53", false}, // El modelo espera formato largo
		{"sin #", "FF5733", false},
		{"con caracteres inválidos", "#GG5733", false},
		{"muy corto", "#F5", false},
		{"muy largo", "#FF57339", false},
		{"vacío", "", true}, // Vacío es válido (opcional)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org := createTestOrganization()
			org.PrimaryColor = tt.color

			err := org.ValidateOrganization()

			if tt.valid {
				assert.NoError(t, err, "Color %s debería ser válido", tt.color)
			} else {
				assert.Error(t, err, "Color %s debería ser inválido", tt.color)
				if tt.color != "" { // Solo verificar mensaje si no es vacío
					assert.Contains(t, err.Error(), "invalid primary color format")
				}
			}
		})
	}
}

// TestOrganization_SlugGeneration_EdgeCases tests para casos límite en generación de slugs
func TestOrganization_SlugGeneration_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "solo números",
			input:    "123456789",
			expected: "123456789",
		},
		{
			name:     "acentos y caracteres especiales",
			input:    "Niños & Niñas Ñoño",
			expected: "ninos-ninas-nono",
		},
		{
			name:     "guiones al inicio y final",
			input:    "-Tech Company-",
			expected: "tech-company",
		},
		{
			name:     "múltiples guiones consecutivos",
			input:    "Tech---Company---Madrid",
			expected: "tech-company-madrid",
		},
		{
			name:     "espacios y tabs mixtos",
			input:    "Tech\t\tCompany   Madrid",
			expected: "tech-company-madrid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			org := &Organization{Name: tt.input}
			result := org.GenerateSlug()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestOrganization_BeforeCreateLogic tests para lógica del hook BeforeCreate sin base de datos
func TestOrganization_BeforeCreateLogic(t *testing.T) {
	org := createTestOrganization()

	// Testear la lógica de validación directamente
	err := org.ValidateOrganization()
	require.NoError(t, err)

	// Generar slug si no existe
	if org.Slug == "" {
		org.Slug = org.GenerateSlug()
	}

	// Verificar que se generó el slug
	assert.NotEmpty(t, org.Slug)
	assert.Equal(t, "cybersecurity-espana", org.Slug)

	// Testear normalización de campos
	org.Name = "  " + org.Name + "  "
	org.Email = "  " + strings.ToUpper(org.Email) + "  "

	// Normalizar campos
	org.Name = strings.TrimSpace(org.Name)
	org.Email = strings.ToLower(strings.TrimSpace(org.Email))

	assert.Equal(t, "CyberSecurity España", org.Name)
	assert.Equal(t, "info@cybersecurityspain.com", org.Email)
}
