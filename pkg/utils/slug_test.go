package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateSlug tests para generación de slugs
func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{
			name:      "texto normal",
			input:     "Workshop de Ciberseguridad",
			maxLength: 0, // Sin límite
			expected:  "workshop-de-ciberseguridad",
		},
		{
			name:      "con acentos españoles",
			input:     "Introducción a la Tecnología",
			maxLength: 0,
			expected:  "introduccion-a-la-tecnologia",
		},
		{
			name:      "con caracteres especiales",
			input:     "¡Workshop Práctico! 100% Gratuito & Útil",
			maxLength: 0,
			expected:  "workshop-practico-100-gratuito-util",
		},
		{
			name:      "con ñ española",
			input:     "Niños y Niñas: Año 2024",
			maxLength: 0,
			expected:  "ninos-y-ninas-ano-2024",
		},
		{
			name:      "texto muy largo con límite",
			input:     "Este es un texto extremadamente largo que debería ser truncado apropiadamente",
			maxLength: 50,
			expected:  "este-es-un-texto-extremadamente-largo-que-deberia",
		},
		{
			name:      "espacios múltiples y tabs",
			input:     "Tech   Company\t\tMadrid   2024",
			maxLength: 0,
			expected:  "tech-company-madrid-2024",
		},
		{
			name:      "guiones y underscores",
			input:     "Tech_Company-Madrid--2024",
			maxLength: 0,
			expected:  "tech-company-madrid-2024",
		},
		{
			name:      "solo caracteres especiales",
			input:     "!@#$%^&*()",
			maxLength: 0,
			expected:  "",
		},
		{
			name:      "texto vacío",
			input:     "",
			maxLength: 0,
			expected:  "",
		},
		{
			name:      "solo espacios",
			input:     "   ",
			maxLength: 0,
			expected:  "",
		},
		{
			name:      "números y letras",
			input:     "Workshop 2024: Python & Go",
			maxLength: 0,
			expected:  "workshop-2024-python-go",
		},
		{
			name:      "truncamiento exacto en límite",
			input:     "12345678901234567890",
			maxLength: 10,
			expected:  "1234567890",
		},
		{
			name:      "truncamiento con guión al final",
			input:     "Workshop de Ciberseguridad",
			maxLength: 15,
			expected:  "workshop-de-cib", // Función trunca a 15 caracteres exactos
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input, tt.maxLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNormalizeText tests para normalización de texto
func TestNormalizeText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "texto normal",
			input:    "Hello World",
			expected: "hello world",
		},
		{
			name:     "con espacios extra",
			input:    "  Hello World  ",
			expected: "hello world",
		},
		{
			name:     "mayúsculas mixtas",
			input:    "HeLLo WoRLd",
			expected: "hello world",
		},
		{
			name:     "texto vacío",
			input:    "",
			expected: "",
		},
		{
			name:     "solo espacios",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeText(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNormalizeEmail tests para normalización de email
func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "email normal",
			input:    "user@example.com",
			expected: "user@example.com",
		},
		{
			name:     "con mayúsculas",
			input:    "USER@EXAMPLE.COM",
			expected: "user@example.com",
		},
		{
			name:     "con espacios",
			input:    "  user@example.com  ",
			expected: "user@example.com",
		},
		{
			name:     "mixto",
			input:    "  UsEr@ExAmPlE.CoM  ",
			expected: "user@example.com",
		},
		{
			name:     "email vacío",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeEmail(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCharacterReplacements tests para el mapa de reemplazos
func TestCharacterReplacements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "acentos en a",
			input:    "áàäâãå",
			expected: "aaaaaa",
		},
		{
			name:     "acentos en e",
			input:    "éèëê",
			expected: "eeee",
		},
		{
			name:     "acentos en i",
			input:    "íìïî",
			expected: "iiii",
		},
		{
			name:     "acentos en o",
			input:    "óòöôõ",
			expected: "ooooo",
		},
		{
			name:     "acentos en u",
			input:    "úùüû",
			expected: "uuuu",
		},
		{
			name:     "ñ española",
			input:    "niño",
			expected: "nino",
		},
		{
			name:     "mayúsculas",
			input:    "ÑOÑO",
			expected: "nono",
		},
		{
			name:     "signos españoles",
			input:    "¿Cómo está?¡Bien!",
			expected: "como-estabien", // Los signos se eliminan sin agregar espacios
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.input, 0)
			assert.Equal(t, tt.expected, result)
		})
	}
}
