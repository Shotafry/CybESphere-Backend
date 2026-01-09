package seeders

import (
	"fmt"
	"strings"
	"time"

	"cybesphere-backend/internal/models"
	"cybesphere-backend/pkg/logger"
	"cybesphere-backend/pkg/utils"

	"gorm.io/gorm"
)

type OrganizationSeeder struct{}

func NewOrganizationSeeder() *OrganizationSeeder {
	return &OrganizationSeeder{}
}

func (os *OrganizationSeeder) Name() string {
	return "OrganizationSeeder"
}

func (os *OrganizationSeeder) Description() string {
	return "Crea organizaciones de ejemplo y asigna organizadores"
}

func (os *OrganizationSeeder) Priority() int {
	return 2 // Ejecutar después de usuarios
}

func (os *OrganizationSeeder) CanRun(db *gorm.DB) bool {
	var count int64
	db.Model(&models.Organization{}).Count(&count)
	return count == 0
}

func (os *OrganizationSeeder) Seed(db *gorm.DB) error {
	logger.Info("🏢 Creando organizaciones de ejemplo...")

	// 1. Organizaciones principales verificadas
	mainOrganizations := []*models.Organization{
		{
			Name:        "CyberSecurity Spain",
			Description: "La mayor comunidad española de profesionales de ciberseguridad. Organizamos eventos, conferencias y meetups para conectar a expertos y principiantes en el mundo de la ciberseguridad.",
			Email:       "info@cybersecurityspain.org",
			Website:     "https://cybersecurityspain.org",
			Phone:       "+34 91 123 45 67",
			Address:     "Calle Gran Vía, 28",
			City:        "Madrid",
			Country:     "Spain",
			PostalCode:  "28013",
			Status:      models.OrgStatusActive,
			IsVerified:  true,
			LogoURL:     "/CloudEvents-logo-2@2x.png",
			BannerURL:   "https://images.unsplash.com/photo-1550751827-4bd374c3f58b?auto=format&fit=crop&w=1600&q=80",
			LinkedIn:    "https://linkedin.com/company/cybersecurity-spain",
			Twitter:     "https://twitter.com/cybersec_spain",
			TaxID:       "G12345678",
			LegalName:   "Asociación CyberSecurity Spain",
		},
		{
			Name:        "Hackingétic",
			Description: "Centro de formación especializado en ethical hacking y pentesting. Ofrecemos cursos avanzados y certificaciones reconocidas internacionalmente.",
			Email:       "contacto@hackingetic.com",
			Website:     "https://hackingetic.com",
			Phone:       "+34 93 987 65 43",
			Address:     "Passeig de Gràcia, 15",
			City:        "Barcelona",
			Country:     "Spain",
			PostalCode:  "08007",
			Status:      models.OrgStatusActive,
			IsVerified:  true,
			LogoURL:     "/cyberLogo-gigapixel-art-scale-2-00x-godpix-1@2x.png",
			BannerURL:   "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?auto=format&fit=crop&w=1600&q=80",
			LinkedIn:    "https://linkedin.com/company/hackingetic",
			Twitter:     "https://twitter.com/hackingetic",
			Facebook:    "https://facebook.com/hackingetic",
			TaxID:       "B98765432",
			LegalName:   "Hackingétic Formación S.L.",
		},
		{
			Name:        "SecurIT All",
			Description: "Consultora especializada en auditorías de seguridad y consultoría en ciberseguridad para empresas medianas y grandes.",
			Email:       "info@securit-all.es",
			Website:     "https://securit-all.es",
			Phone:       "+34 96 321 54 87",
			Address:     "Av. del Puerto, 45",
			City:        "Valencia",
			Country:     "Spain",
			PostalCode:  "46021",
			Status:      models.OrgStatusActive,
			IsVerified:  true,
			LogoURL:     "/asturcon-low-1@2x.png",
			BannerURL:   "https://images.unsplash.com/photo-1563206767-5b1d972d9323?auto=format&fit=crop&w=1600&q=80",
			LinkedIn:    "https://linkedin.com/company/securit-all",
			TaxID:       "A11223344",
			LegalName:   "SecurIT All S.A.",
		},
		{
			Name:        "Universidad Complutense - Ciberseguridad",
			Description: "Departamento de Ciberseguridad de la Universidad Complutense de Madrid. Organizamos seminarios académicos y conferencias de investigación.",
			Email:       "cybersec@ucm.es",
			Website:     "https://ucm.es/cybersecurity",
			Phone:       "+34 91 394 52 00",
			Address:     "Ciudad Universitaria, s/n",
			City:        "Madrid",
			Country:     "Spain",
			PostalCode:  "28040",
			Status:      models.OrgStatusActive,
			IsVerified:  true,
			LogoURL:     "/CloudEvents-logo-1@2x.png",
			LinkedIn:    "https://linkedin.com/school/universidad-complutense-madrid",
			TaxID:       "Q2818018D",
			LegalName:   "Universidad Complutense de Madrid",
		},
	}

	// Asignar coordenadas y crear organizaciones principales
	coordinates := []struct{ lat, lng float64 }{
		{40.4168, -3.7038}, // Madrid
		{41.3851, 2.1734},  // Barcelona
		{39.4699, -0.3763}, // Valencia
		{40.4168, -3.7038}, // Madrid (UCM)
	}

	// Configurar branding para algunas organizaciones
	brandingColors := []struct{ primary, secondary string }{
		{"#FF6B35", "#004E89"}, // Naranja y Azul
		{"#2ECC71", "#34495E"}, // Verde y Gris
		{"#9B59B6", "#E74C3C"}, // Morado y Rojo
		{"#3498DB", "#2C3E50"}, // Azul y Azul Oscuro
	}

	for i, org := range mainOrganizations {
		coord := coordinates[i]
		org.SetLocation(coord.lat, coord.lng) // No retorna error

		// Configurar branding
		colors := brandingColors[i]
		if err := org.SetBranding(colors.primary, colors.secondary); err != nil {
			logger.Warnf("Error asignando branding a la organización %s: %v", org.Name, err)
		}

		// Marcar como verificadas con fecha
		now := time.Now().Format(time.RFC3339)
		org.VerifiedAt = &now
		adminUserID := "admin-user-id" // Se actualizará después
		org.VerifiedBy = &adminUserID

		if err := db.Create(org).Error; err != nil {
			return err
		}
	}

	// 2. Asignar organizadores a organizaciones principales
	if err := os.assignOrganizersToOrganizations(db); err != nil {
		return err
	}

	// 3. Crear organizaciones adicionales en diferentes estados
	additionalOrgs := []*models.Organization{
		{
			Name:        "InfoSec Sevilla",
			Description: "Grupo local de profesionales de seguridad informática en Sevilla. Organizamos meetups mensuales y talleres prácticos.",
			Email:       "hola@infosecsevilla.es",
			Website:     "https://infosecsevilla.es",
			City:        "Sevilla",
			Country:     "Spain",
			Status:      models.OrgStatusActive,
			IsVerified:  true,
		},
		{
			Name:        "CyberStartup Incubator",
			Description: "Incubadora de startups especializadas en ciberseguridad. Apoyamos a emprendedores con ideas innovadoras.",
			Email:       "info@cyberstartup.es",
			Website:     "https://cyberstartup.es",
			City:        "Barcelona",
			Country:     "Spain",
			Status:      models.OrgStatusActive,
			IsVerified:  false, // Organización nueva, pendiente de verificación
		},
		{
			Name:        "Red Team Málaga",
			Description: "Comunidad de red teamers y pentesters de Málaga. Nos enfocamos en técnicas avanzadas de hacking ético.",
			Email:       "contact@redteammalaga.com",
			City:        "Málaga",
			Country:     "Spain",
			Status:      models.OrgStatusPending, // Pendiente de revisión
			IsVerified:  false,
		},
		{
			Name:        "Org Suspendida Test",
			Description: "Organización de prueba para testing del estado suspendido.",
			Email:       "test@suspended.example.com",
			City:        "Madrid",
			Country:     "Spain",
			Status:      models.OrgStatusSuspended, // Para testing
			IsVerified:  false,
		},
		{
			Name:        "ForensicsLab Zaragoza",
			Description: "Laboratorio especializado en análisis forense digital y respuesta a incidentes.",
			Email:       "lab@forensicszgz.es",
			Website:     "https://forensicszgz.es",
			City:        "Zaragoza",
			Country:     "Spain",
			Status:      models.OrgStatusActive,
			IsVerified:  true,
		},
	}

	// Coordenadas para organizaciones adicionales
	additionalCoords := []struct{ lat, lng float64 }{
		{37.3886, -5.9823}, // Sevilla
		{41.3851, 2.1734},  // Barcelona
		{36.7213, -4.4214}, // Málaga
		{40.4168, -3.7038}, // Madrid
		{41.6518, -0.8759}, // Zaragoza
	}

	for i, org := range additionalOrgs {
		if i < len(additionalCoords) {
			coord := additionalCoords[i]
			org.SetLocation(coord.lat, coord.lng)
		}

		if err := db.Create(org).Error; err != nil {
			return err
		}
	}

	// 4. Generar organizaciones aleatorias
	if err := os.generateRandomOrganizations(db, 15); err != nil {
		return err
	}

	logger.Info("✅ Organizaciones creadas exitosamente")
	return nil
}

// assignOrganizersToOrganizations asigna usuarios organizadores a organizaciones
func (os *OrganizationSeeder) assignOrganizersToOrganizations(db *gorm.DB) error {
	// Mapear organizadores con organizaciones
	assignments := map[string]string{
		"maria.garcia@cybersec-spain.org":  "CyberSecurity Spain",
		"carlos.rodriguez@hackingetic.com": "Hackingétic",
		"ana.martinez@secureit.es":         "SecureIT Valencia",
	}

	for email, orgName := range assignments {
		// Buscar usuario organizador
		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			logger.Warnf("Usuario organizador no encontrado: %s", email)
			continue
		}

		// Buscar organización
		var org models.Organization
		if err := db.Where("name = ?", orgName).First(&org).Error; err != nil {
			logger.Warnf("Organización no encontrada: %s", orgName)
			continue
		}

		// Asignar usuario a organización
		orgIDStr := org.ID.String()
		user.OrganizationID = &orgIDStr

		if err := db.Save(&user).Error; err != nil {
			logger.Errorf("Error asignando usuario %s a organización %s: %v", email, orgName, err)
			continue
		}

		logger.Debugf("Usuario %s asignado a organización %s", email, orgName)
	}

	return nil
}

// generateRandomOrganizations crea organizaciones adicionales con datos aleatorios
func (os *OrganizationSeeder) generateRandomOrganizations(db *gorm.DB, count int) error {
	orgTypes := []string{
		"Asociación", "Consultora", "Universidad", "Centro de Formación",
		"Comunidad", "Startup", "Laboratorio", "Instituto",
	}

	focusAreas := []string{
		"Ethical Hacking", "Forensics", "Compliance", "Red Team", "Blue Team",
		"Cloud Security", "IoT Security", "Blockchain", "AI Security", "DevSecOps",
		"Industrial Security", "Mobile Security", "Network Security", "Cryptography",
	}

	cities := []struct {
		name     string
		lat, lng float64
	}{
		{"Bilbao", 43.2627, -2.9253},
		{"Palma", 39.5696, 2.6502},
		{"Murcia", 37.9922, -1.1307},
		{"Alicante", 38.3460, -0.4907},
		{"Córdoba", 37.8882, -4.7794},
		{"Valladolid", 41.6523, -4.7245},
		{"Vigo", 42.2406, -8.7207},
		{"Gijón", 43.5322, -5.6611},
		{"Vitoria", 42.8467, -2.6716},
		{"Santander", 43.4623, -3.8099},
	}

	statuses := []models.OrganizationStatus{
		models.OrgStatusActive,
		models.OrgStatusActive,
		models.OrgStatusActive, // Más probabilidad de activas
		models.OrgStatusPending,
		models.OrgStatusInactive,
	}

	for i := 0; i < count; i++ {
		orgType := orgTypes[utils.SecureRandInt(len(orgTypes))]
		focusArea := focusAreas[utils.SecureRandInt(len(focusAreas))]
		city := cities[utils.SecureRandInt(len(cities))]

		orgName := fmt.Sprintf("%s %s %s", orgType, focusArea, city.name)

		// Generar slug único
		baseSlug := strings.ToLower(strings.ReplaceAll(orgName, " ", "-"))
		baseSlug = strings.ReplaceAll(baseSlug, "ñ", "n")
		slug := fmt.Sprintf("%s-%d", baseSlug, i)

		org := &models.Organization{
			Name: orgName,
			Slug: slug,
			Description: fmt.Sprintf("%s especializado en %s ubicado en %s. Organizamos eventos y actividades relacionadas con la ciberseguridad.",
				orgType, focusArea, city.name),
			Email:      fmt.Sprintf("info@%s.es", strings.ReplaceAll(baseSlug, "-", "")),
			City:       city.name,
			Country:    "Spain",
			Status:     statuses[utils.SecureRandInt(len(statuses))],
			IsVerified: utils.SecureRandFloat32() > 0.4, // 60% verificadas
		}

		// Configurar website y redes sociales para algunas
		if utils.SecureRandFloat32() > 0.5 {
			org.Website = fmt.Sprintf("https://%s.es", strings.ReplaceAll(baseSlug, "-", ""))
		}

		if utils.SecureRandFloat32() > 0.7 {
			org.LinkedIn = fmt.Sprintf("https://linkedin.com/company/%s", baseSlug)
		}

		if utils.SecureRandFloat32() > 0.8 {
			org.Twitter = fmt.Sprintf("https://twitter.com/%s", strings.ReplaceAll(baseSlug, "-", ""))
		}

		// Configurar ubicación
		org.SetLocation(city.lat, city.lng)

		// Configurar algunos con branding aleatorio
		if utils.SecureRandFloat32() > 0.6 {
			colors := []string{"#FF6B35", "#2ECC71", "#9B59B6", "#3498DB", "#E74C3C", "#F39C12"}
			primary := colors[utils.SecureRandInt(len(colors))]
			secondary := colors[utils.SecureRandInt(len(colors))]

			if err := org.SetBranding(primary, secondary); err != nil {
				logger.Warnf("Error asignando branding a la organización %s: %v", org.Name, err)
			}
		}

		if err := db.Create(org).Error; err != nil {
			return err
		}
	}

	return nil
}
