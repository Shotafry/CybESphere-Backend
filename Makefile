# CybESphere Backend - Makefile
# =============================
# Comandos para desarrollo, build, testing y deployment

.PHONY: help setup deps deps-test dev dev-watch build build-linux \
        test test-unit test-integration test-coverage test-race test-benchmark \
        docker-up docker-down docker-restart docker-logs docker-clean docker-test-db \
        db-create db-drop db-reset db-migrate db-seed db-seed-force db-seed-fresh \
        lint lint-fix format vet quality security \
        logs logs-clear logs-test clean clean-all update mod-verify

APP_NAME=cybesphere-backend
DOCKER_COMPOSE_DEV=docker-compose.dev.yml
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# -------------------------
# HELP
# -------------------------
help: ## Mostrar ayuda disponible
	@echo "CybESphere Backend - Comandos disponibles:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / \
	{printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# -------------------------
# SETUP & DEPENDENCIAS
# -------------------------
setup: ## Setup inicial del proyecto
	@echo "Configurando entorno..."
	@if [ ! -f .env ]; then cp .env.example .env; echo ".env creado"; fi
	@go mod download && go mod tidy
	@mkdir -p logs uploads migrations scripts/db/init
	@echo "Setup completado"

deps: ## Descargar dependencias de Go
	@go mod download && go mod tidy
	@echo "Dependencias actualizadas"

deps-test: ## Instalar dependencias para testing
	@go get github.com/stretchr/testify/{assert,require}
	@go get gorm.io/driver/sqlite
	@go mod tidy
	@echo "Dependencias de testing listas"

# -------------------------
# DESARROLLO
# -------------------------
dev: ## Ejecutar en modo desarrollo
	@go run ./cmd/cybesphere

dev-watch: ## Ejecutar con hot reload (Air)
	@air

# -------------------------
# BUILD
# -------------------------
build: ## Compilar la aplicación
	@go build -o bin/$(APP_NAME) ./cmd/cybesphere

build-linux: ## Compilar para Linux
	@GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux ./cmd/cybesphere

# -------------------------
# TESTING
# -------------------------
test: ## Ejecutar todos los tests
	@go test -v ./...

test-unit: ## Ejecutar tests unitarios
	@go test -v -short ./...

test-integration: ## Ejecutar tests de integración
	@go test -v -run Integration ./...

test-coverage: ## Ejecutar tests con coverage
	@go test -v -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

test-race: ## Detectar race conditions
	@go test -race ./...

test-benchmark: ## Ejecutar benchmarks
	@go test -bench=. -benchmem ./...

test-all: test-unit test-integration test-race test-coverage ## Suite completa de testing

# -------------------------
# DOCKER DESARROLLO LOCAL
# -------------------------
docker-local-up: ## Levantar servicios de Docker Compose
	@docker-compose -f $(DOCKER_COMPOSE_DEV) up -d

docker-local-down: ## Detener servicios
	@docker-compose -f $(DOCKER_COMPOSE_DEV) down

docker-local-restart: docker-local-down docker-local-up ## Reiniciar servicios Docker

docker-local-logs: ## Ver logs de servicios
	@docker-compose -f $(DOCKER_COMPOSE_DEV) logs -f

docker-local-clean: ## Limpiar contenedores e imágenes
	@docker-compose -f $(DOCKER_COMPOSE_DEV) down -v --remove-orphans
	@docker system prune -f

# -------------------------
# DATABASE
# -------------------------
db-create: ## Crear base de datos
	@createdb -h localhost -U cybesphere_user cybesphere_dev

db-drop: ## Eliminar base de datos
	@dropdb -h localhost -U cybesphere_user cybesphere_dev

db-reset: db-drop db-create ## Reset completo de base de datos

db-migrate: ## Ejecutar migraciones
	@go run ./cmd/cybesphere migrate

# -------------------------
# SEEDERS
# -------------------------
db-seed: ## Ejecutar seeders
	@go run ./cmd/seed/main.go

db-seed-fresh: ## Recrear tablas y seeders
	@go run ./cmd/seed/main.go -fresh

db-seed-list: ## Listar seeders disponibles
	@go run ./cmd/seed/main.go -list

seed-development: ## Seeders para desarrollo
	@go run ./cmd/seed/main.go -fresh

seed-testing: ## Seeders para testing
	@go run ./cmd/seed/main.go -fresh -priority 3

# -------------------------
# CALIDAD DE CÓDIGO
# -------------------------
lint: ## Ejecutar linter
	@golangci-lint run

lint-fix: ## Linter con auto-fix
	@golangci-lint run --fix

format: ## Formatear código
	@go fmt ./...
	@goimports -w . 2>/dev/null || true

vet: ## Ejecutar go vet
	@go vet ./...

quality: format vet lint test-unit ## Verificación completa

security: ## Verificar vulnerabilidades
	@gosec ./... 2>/dev/null || echo "gosec no instalado"

# -------------------------
# LOGS & UTILIDADES
# -------------------------
logs: ## Ver logs de aplicación
	@tail -f logs/app.log

logs-clear: ## Limpiar logs
	@rm -f logs/*.log

# -------------------------
# LIMPIEZA
# -------------------------
clean: ## Limpiar archivos generados
	@rm -rf bin/ $(COVERAGE_FILE) $(COVERAGE_HTML) test_*.db
	@go clean -testcache

clean-all: clean docker-clean ## Limpieza completa

update: ## Actualizar dependencias
	@go get -u ./...
	@go mod tidy

mod-verify: ## Verificar integridad de módulos
	@go mod verify
