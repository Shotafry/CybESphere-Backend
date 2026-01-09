# API CybESphere - Documentación General

## Descripción General

La API de CybESphere es una REST API completa para la gestión de eventos de ciberseguridad y tecnología. Permite a organizaciones crear y gestionar eventos, a usuarios registrarse y participar, y proporciona un sistema completo de autenticación y autorización.

## Información Básica

- **Versión**: v1
- **Base URL**: `https://api.cybesphere.com/api/v1`
- **Protocolo**: HTTPS
- **Formato**: JSON
- **Autenticación**: JWT Bearer Token

## Arquitectura

La API está organizada en los siguientes módulos principales:

### 🔐 Autenticación (`/auth`)

- Registro y login de usuarios
- Gestión de tokens JWT
- Renovación de tokens
- Gestión de sesiones

### 👥 Usuarios (`/users`)

- Gestión de perfiles de usuario
- Sistema de roles y permisos
- Capacidades y autorizaciones
- Sesiones activas

### 🏢 Organizaciones (`/organizations`)

- Creación y gestión de organizaciones
- Verificación de organizaciones
- Gestión de miembros
- Branding y configuración

### 📅 Eventos (`/events`)

- Creación y gestión de eventos
- Publicación y moderación
- Sistema de favoritos
- Registro y asistencia

## Endpoints Principales

### Documentación Específica

- [📖 Autenticación](auth_endpoints.md)
- [📖 Eventos](events_endpoints.md)
- [📖 Organizaciones](organizations_endpoints.md)
- [📖 Usuarios](users_endpoints.md)

### Endpoints Públicos (Sin Autenticación)

```
GET  /health                             # Health check
GET  /api/v1/public/ping                 # Ping test
GET  /api/v1/public/events               # Eventos públicos
GET  /api/v1/public/events/{id}          # Detalle de evento
GET  /api/v1/public/organizations        # Organizaciones públicas
GET  /api/v1/public/stats                # Estadísticas públicas
POST /api/v1/auth/register               # Registro
POST /api/v1/auth/login                  # Login
POST /api/v1/auth/refresh                # Renovar token
```

### Endpoints Protegidos (Requieren Autenticación)

```
GET  /api/v1/auth/me                     # Perfil actual
GET  /api/v1/user/capabilities           # Capacidades del usuario
GET  /api/v1/events                      # CRUD eventos
GET  /api/v1/organizations               # CRUD organizaciones
GET  /api/v1/users                       # Gestión usuarios (admin)
```

### Endpoints de Administración (Solo Admin)

```
GET  /api/v1/admin/dashboard             # Panel de admin
GET  /api/v1/admin/system/stats          # Estadísticas del sistema
POST /api/v1/organizations/{id}/verify   # Verificar organización
PUT  /api/v1/users/{id}/role             # Cambiar rol de usuario
```

## Autenticación

### Flujo de Autenticación

1. **Registro**: `POST /auth/register`
2. **Login**: `POST /auth/login` → Recibe `access_token` y `refresh_token`
3. **Uso**: Incluir header `Authorization: Bearer {access_token}`
4. **Renovación**: `POST /auth/refresh` cuando el token expire

### Tokens

- **Access Token**: Válido por 1 hora
- **Refresh Token**: Válido por 30 días, se rota en cada renovación
- **Algoritmo**: HS256

### Headers de Autenticación

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

## Sistema de Roles

### Roles Disponibles

- **`user`**: Usuario básico
- **`organizer`**: Organizador de eventos
- **`admin`**: Administrador del sistema

### Permisos por Rol

#### Usuario (user)

- Ver contenido público
- Registrarse a eventos
- Gestionar su perfil
- Agregar favoritos

#### Organizador (organizer)

- Todos los permisos de usuario
- Crear eventos de su organización
- Gestionar eventos de su organización
- Ver analytics de organización
- Gestionar miembros de organización

#### Administrador (admin)

- Todos los permisos
- Gestionar usuarios y roles
- Verificar organizaciones
- Moderar contenido
- Destacar eventos
- Acceso a panel de administración

## Formato de Respuesta

### Respuesta Exitosa

```json
{
  "success": true,
  "message": "Descripción de la operación",
  "data": {
    // Datos de respuesta
  }
}
```

### Respuesta con Paginación

```json
{
  "success": true,
  "message": "Datos obtenidos exitosamente",
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 156,
      "total_pages": 8,
      "has_prev": false,
      "has_next": true,
      "prev_page": null,
      "next_page": 2
    }
  }
}
```

### Respuesta de Error

```json
{
  "success": false,
  "error": "error_code",
  "message": "Descripción del error",
  "details": {
    // Información adicional del error
  }
}
```

### Errores de Validación

```json
{
  "success": false,
  "error": "validation_error",
  "message": "Error de validación",
  "errors": {
    "field_name": ["Error 1", "Error 2"],
    "other_field": ["Error en este campo"]
  }
}
```

## Códigos de Estado HTTP

### Exitosos (2xx)

- **200**: OK - Operación exitosa
- **201**: Created - Recurso creado exitosamente
- **204**: No Content - Operación exitosa sin contenido

### Errores del Cliente (4xx)

- **400**: Bad Request - Datos inválidos o faltantes
- **401**: Unauthorized - No autenticado o token inválido
- **403**: Forbidden - No tienes permisos para esta acción
- **404**: Not Found - Recurso no encontrado
- **409**: Conflict - Conflicto con el estado actual
- **422**: Unprocessable Entity - Error de validación
- **429**: Too Many Requests - Límite de requests excedido

### Errores del Servidor (5xx)

- **500**: Internal Server Error - Error interno del servidor
- **502**: Bad Gateway - Error de gateway
- **503**: Service Unavailable - Servicio no disponible

## Paginación

### Parámetros de Query

```
?page=1          # Página (por defecto: 1)
&limit=20        # Elementos por página (por defecto: 20, máx: 100)
&order_by=field  # Campo para ordenar
&order_dir=asc   # Dirección: asc, desc
```

### Metadata de Paginación

```json
{
  "page": 1,
  "limit": 20,
  "total": 156,
  "total_pages": 8,
  "has_prev": false,
  "has_next": true,
  "prev_page": null,
  "next_page": 2
}
```

## Filtros y Búsqueda

### Parámetros Comunes

```
?search=término              # Búsqueda de texto
&status=published            # Filtro por estado
&is_active=true              # Filtro booleano
&start_date_from=2024-01-01  # Filtro de fecha desde
&start_date_to=2024-12-31    # Filtro de fecha hasta
&city=Madrid                 # Filtro por ciudad
&country=España              # Filtro por país
&tags=tag1,tag2              # Filtro por tags (separados por coma)
```

## Rate Limiting

### Límites por Endpoint

- **Auth endpoints**: 5 requests/minuto por IP
- **Public endpoints**: 100 requests/minuto por IP
- **Protected endpoints**: 1000 requests/minuto por usuario
- **Admin endpoints**: 500 requests/minuto por usuario

### Headers de Rate Limit

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640995200
```

## Formatos de Datos

### Fechas

- **Formato**: ISO 8601 UTC
- **Ejemplo**: `2024-01-15T10:30:00Z`

### Coordenadas Geográficas

- **Latitud**: -90.0 a 90.0
- **Longitud**: -180.0 a 180.0
- **Ejemplo**: `{"latitude": 40.4168, "longitude": -3.7038}`

### Precios

- **Formato**: Entero en céntimos
- **Ejemplo**: `15000` = 150.00 EUR

### Colores

- **Formato**: Hexadecimal
- **Ejemplo**: `#1a73e8`

### URLs

- **Formato**: URL válida con protocolo
- **Ejemplo**: `https://ejemplo.com/imagen.jpg`

## Seguridad

### HTTPS

- Todas las comunicaciones deben usar HTTPS
- HTTP será redirigido a HTTPS en producción

### Autenticación

- JWT con expiración corta (1 hora)
- Refresh tokens con rotación
- Logout invalida tokens

### Autorización

- Verificación de permisos por endpoint
- Control de acceso basado en roles
- Validación de ownership para recursos

### Validación

- Validación estricta de entrada
- Sanitización de datos
- Protección contra XSS y SQL injection

## Monitoreo y Logs

### Health Check

```
GET /health
```

Respuesta:

```json
{
  "status": "ok",
  "version": "0.1.0",
  "timestamp": "2024-01-15T10:00:00Z",
  "database": {
    "connected": true,
    "stats": {
      "open_connections": 5,
      "max_connections": 100
    }
  }
}
```

### Logs de Auditoría

- Todas las operaciones críticas se registran
- Incluye información del usuario, IP, timestamp y acción
- Disponible para admin en `/admin/audit-logs`

## Ambientes

### Desarrollo

- **URL**: `http://localhost:8080/api/v1`
- **Base de datos**: Local PostgreSQL
- **Logs**: Nivel DEBUG
- **Rate limiting**: Deshabilitado

### Staging

- **URL**: `https://staging-api.cybesphere.com/api/v1`
- **Base de datos**: PostgreSQL en cloud
- **Logs**: Nivel INFO
- **Rate limiting**: Habilitado (límites relajados)

### Producción

- **URL**: `https://api.cybesphere.com/api/v1`
- **Base de datos**: PostgreSQL en cloud con réplicas
- **Logs**: Nivel WARN
- **Rate limiting**: Habilitado (límites estrictos)

## Herramientas Recomendadas

### Clientes HTTP

- **Postman**: Colección disponible
- **Insomnia**: Workspace exportable
- **cURL**: Ejemplos en documentación
- **HTTPie**: Sintaxis simplificada

### SDKs y Librerías

```javascript
// JavaScript/TypeScript
import { CybESphereAPI } from "@cybesphere/api-client";

const api = new CybESphereAPI({
  baseURL: "https://api.cybesphere.com/api/v1",
  apiKey: "your-access-token",
});
```

## Ejemplos de Uso

### Autenticación Básica

```bash
# Registro
curl -X POST https://api.cybesphere.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com",
    "password": "password123",
    "first_name": "Juan",
    "last_name": "Pérez"
  }'

# Login
curl -X POST https://api.cybesphere.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com",
    "password": "password123"
  }'
```

### Crear Evento

```bash
curl -X POST https://api.cybesphere.com/api/v1/events \
  -H "Authorization: Bearer your-access-token" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Workshop de Ciberseguridad",
    "description": "Workshop práctico sobre las últimas técnicas de ciberseguridad",
    "type": "workshop",
    "start_date": "2024-03-15T10:00:00Z",
    "end_date": "2024-03-15T17:00:00Z",
    "is_online": false,
    "venue_address": "Calle Ejemplo, 123, Madrid",
    "max_attendees": 50,
    "is_free": true
  }'
```

### Obtener Eventos Públicos

```bash
curl -X GET "https://api.cybesphere.com/api/v1/public/events?page=1&limit=10&city=Madrid&is_free=true"
```

## Mejores Prácticas

### Manejo de Errores

```javascript
async function apiRequest(url, options = {}) {
  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || `HTTP ${response.status}`);
    }

    return data;
  } catch (error) {
    console.error("API Error:", error);
    throw error;
  }
}
```

### Renovación Automática de Tokens

```javascript
class TokenManager {
  constructor() {
    this.accessToken = localStorage.getItem("access_token");
    this.refreshToken = localStorage.getItem("refresh_token");
  }

  async refreshAccessToken() {
    try {
      const response = await fetch("/api/v1/auth/refresh", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: this.refreshToken }),
      });

      const data = await response.json();
      if (data.success) {
        this.accessToken = data.data.access_token;
        this.refreshToken = data.data.refresh_token;
        localStorage.setItem("access_token", this.accessToken);
        localStorage.setItem("refresh_token", this.refreshToken);
        return true;
      }
    } catch (error) {
      console.error("Token refresh failed:", error);
      this.logout();
    }
    return false;
  }

  async makeAuthenticatedRequest(url, options = {}) {
    let response = await fetch(url, {
      ...options,
      headers: {
        Authorization: `Bearer ${this.accessToken}`,
        ...options.headers,
      },
    });

    if (response.status === 401) {
      const refreshed = await this.refreshAccessToken();
      if (refreshed) {
        response = await fetch(url, {
          ...options,
          headers: {
            Authorization: `Bearer ${this.accessToken}`,
            ...options.headers,
          },
        });
      }
    }

    return response;
  }

  logout() {
    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");
    this.accessToken = null;
    this.refreshToken = null;
    window.location.href = "/login";
  }
}
```

### Paginación Eficiente

```javascript
class EventsPaginator {
  constructor(apiClient) {
    this.apiClient = apiClient;
    this.events = [];
    this.currentPage = 1;
    this.hasMore = true;
  }

  async loadNextPage(filters = {}) {
    if (!this.hasMore) return [];

    const response = await this.apiClient.getEvents({
      ...filters,
      page: this.currentPage,
      limit: 20,
    });

    if (response.success) {
      const newEvents = response.data.events;
      this.events.push(...newEvents);
      this.currentPage++;
      this.hasMore = response.data.pagination.has_next;
      return newEvents;
    }

    return [];
  }

  reset() {
    this.events = [];
    this.currentPage = 1;
    this.hasMore = true;
  }
}
```

## Migración y Versionado

### Versionado de API

- La API usa versionado semántico
- Cambios breaking requieren nueva versión major
- Cambios compatibles incrementan versión minor
- Bug fixes incrementan versión patch

### Deprecación

- Features deprecados se marcan 6 meses antes de eliminación
- Headers `X-API-Deprecated` indican funcionalidad deprecada
- Documentación incluye alternativas recomendadas

### Migración entre Versiones

- Soporte simultáneo de 2 versiones major
- Guías de migración disponibles
- Endpoints de compatibilidad cuando sea posible

## Conclusión

La API de CybESphere está diseñada para ser robusta, escalable y fácil de usar. Proporciona todas las herramientas necesarias para construir aplicaciones completas de gestión de eventos de ciberseguridad.

Para comenzar:

1. Regístrate en la plataforma
2. Obtén tus tokens de acceso
3. Explora la documentación específica de cada módulo
4. Implementa tu primera integración
5. Únete a la comunidad para obtener soporte

¡Bienvenido al ecosistema CybESphere!
