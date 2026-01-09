# API de Autenticación - CybESphere

Este documento describe todos los endpoints relacionados con autenticación y gestión de sesiones.

## Base URL

```
/api/v1/auth
```

## Endpoints Públicos

### 1. Registro de Usuario

**POST** `/auth/register`

Registra un nuevo usuario en el sistema.

#### Request Body

```json
{
  "email": "usuario@ejemplo.com",
  "password": "password123",
  "first_name": "Juan",
  "last_name": "Pérez",
  "company": "Empresa SA", // Opcional
  "position": "Desarrollador", // Opcional
  "city": "Madrid", // Opcional
  "country": "España" // Opcional
}
```

#### Validaciones

- `email`: Requerido, formato email válido
- `password`: Requerido, mínimo 8 caracteres, máximo 128
- `first_name`: Requerido, mínimo 1 carácter, máximo 100
- `last_name`: Requerido, mínimo 1 carácter, máximo 100
- `company`: Opcional, máximo 200 caracteres
- `position`: Opcional, máximo 200 caracteres
- `city`: Opcional, máximo 100 caracteres
- `country`: Opcional, máximo 100 caracteres

#### Response Success (201)

```json
{
  "success": true,
  "message": "Usuario registrado exitosamente. Se ha enviado un email de verificación.",
  "data": {
    "message": "Usuario registrado exitosamente. Se ha enviado un email de verificación.",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "usuario@ejemplo.com",
      "first_name": "Juan",
      "last_name": "Pérez",
      "full_name": "Juan Pérez",
      "role": "user",
      "company": "Empresa SA",
      "position": "Desarrollador",
      "city": "Madrid",
      "country": "España",
      "is_active": true,
      "is_verified": false,
      "timezone": "Europe/Madrid",
      "language": "es",
      "created_at": "2024-01-01T10:00:00Z",
      "last_login_at": null,
      "organization": null
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "def502001234567890abcdef...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "requires_verification": true
  }
}
```

#### Response Error (400)

```json
{
  "success": false,
  "error": "validation_error",
  "message": "Error de validación",
  "errors": {
    "email": ["El email ya está en uso"],
    "password": ["La contraseña debe tener al menos 8 caracteres"]
  }
}
```

---

### 2. Inicio de Sesión

**POST** `/auth/login`

Autentica un usuario existente.

#### Request Body

```json
{
  "email": "usuario@ejemplo.com",
  "password": "password123"
}
```

#### Validaciones

- `email`: Requerido, formato email válido
- `password`: Requerido

#### Response Success (200)

```json
{
  "success": true,
  "message": "Inicio de sesión exitoso",
  "data": {
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "usuario@ejemplo.com",
      "first_name": "Juan",
      "last_name": "Pérez",
      "full_name": "Juan Pérez",
      "role": "user",
      "company": "Empresa SA",
      "position": "Desarrollador",
      "city": "Madrid",
      "country": "España",
      "is_active": true,
      "is_verified": true,
      "timezone": "Europe/Madrid",
      "language": "es",
      "created_at": "2024-01-01T10:00:00Z",
      "last_login_at": "2024-01-15T09:30:00Z",
      "organization": {
        "id": "456e7890-e12b-34d5-b678-901234567890",
        "slug": "empresa-sa",
        "name": "Empresa SA",
        "logo_url": "https://ejemplo.com/logo.png",
        "is_verified": true,
        "events_count": 15,
        "city": "Madrid",
        "country": "España"
      }
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "def502001234567890abcdef...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

#### Response Error (401)

```json
{
  "success": false,
  "error": "invalid_credentials",
  "message": "Email o contraseña incorrectos"
}
```

---

### 3. Renovar Token

**POST** `/auth/refresh`

Renueva un access token usando el refresh token.

#### Request Body

```json
{
  "refresh_token": "def502001234567890abcdef..."
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Token renovado exitosamente",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "def502001234567890abcdef...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "expires_at": "2024-01-15T11:30:00Z"
  }
}
```

#### Response Error (401)

```json
{
  "success": false,
  "error": "invalid_token",
  "message": "Refresh token inválido o expirado"
}
```

---

### 4. Cerrar Sesión

**POST** `/auth/logout`

Cierra la sesión actual revocando el refresh token.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "refresh_token": "def502001234567890abcdef..."
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Sesión cerrada exitosamente",
  "data": {
    "message": "Sesión cerrada exitosamente",
    "success": true
  }
}
```

---

## Endpoints Protegidos

### 5. Información del Usuario Actual

**GET** `/auth/me`

Obtiene la información del usuario autenticado actualmente.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Información del usuario obtenida",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "usuario@ejemplo.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "full_name": "Juan Pérez",
    "role": "user",
    "company": "Empresa SA",
    "position": "Desarrollador",
    "city": "Madrid",
    "country": "España",
    "is_active": true,
    "is_verified": true,
    "timezone": "Europe/Madrid",
    "language": "es",
    "newsletter_enabled": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-15T09:30:00Z",
    "last_login_at": "2024-01-15T09:30:00Z",
    "organization": {
      "id": "456e7890-e12b-34d5-b678-901234567890",
      "slug": "empresa-sa",
      "name": "Empresa SA",
      "logo_url": "https://ejemplo.com/logo.png",
      "is_verified": true,
      "events_count": 15,
      "city": "Madrid",
      "country": "España"
    },
    "can_edit": true,
    "can_manage": false
  }
}
```

---

### 6. Cerrar Todas las Sesiones

**POST** `/auth/logout-all`

Cierra todas las sesiones del usuario revocando todos sus refresh tokens.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "confirm_logout_all": true
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Todas las sesiones han sido cerradas",
  "data": {
    "message": "Todas las sesiones han sido cerradas",
    "success": true,
    "sessions_revoked": 3
  }
}
```

---

## Gestión de Sesiones

### 7. Listar Sesiones Activas

**GET** `/user/sessions`

Obtiene todas las sesiones activas del usuario.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Sesiones obtenidas exitosamente",
  "data": {
    "sessions": [
      {
        "id": "sess_123456789",
        "token_id": "token_abcdef123",
        "device_info": "Windows - Google Chrome",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36...",
        "created_at": "2024-01-15T09:00:00Z",
        "last_used_at": "2024-01-15T10:30:00Z",
        "expires_at": "2024-02-15T09:00:00Z",
        "is_current": true
      },
      {
        "id": "sess_987654321",
        "token_id": "token_fedcba987",
        "device_info": "iPhone",
        "ip_address": "192.168.1.101",
        "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X)...",
        "created_at": "2024-01-14T15:00:00Z",
        "last_used_at": "2024-01-15T08:00:00Z",
        "expires_at": "2024-02-14T15:00:00Z",
        "is_current": false
      }
    ],
    "total_sessions": 2,
    "active_sessions": 2
  }
}
```

---

### 8. Revocar Sesión Específica

**DELETE** `/user/sessions/{sessionId}`

Revoca una sesión específica.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Sesión revocada exitosamente",
  "data": {
    "session_id": "sess_987654321",
    "revoked_at": "2024-01-15T10:45:00Z"
  }
}
```

---

## Códigos de Error Comunes

### 400 - Bad Request

- `validation_error`: Error en la validación de datos
- `missing_field`: Campo requerido faltante
- `invalid_format`: Formato de datos inválido

### 401 - Unauthorized

- `invalid_credentials`: Email o contraseña incorrectos
- `invalid_token`: Token de acceso inválido o expirado
- `token_expired`: Token expirado
- `account_disabled`: Cuenta desactivada

### 403 - Forbidden

- `account_not_verified`: Cuenta no verificada
- `insufficient_permissions`: Permisos insuficientes

### 409 - Conflict

- `email_already_exists`: El email ya está registrado

### 429 - Too Many Requests

- `rate_limit_exceeded`: Límite de solicitudes excedido

### 500 - Internal Server Error

- `internal_error`: Error interno del servidor
- `database_error`: Error de base de datos

---

## Notas Importantes

1. **Tokens de Acceso**: Expiran en 1 hora (3600 segundos)
2. **Refresh Tokens**: Expiran en 30 días y se rotan en cada renovación
3. **Rate Limiting**: Máximo 5 intentos de login por minuto por IP
4. **Verificación de Email**: Los usuarios deben verificar su email para acceder a funcionalidades completas
5. **Timezone**: Por defecto se usa "Europe/Madrid"
6. **Idioma**: Por defecto se usa "es" (español)

## Headers de Autenticación

Para endpoints protegidos, incluir el header:

```
Authorization: Bearer {access_token}
```

## Formato de Respuesta Estándar

Todas las respuestas siguen este formato:

```json
{
  "success": boolean,
  "message": "string",
  "data": object | array,
  "error": "string",      // Solo en errores
  "errors": object        // Solo en errores de validación
}
```
