package main

import (
	"fmt"
	"os"

	"cybesphere-backend/internal/config"
)

func main() {
	fmt.Println("🔍 Debugging Configuration...")
	fmt.Println("==============================")

	// Mostrar variables de entorno relevantes
	fmt.Println("\n📋 Environment Variables:")
	envVars := []string{
		"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD",
		"SERVER_HOST", "SERVER_PORT", "JWT_SECRET",
	}

	for _, env := range envVars {
		value := os.Getenv(env)
		if value == "" {
			fmt.Printf("❌ %s: NOT SET\n", env)
		} else {
			// Ocultar passwords
			if env == "DB_PASSWORD" || env == "JWT_SECRET" {
				fmt.Printf("✅ %s: [HIDDEN - %d chars]\n", env, len(value))
			} else {
				fmt.Printf("✅ %s: %s\n", env, value)
			}
		}
	}

	fmt.Println("\n🔧 Loading Configuration...")

	// Intentar cargar configuración
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Configuration loaded successfully!")

	// Mostrar configuración (sin passwords)
	fmt.Println("\n📊 Configuration Summary:")
	fmt.Printf("  Server: %s:%s (mode: %s)\n",
		cfg.Server.Host, fmt.Sprintf("%d", cfg.Server.Port), cfg.Server.Mode)
	fmt.Printf("  Database: %s@%s:%s/%s (ssl: %s)\n",
		cfg.Database.User, cfg.Database.Host, fmt.Sprintf("%d", cfg.Database.Port),
		cfg.Database.Name, cfg.Database.SSLMode)
	fmt.Printf("  JWT Issuer: %s (exp: %s)\n",
		cfg.JWT.Issuer, cfg.JWT.AccessTokenDuration)
	fmt.Printf("  Environment: %s\n", cfg.Monitoring.Environment)
	fmt.Printf("  Logging: %s (%s)\n", cfg.Logging.Level, cfg.Logging.Format)

	fmt.Println("\n🔗 Database DSN (without password):")
	// Crear versión segura del DSN para mostrar
	safeDSN := fmt.Sprintf("host=%s user=%s password=[HIDDEN] dbname=%s port=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Name,
		fmt.Sprintf("%d", cfg.Database.Port), cfg.Database.SSLMode)
	fmt.Printf("  %s\n", safeDSN)

	fmt.Println("\n✅ Configuration debug completed!")
}
