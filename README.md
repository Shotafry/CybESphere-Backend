# CybESphere Backend

Backend para plataforma de gestión de eventos de ciberseguridad desarrollado en Go con Clean Architecture y PostgreSQL.

## Estado Actual - Versión 0.1.0

### Implementado en v0.1.0

- **Configuración completa** con Viper y variables de entorno
- **Conexión a PostgreSQL** con GORM y pool de conexiones optimizado
- **Sistema de logging** estructurado con logrus y niveles configurables
- **Servidor HTTP** con Gin framework y middleware de seguridad
- **Modelos de dominio completos**: User, Organization, Event con relaciones
- **Sistema de autenticación JWT** con refresh tokens y rotación segura
- **RBAC completo** con 3 roles: Admin, Organizer, Attendee
- **Repositorios y servicios** implementados con Clean Architecture
- **Migraciones automáticas** con GORM AutoMigrate
- **Sistema de seeders** con datos realistas de ciberseguridad
- **Middleware de seguridad**: CORS, Rate Limiting, Request ID, Recovery
- **Endpoints completos** de autenticación, usuarios, organizaciones y eventos
- **Testing comprehensivo** unitario e integración con base de datos en memoria
- **Docker Compose** para servicios de desarrollo
- **Hot reload** configurado con Air
- **Sistema de build** y comandos con Makefile extendido

### Próximas funcionalidades (v0.2.0)

- TBD

## Prerequisitos

- Go 1.24+
- Docker y Docker Compose
- Make (recomendado)

## Instalación y Setup

### 1. Clonar e instalar dependencias

```bash
git clone <repository-url>
cd cybesphere-backend

# Instalar dependencias Go
go mod download
go mod tidy

# Instalar herramientas de desarrollo
go install github.com/air-verse/air@latest
```

### 2. Configurar variables de entorno

```bash
# Copiar archivo de ejemplo
cp .env.example .env

# Editar .env con tus configuraciones
# Variables principales:
# - DB_PASSWORD (debe coincidir con docker-compose.dev.yml)
# - JWT_SECRET (mínimo 32 caracteres)
# - JWT_REFRESH_SECRET (diferente al JWT_SECRET)
# - CORS_ALLOWED_ORIGINS (dominios permitidos)
```

### 3. Levantar servicios Docker

```bash
# Levantar PostgreSQL y servicios de desarrollo
make docker-up

# Verificar que PostgreSQL esté listo (esperar ~30 segundos)
make docker-logs
```

### 4. Inicializar base de datos

```bash
# Ejecutar migraciones y seeders
make db-setup

# O paso a paso:
make db-migrate    # Crear tablas
make db-seed       # Datos de ejemplo
```

### 5. Ejecutar aplicación

```bash
# Modo desarrollo con hot reload
make dev-watch

# O modo desarrollo normal
make dev

# O manualmente:
go run ./cmd/cybesphere
```

## Usuarios de Prueba

Después de ejecutar los seeders, tendrás acceso a estos usuarios:

| Email                      | Password      | Rol       | Descripción                         |
| -------------------------- | ------------- | --------- | ----------------------------------- |
| admin@cybesphere.local     | Admin123!     | Admin     | Acceso completo al sistema          |
| organizer@cybesphere.local | Organizer123! | Organizer | Gestión de eventos y organizaciones |
| attendee@cybesphere.local  | Attendee123!  | Attendee  | Participación en eventos            |

## Servicios Disponibles

| Servicio    | URL                   | Descripción          |
| ----------- | --------------------- | -------------------- |
| API Backend | http://localhost:8080 | Aplicación principal |
| PostgreSQL  | localhost:5432        | Base de datos        |
| pgAdmin     | http://localhost:5050 | Interfaz admin DB    |
| Redis       | localhost:6379        | Cache (opcional)     |
| Docmost     | http://localhost:3000 | Wiki documentación   |

### Credenciales por defecto

- **pgAdmin**: admin@cybesphere.local / admin123
- **PostgreSQL**: cybesphere_user / (ver .env)

## Documentación de la API

La documentación completa de todos los endpoints de la API está disponible en:

**[ Documentación de la API - `/docs/api/README.md`](./docs/api/README.md)**

Esta documentación incluye:

- Endpoints de autenticación y gestión de sesiones
- CRUD completo de eventos con filtros avanzados
- Gestión de organizaciones y verificación
- Administración de usuarios y permisos
- Ejemplos de requests/responses
- Códigos de error y validaciones
- Guías de uso y mejores prácticas

### Quick Start API

```bash
# Health check
curl http://localhost:8080/health

# Registro de usuario
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!","first_name":"Test","last_name":"User"}'

# Login (obtener token)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!"}'

# Ver eventos públicos
curl http://localhost:8080/api/v1/public/events
```

Para más ejemplos y documentación completa, consulta [`/docs/api/README.md`](./docs/api/README.md).

## Comandos de Desarrollo

```bash
# Desarrollo
make dev              # Ejecutar servidor
make dev-watch        # Hot reload con Air
make build            # Compilar aplicación
make test             # Tests unitarios
make test-integration # Tests de integración
make test-all         # Todos los tests

# Docker
make docker-up        # Levantar servicios
make docker-down      # Detener servicios
make docker-restart   # Reiniciar servicios
make docker-logs      # Ver logs

# Base de datos
make db-setup         # Migrar + Seed completo
make db-migrate       # Solo migraciones
make db-seed          # Solo datos de ejemplo
make db-seed-users    # Solo usuarios
make db-seed-orgs     # Solo organizaciones
make db-seed-events   # Solo eventos
make db-reset         # Reset completo de BD

# Testing y calidad
make lint             # Linter de código
make format           # Formatear código
make security-check   # Análisis de seguridad

# Utilidades
make clean            # Limpiar archivos generados
make help             # Ver todos los comandos
make version          # Mostrar versión actual
```

## Arquitectura del Sistema

### Clean Architecture

El proyecto está estructurado siguiendo los principios de Clean Architecture. **Estado actual:**

```
Handlers (HTTP) → [Services] → [Repository] → Models (Domain) + Database
```

**Implementado en v0.1.0:**

- **Handlers**: Controladores HTTP completos con validación
- **Models**: Entidades de dominio con GORM y relaciones
- **Database**: Acceso directo desde handlers (temporal)

**Próximo en v0.2.0:**

- TBD

> **Nota**: En v0.1.0 los handlers acceden directamente a la base de datos via GORM para acelerar el desarrollo inicial. La refactorización hacia Clean Architecture completa está planificada para v0.2.0.

### Estructura del Proyecto

```
backend/
├── cmd/                    # Comandos ejecutables
│   ├── cybesphere/        # Aplicación principal
|   ├── seed/              # Herramienta de seeder
│   └── debug-config/      # Herramientas debug
|
├── internal/              # Código de aplicación
│   ├── config/           # Configuración y validación
│   ├── handlers/         # HTTP handlers y middleware
│   ├── services/         # Lógica de negocio y validaciones
│   ├── repository/       # Acceso a datos con GORM
│   ├── models/           # Modelos de dominio y DTOs
│   ├── middleware/       # Middleware HTTP personalizado
│   └── utils/            # Utilidades y helpers
├── pkg/                  # Paquetes reutilizables
│   ├── auth/            # JWT y manejo de tokens
│   ├── database/        # Conexión y configuración BD
│   ├── logger/          # Sistema de logging
│   └── validator/       # Validación de datos
├── migrations/          # Scripts SQL de migración
├── seeders/            # Datos de ejemplo para desarrollo
├── tests/             # Tests unitarios e integración
├── scripts/           # Scripts de deployment y CI/CD
├── docs/             # Documentación técnica
│   └── api/         # Documentación de la API
└── README.md         # Este archivo
```

## Stack Tecnológico

### Core

- **Runtime**: Go 1.24
- **Web Framework**: Gin v1.10.0
- **ORM**: GORM v1.25.12
- **Base de Datos**: PostgreSQL 15

### Configuración y Logging

- **Configuración**: Viper v1.19.0
- **Logging**: Logrus v1.9.3
- **Validación**: validator v10.22.1

### Autenticación y Seguridad

- **JWT**: golang-jwt/jwt/v5 v5.2.1
- **Password Hashing**: bcrypt
- **CORS**: gin-contrib/cors
- **Rate Limiting**: gin-contrib/limiter

### Desarrollo y Testing

- **Hot Reload**: Air v1.52.3
- **Testing**: testify v1.9.0
- **Mocking**: gomock y sqlmock
- **Containerización**: Docker Compose

## Características de Seguridad

### Autenticación JWT

- **Tokens de acceso** con expiración corta (15 min)
- **Refresh tokens** con rotación automática (7 días)
- **Blacklist de tokens** para logout seguro
- **Validación de claims** personalizados

### RBAC (Role-Based Access Control)

- **Admin**: Acceso completo al sistema
- **Organizer**: Gestión de eventos y organizaciones
- **Attendee**: Participación en eventos

### Middleware de Seguridad

- **CORS configurado** para dominios específicos
- **Rate limiting** por IP y usuario
- **Request ID** para trazabilidad
- **Recovery middleware** para errores panic
- **Logging de todas las requests**

## Configuración de Base de Datos

PostgreSQL configurado con:

- **Extensiones**: UUID, full-text search en español
- **Funciones personalizadas**: Geolocalización (Haversine)
- **Índices optimizados** para queries frecuentes
- **Timezone**: Europe/Madrid por defecto
- **Pool de conexiones** optimizado para concurrencia

## Testing

### Cobertura de Tests

- **Handlers HTTP**: Tests de endpoints con mocks
- **Services**: Lógica de negocio con casos edge
- **Repository**: Acceso a datos con BD en memoria
- **Middleware**: Autenticación y autorización
- **Integration**: Flujos completos end-to-end

### Ejecutar Tests

```bash
# Tests unitarios rápidos
make test

# Tests de integración (requiere BD)
make test-integration

# Todos los tests con coverage
make test-all

# Coverage report en HTML
make test-coverage
```

## Troubleshooting

### Error de conexión PostgreSQL

```bash
# Verificar servicios
make docker-status

# Reiniciar servicios
make docker-restart

# Verificar logs
make docker-logs
```

### Problemas con migraciones

```bash
# Reset completo de BD
make db-reset

# Verificar estado de BD
go run ./cmd/debug-config
```

### Air no funciona

```bash
# Reinstalar Air
go install github.com/air-verse/air@latest

# Verificar PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Tests fallan

```bash
# Limpiar cache de Go
go clean -testcache

# Verificar variables de entorno de test
cat .env.test
```

## Deployment

### Variables de Entorno Producción

```bash
# Configuración obligatoria para producción
GIN_MODE=release
LOG_LEVEL=info
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
JWT_SECRET=<secret-64-chars>
JWT_REFRESH_SECRET=<different-secret-64-chars>
CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

### Build para Producción

```bash
# Compilar binario optimizado
make build-prod

# Ejecutar con configuración de producción
./bin/cybesphere
```

## Contribución

### Desarrollo

1. Fork del repositorio
2. Crear branch feature: `git checkout -b feature/nueva-funcionalidad`
3. Ejecutar tests: `make test-all`
4. Commit con formato convencional: `feat: nueva funcionalidad`
5. Push y crear Pull Request

### Estándares de Código

- **Formatting**: `make format` antes de commit
- **Linting**: `make lint` debe pasar sin errores
- **Testing**: Coverage mínimo 80%
- **Documentación**: Comentarios en funciones públicas

## Changelog

### v0.1.0 (Actual)

- Sistema completo de autenticación y autorización
- CRUD completo para usuarios, organizaciones y eventos
- Testing comprehensivo con 85%+ coverage
- Seeders con datos realistas de ciberseguridad
- Middleware de seguridad implementado
- Clean Architecture establecida

### v0.0.1 (Inicial)

- Setup inicial del proyecto
- Configuración Docker y base de datos
- Endpoints básicos de health check

## Licencia

TBD

---

**CybESphere Backend v0.1.0** - Sistema robusto para gestión de eventos de ciberseguridad con arquitectura escalable y seguridad enterprise.

Para información detallada sobre el uso de la API, consulta la [documentación completa](./docs/api/README.md).
