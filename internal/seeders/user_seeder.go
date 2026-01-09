package seeders

import (
	"fmt"
	"strings"

	"cybesphere-backend/internal/models"
	"cybesphere-backend/pkg/logger"
	"cybesphere-backend/pkg/utils"

	"gorm.io/gorm"
)

type UserSeeder struct{}

func NewUserSeeder() *UserSeeder {
	return &UserSeeder{}
}

func (us *UserSeeder) Name() string {
	return "UserSeeder"
}

func (us *UserSeeder) Description() string {
	return "Crea usuarios de ejemplo: admin, organizadores y usuarios regulares"
}

func (us *UserSeeder) Priority() int {
	return 1 // Ejecutar primero
}

func (us *UserSeeder) CanRun(db *gorm.DB) bool {
	var count int64
	db.Model(&models.User{}).Count(&count)
	return count == 0
}

func (us *UserSeeder) Seed(db *gorm.DB) error {
	logger.Info("🌱 Creando usuarios de ejemplo...")

	// 1. Usuario administrador principal
	adminUser := &models.User{
		Email:             "admin@cybesphere.com",
		Password:          "admin123456", // Se hasheará automáticamente
		FirstName:         "Administrator",
		LastName:          "CybESphere",
		Role:              models.RoleAdmin,
		IsActive:          true,
		IsVerified:        true,
		Company:           "CybESphere",
		Position:          "System Administrator",
		Bio:               "Administrador principal de la plataforma CybESphere",
		City:              "Madrid",
		Country:           "Spain",
		Timezone:          "Europe/Madrid",
		Language:          "es",
		NewsletterEnabled: true,
	}
	adminUser.SetLocation(40.4168, -3.7038, "Madrid", "Spain") // Coordenadas de Madrid

	if err := db.Create(adminUser).Error; err != nil {
		return err
	}

	// 2. Usuarios organizadores
	organizers := []*models.User{
		{
			Email:             "maria.garcia@cybersec-spain.org",
			Password:          "organizer123456",
			FirstName:         "María",
			LastName:          "García",
			Role:              models.RoleOrganizer,
			IsActive:          true,
			IsVerified:        true,
			Company:           "CyberSecurity Spain",
			Position:          "Event Manager",
			Bio:               "Especialista en gestión de eventos de ciberseguridad",
			Website:           "https://linkedin.com/in/maria-garcia-cybersec",
			LinkedIn:          "https://linkedin.com/in/maria-garcia-cybersec",
			City:              "Madrid",
			Country:           "Spain",
			Timezone:          "Europe/Madrid",
			Language:          "es",
			NewsletterEnabled: true,
		},
		{
			Email:             "carlos.rodriguez@hackingetic.com",
			Password:          "organizer123456",
			FirstName:         "Carlos",
			LastName:          "Rodríguez",
			Role:              models.RoleOrganizer,
			IsActive:          true,
			IsVerified:        true,
			Company:           "Hackingétic",
			Position:          "CTO",
			Bio:               "Experto en ethical hacking y formación en ciberseguridad",
			Website:           "https://hackingetic.com",
			Twitter:           "https://twitter.com/carlos_cybersec",
			LinkedIn:          "https://linkedin.com/in/carlos-rodriguez-cybersec",
			City:              "Barcelona",
			Country:           "Spain",
			Timezone:          "Europe/Madrid",
			Language:          "es",
			NewsletterEnabled: true,
		},
		{
			Email:             "ana.martinez@secureit.es",
			Password:          "organizer123456",
			FirstName:         "Ana",
			LastName:          "Martínez",
			Role:              models.RoleOrganizer,
			IsActive:          true,
			IsVerified:        true,
			Company:           "SecureIT",
			Position:          "Security Consultant",
			Bio:               "Consultora especializada en auditorías de seguridad",
			Website:           "https://secureit.es",
			LinkedIn:          "https://linkedin.com/in/ana-martinez-security",
			City:              "Valencia",
			Country:           "Spain",
			Timezone:          "Europe/Madrid",
			Language:          "es",
			NewsletterEnabled: true,
		},
	}

	// Añadir coordenadas a organizadores
	organizers[0].SetLocation(40.4168, -3.7038, "Madrid", "Spain")   // Madrid
	organizers[1].SetLocation(41.3851, 2.1734, "Barcelona", "Spain") // Barcelona
	organizers[2].SetLocation(39.4699, -0.3763, "Valencia", "Spain") // Valencia

	for _, organizer := range organizers {
		if err := db.Create(organizer).Error; err != nil {
			return err
		}
	}

	// 3. Usuarios regulares variados
	regularUsers := []*models.User{
		{
			Email:             "juan.perez@techcorp.com",
			Password:          "user123456",
			FirstName:         "Juan",
			LastName:          "Pérez",
			Role:              models.RoleUser,
			IsActive:          true,
			IsVerified:        true,
			Company:           "TechCorp",
			Position:          "Security Analyst",
			Bio:               "Analista de seguridad con 3 años de experiencia",
			City:              "Sevilla",
			Country:           "Spain",
			Timezone:          "Europe/Madrid",
			Language:          "es",
			NewsletterEnabled: true,
		},
		{
			Email:             "laura.gomez@university.edu",
			Password:          "user123456",
			FirstName:         "Laura",
			LastName:          "Gómez",
			Role:              models.RoleUser,
			IsActive:          true,
			IsVerified:        true,
			Company:           "Universidad Complutense",
			Position:          "Estudiante de Doctorado",
			Bio:               "Investigando en criptografía post-cuántica",
			City:              "Madrid",
			Country:           "Spain",
			Timezone:          "Europe/Madrid",
			Language:          "es",
			NewsletterEnabled: true,
		},
		{
			Email:             "miguel.santos@freelancer.com",
			Password:          "user123456",
			FirstName:         "Miguel",
			LastName:          "Santos",
			Role:              models.RoleUser,
			IsActive:          true,
			IsVerified:        false, // Usuario no verificado para testing
			Company:           "Freelance",
			Position:          "Pentester",
			Bio:               "Pentester independiente especializado en aplicaciones web",
			Website:           "https://miguel-santos-security.com",
			Twitter:           "https://twitter.com/miguel_pentest",
			City:              "Bilbao",
			Country:           "Spain",
			Timezone:          "Europe/Madrid",
			Language:          "es",
			NewsletterEnabled: false,
		},
		{
			Email:             "sofia.ruiz@startup.io",
			Password:          "user123456",
			FirstName:         "Sofía",
			LastName:          "Ruiz",
			Role:              models.RoleUser,
			IsActive:          true,
			IsVerified:        true,
			Company:           "CyberStartup",
			Position:          "DevSecOps Engineer",
			Bio:               "Ingeniera DevSecOps apasionada por la automatización de seguridad",
			LinkedIn:          "https://linkedin.com/in/sofia-ruiz-devsecops",
			City:              "Barcelona",
			Country:           "Spain",
			Timezone:          "Europe/Madrid",
			Language:          "es",
			NewsletterEnabled: true,
		},
		{
			Email:             "inactive@example.com",
			Password:          "user123456",
			FirstName:         "Usuario",
			LastName:          "Inactivo",
			Role:              models.RoleUser,
			IsActive:          false, // Usuario inactivo para testing
			IsVerified:        true,
			Company:           "Ex-Company",
			Position:          "Ex-Position",
			Bio:               "Usuario desactivado para pruebas de funcionalidad",
			City:              "Madrid",
			Country:           "Spain",
			Timezone:          "Europe/Madrid",
			Language:          "es",
			NewsletterEnabled: false,
		},
	}

	// Coordenadas para usuarios regulares
	coordinates := []struct {
		lat, lng float64
		city     string
	}{
		{37.3886, -5.9823, "Sevilla"},  // Sevilla
		{40.4168, -3.7038, "Madrid"},   // Madrid
		{43.2627, -2.9253, "Bilbao"},   // Bilbao
		{41.3851, 2.1734, "Barcelona"}, // Barcelona
		{40.4168, -3.7038, "Madrid"},   // Madrid (usuario inactivo)
	}

	for i, user := range regularUsers {
		coord := coordinates[i]
		user.SetLocation(coord.lat, coord.lng, coord.city, "Spain")

		if err := db.Create(user).Error; err != nil {
			return err
		}
	}

	// 4. Generar usuarios adicionales aleatorios
	if err := us.generateRandomUsers(db, 25); err != nil {
		return err
	}

	logger.Info("✅ Usuarios creados exitosamente")
	return nil
}

// generateRandomUsers crea usuarios adicionales con datos aleatorios
func (us *UserSeeder) generateRandomUsers(db *gorm.DB, count int) error {
	firstNames := []string{
		"Alejandro", "Beatriz", "Carmen", "Diego", "Elena", "Fernando", "Gloria", "Héctor",
		"Isabel", "Javier", "Lucía", "Manuel", "Natalia", "Óscar", "Patricia", "Raúl",
		"Silvia", "Tomás", "Verónica", "Xavier", "Yolanda", "Zacarías", "Amparo", "Bernardo",
		"Cristina", "Daniel", "Esperanza", "Francisco", "Guadalupe", "Ignacio",
	}

	lastNames := []string{
		"García", "Rodríguez", "González", "Fernández", "López", "Martínez", "Sánchez", "Pérez",
		"Gómez", "Martín", "Jiménez", "Ruiz", "Hernández", "Díaz", "Moreno", "Muñoz",
		"Álvarez", "Romero", "Alonso", "Gutiérrez", "Navarro", "Torres", "Domínguez", "Vázquez",
		"Ramos", "Gil", "Ramírez", "Serrano", "Blanco", "Suárez",
	}

	companies := []string{
		"IBM", "Accenture", "Deloitte", "Telefónica", "BBVA", "Santander", "Indra", "Capgemini",
		"Everis", "Atos", "Sopra Steria", "NTT Data", "DXC Technology", "Bankia", "Mapfre",
		"Repsol", "Endesa", "Ferrovial", "ACS", "Amadeus", "Freelance", "Startup", "Consulting",
	}

	positions := []string{
		"Security Analyst", "Cybersecurity Specialist", "IT Security Consultant", "SOC Analyst",
		"Penetration Tester", "Security Engineer", "CISO", "Security Architect", "Forensic Analyst",
		"Incident Response Analyst", "Compliance Officer", "Risk Analyst", "Security Researcher",
		"DevSecOps Engineer", "Security Auditor", "Network Security Engineer", "Cloud Security Specialist",
	}

	cities := []struct {
		name     string
		lat, lng float64
	}{
		{"Madrid", 40.4168, -3.7038},
		{"Barcelona", 41.3851, 2.1734},
		{"Valencia", 39.4699, -0.3763},
		{"Sevilla", 37.3886, -5.9823},
		{"Zaragoza", 41.6518, -0.8759},
		{"Málaga", 36.7213, -4.4214},
		{"Murcia", 37.9922, -1.1307},
		{"Palma", 39.5696, 2.6502},
		{"Bilbao", 43.2627, -2.9253},
		{"Alicante", 38.3460, -0.4907},
	}

	for i := 0; i < count; i++ {
		firstName := firstNames[utils.SecureRandInt(len(firstNames))]
		lastName := lastNames[utils.SecureRandInt(len(lastNames))]
		email := fmt.Sprintf("%s.%s%d@example.com",
			strings.ToLower(firstName), strings.ToLower(lastName), utils.SecureRandInt(100))

		company := companies[utils.SecureRandInt(len(companies))]
		position := positions[utils.SecureRandInt(len(positions))]
		city := cities[utils.SecureRandInt(len(cities))]

		user := &models.User{
			Email:             email,
			Password:          "user123456",
			FirstName:         firstName,
			LastName:          lastName,
			Role:              models.RoleUser,
			IsActive:          utils.SecureRandFloat32() > 0.1, // 90% activos
			IsVerified:        utils.SecureRandFloat32() > 0.3, // 70% verificados
			Company:           company,
			Position:          position,
			City:              city.name,
			Country:           "Spain",
			Timezone:          "Europe/Madrid",
			Language:          "es",
			NewsletterEnabled: utils.SecureRandFloat32() > 0.4, // 60% con newsletter
		}

		user.SetLocation(city.lat, city.lng, city.name, "Spain")

		// Algunos usuarios con bio y redes sociales
		if utils.SecureRandFloat32() > 0.7 {
			user.Bio = fmt.Sprintf("%s especializado en ciberseguridad con experiencia en %s",
				position, company)
		}

		if utils.SecureRandFloat32() > 0.8 {
			user.LinkedIn = fmt.Sprintf("https://linkedin.com/in/%s-%s",
				strings.ToLower(firstName), strings.ToLower(lastName))
		}

		if err := db.Create(user).Error; err != nil {
			return err
		}
	}

	return nil
}
