package utils

import "strings"

// CharacterReplacements mapa de caracteres especiales y sus equivalentes ASCII
var CharacterReplacements = map[rune]string{
	// Vocales con acentos
	'á': "a", 'à': "a", 'ä': "a", 'â': "a", 'ã': "a", 'å': "a",
	'Á': "a", 'À': "a", 'Ä': "a", 'Â': "a", 'Ã': "a", 'Å': "a",
	'é': "e", 'è': "e", 'ë': "e", 'ê': "e",
	'É': "e", 'È': "e", 'Ë': "e", 'Ê': "e",
	'í': "i", 'ì': "i", 'ï': "i", 'î': "i",
	'Í': "i", 'Ì': "i", 'Ï': "i", 'Î': "i",
	'ó': "o", 'ò': "o", 'ö': "o", 'ô': "o", 'õ': "o",
	'Ó': "o", 'Ò': "o", 'Ö': "o", 'Ô': "o", 'Õ': "o",
	'ú': "u", 'ù': "u", 'ü': "u", 'û': "u",
	'Ú': "u", 'Ù': "u", 'Ü': "u", 'Û': "u",
	// Caracteres específicos del español
	'ñ': "n", 'Ñ': "n",
	'ç': "c", 'Ç': "c",
	// Caracteres alemanes y otros
	'ß': "ss",
	// Signos de puntuación españoles
	'¿': "", '¡': "",
	// Símbolos comunes en títulos
	'%': "", '&': "",
}

// GenerateSlug genera un slug limpio a partir de cualquier string
// Parámetros:
//   - text: el texto a convertir en slug
//   - maxLength: longitud máxima del slug (0 = sin límite)
func GenerateSlug(text string, maxLength int) string {
	if text == "" {
		return ""
	}

	// Convertir a minúsculas
	slug := strings.ToLower(text)

	// Reemplazar caracteres especiales y acentuados
	var result strings.Builder
	for _, char := range slug {
		if replacement, exists := CharacterReplacements[char]; exists {
			result.WriteString(replacement)
		} else {
			result.WriteRune(char)
		}
	}
	slug = result.String()

	// Reemplazar espacios, tabs y underscores por guiones
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "\t", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Mantener solo caracteres alfanuméricos y guiones
	var cleanSlug strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			cleanSlug.WriteRune(char)
		}
	}
	slug = cleanSlug.String()

	// Limpiar guiones duplicados
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Remover guiones al inicio y final
	slug = strings.Trim(slug, "-")

	// Truncar si se especifica longitud máxima
	if maxLength > 0 && len(slug) > maxLength {
		slug = slug[:maxLength]
		// Asegurar que no termine en guión después del truncamiento
		slug = strings.TrimSuffix(slug, "-")
	}

	return slug
}

// NormalizeText normaliza un texto eliminando espacios extra y convirtiendo a minúsculas
func NormalizeText(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}

// NormalizeEmail normaliza un email (minúsculas y sin espacios)
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
