#!/bin/bash

# CybESphere Backend - Script de Setup Inicial
# Este script prepara el entorno de desarrollo completo

set -e  # Exit on any error

# Colors para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logo ASCII
echo -e "${BLUE}"
echo " CYBERSPHERE BACKEND SETUP"
echo " \_______________________/ "
echo -e "${NC}"
echo -e "${GREEN}CybESphere Backend - Setup Inicial${NC}"
echo ""

# Función para mostrar mensajes
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✅]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[⚠️]${NC} $1"
}

log_error() {
    echo -e "${RED}[❌]${NC} $1"
}

# Verificar prerequisitos
check_prerequisites() {
    log_info "Verificando prerequisitos..."
    
    # Verificar Go
    if ! command -v go &> /dev/null; then
        log_error "Go no está instalado. Instala Go 1.21+ desde https://golang.org/dl/"
        exit 1
    fi
    
    GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
    if [[ $(echo "$GO_VERSION < 1.24" | bc) -eq 1 ]]; then
        log_error "Se requiere Go 1.24+. Versión actual: $GO_VERSION"
        exit 1
    fi
    log_success "Go $GO_VERSION encontrado"
    
    # Verificar Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker no está instalado. Instala Docker desde https://docker.com/get-started"
        exit 1
    fi
    log_success "Docker encontrado"
    
    # Verificar Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose no está instalado. Instala Docker Compose"
        exit 1
    fi
    log_success "Docker Compose encontrado"
    
    # Verificar Make
    if ! command -v make &> /dev/null; then
        log_warning "Make no está instalado. Se recomienda instalarlo para usar los comandos del Makefile"
    else
        log_success "Make encontrado"
    fi
}

# Crear estructura de directorios
create_directories() {
    log_info "Creando estructura de directorios..."
    
    directories=(
        "cmd"
        "internal/handlers"
        "internal/services"
        "internal/repository"
        "internal/models"
        "internal/middleware"
        "internal/config"
        "internal/utils"
        "pkg/auth"
        "pkg/database"
        "pkg/logger"
        "pkg/validator"
        "migrations"
        "scripts/db/init"
        "scripts/deploy"
        "logs"
        "uploads"
        "tmp"
        "docs"
        "tests/integration"
        "tests/unit"
    )
    
    for dir in "${directories[@]}"; do
        mkdir -p "$dir"
        log_success "Directorio creado: $dir"
    done
}

# Configurar archivo .env
setup_environment() {
    log_info "Configurando variables de entorno..."
    
    if [ ! -f .env ]; then
        if [ -f .env.example ]; then
            cp .env.example .env
            log_success "Archivo .env creado desde .env.example"
            
            # Generar JWT secret aleatorio
            JWT_SECRET=$(openssl rand -hex 32 2>/dev/null || head -c 32 /dev/urandom | base64)
            sed -i.bak "s/your_super_secret_jwt_key_minimum_32_characters_long/$JWT_SECRET/" .env
            
            # Generar password seguro para DB
            DB_PASSWORD=$(openssl rand -hex 16 2>/dev/null || head -c 16 /dev/urandom | base64)
            sed -i.bak "s/your_secure_password_here/$DB_PASSWORD/" .env
            
            # Actualizar docker-compose con la misma password
            if [ -f docker-compose.dev.yml ]; then
                sed -i.bak "s/your_secure_password_here/$DB_PASSWORD/" docker-compose.dev.yml
            fi
            
            rm -f .env.bak docker-compose.dev.yml.bak
            
            log_success "Credenciales generadas automáticamente"
            log_warning "Revisa y ajusta las variables en .env según tu entorno"
        else
            log_error ".env.example no encontrado"
            exit 1
        fi
    else
        log_warning ".env ya existe, omitiendo..."
    fi
}

# Instalar dependencias Go
install_go_dependencies() {
    log_info "Instalando dependencias de Go..."
    
    go mod init cybesphere-backend 2>/dev/null || true
    go mod download
    go mod tidy
    
    log_success "Dependencias de Go instaladas"
}

# Instalar herramientas de desarrollo
install_dev_tools() {
    log_info "Instalando herramientas de desarrollo..."
    
    # Air para hot reload
    if ! command -v air &> /dev/null; then
        log_info "Instalando Air (hot reload)..."
        go install github.com/cosmtrek/air@latest
        log_success "Air instalado"
    else
        log_success "Air ya está instalado"
    fi
    
    # golangci-lint para linting
    if ! command -v golangci-lint &> /dev/null; then
        log_info "Instalando golangci-lint..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
        log_success "golangci-lint instalado"
    else
        log_success "golangci-lint ya está instalado"
    fi
    
    # goimports para formateo
    if ! command -v goimports &> /dev/null; then
        log_info "Instalando goimports..."
        go install golang.org/x/tools/cmd/goimports@latest
        log_success "goimports instalado"
    else
        log_success "goimports ya está instalado"
    fi
}

# Configurar Docker
setup_docker() {
    log_info "Configurando entorno Docker..."
    
    # Verificar que Docker esté corriendo
    if ! docker info &> /dev/null; then
        log_error "Docker no está corriendo. Inicia Docker y vuelve a ejecutar este script"
        exit 1
    fi
    
    # Levantar servicios de desarrollo
    log_info "Levantando servicios de desarrollo..."
    docker-compose -f docker-compose.dev.yml up -d
    
    # Esperar a que PostgreSQL esté listo
    log_info "Esperando a que PostgreSQL esté listo..."
    sleep 15
    
    # Verificar conexión a la base de datos
    if docker-compose -f docker-compose.dev.yml exec -T postgres pg_isready -U cybesphere_user -d cybesphere_dev &> /dev/null; then
        log_success "PostgreSQL está listo"
    else
        log_warning "PostgreSQL puede no estar completamente listo. Verifica manualmente"
    fi
    
    log_success "Servicios Docker configurados"
}

# Crear archivo de configuración golangci-lint
create_golangci_config() {
    log_info "Creando configuración de golangci-lint..."
    
    cat > .golangci.yml << 'EOF'
run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode
    - typecheck
    - gosec
    - gocyclo
    - misspell
    - unparam

linters-settings:
  gocyclo:
    min-complexity: 15
  gosec:
    excludes:
      - G104 # Audit errors not checked

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gosec
        - errcheck
EOF
    
    log_success "Configuración de golangci-lint creada"
}

# Mostrar información final
show_final_info() {
    echo ""
    echo -e "${GREEN}🎉 Setup completado exitosamente!${NC}"
    echo ""
    echo -e "${BLUE}Servicios disponibles:${NC}"
    echo "  📦 PostgreSQL:     localhost:5432"
    echo "  🔴 Redis:          localhost:6379"
    echo "  🛠️  pgAdmin:        http://localhost:5050"
    echo "  📖 Docmost:        http://localhost:3000"
    echo ""
    echo -e "${BLUE}Próximos pasos:${NC}"
    echo "  1. Revisa y ajusta las variables en .env"
    echo "  2. Ejecuta 'make dev' para iniciar el servidor"
    echo "  3. O ejecuta 'make dev-watch' para desarrollo con hot reload"
    echo ""
    echo -e "${BLUE}Comandos útiles:${NC}"
    echo "  make help          - Ver todos los comandos disponibles"
    echo "  make docker-up     - Levantar servicios Docker"
    echo "  make docker-down   - Detener servicios Docker"
    echo "  make test          - Ejecutar tests"
    echo "  make build         - Compilar aplicación"
    echo ""
    echo -e "${YELLOW}Credenciales por defecto:${NC}"
    echo "  pgAdmin: admin@cybesphere.local / admin123"
    echo "  Docmost: Configurar en primer acceso"
    echo ""
    echo -e "${GREEN}¡Listo para desarrollar! 🚀${NC}"
}

# Función principal
main() {
    echo "Iniciando setup de CybESphere Backend..."
    echo ""
    
    check_prerequisites
    create_directories
    setup_environment
    install_go_dependencies
    install_dev_tools
    create_golangci_config
    setup_docker
    show_final_info
}

# Ejecutar si se llama directamente
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi