# API de Usuarios - CybESphere

Este documento describe todos los endpoints relacionados con la gestión de usuarios y perfiles.

## Base URL

```
/api/v1
```

## Endpoints Protegidos de Usuario

### 1. Obtener Perfil del Usuario Actual

**GET** `/user/profile`

Obtiene el perfil completo del usuario autenticado.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Perfil del usuario obtenido",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "usuario@ejemplo.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "full_name": "Juan Pérez",
    "company": "Empresa SA",
    "position": "Desarrollador Senior",
    "bio": "Desarrollador apasionado por la ciberseguridad y la tecnología",
    "website": "https://juan-perez.dev",
    "linkedin": "https://linkedin.com/in/juan-perez",
    "twitter": "https://twitter.com/juan_perez",
    "role": "user",
    "is_active": true,
    "is_verified": true,
    "city": "Madrid",
    "country": "España",
    "latitude": 40.4168,
    "longitude": -3.7038,
    "timezone": "Europe/Madrid",
    "language": "es",
    "newsletter_enabled": true,
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
    "statistics": {
      "events_attended": 12,
      "events_organized": 3,
      "favorite_events": 8,
      "total_connections": 45,
      "member_since": "2024-01-01T10:00:00Z",
      "last_active": "2024-01-15T09:30:00Z"
    },
    "favorite_events": [
      {
        "id": "789e0123-e45f-67g8-h901-234567890123",
        "slug": "evento-favorito",
        "title": "Evento Favorito",
        "short_desc": "Descripción del evento favorito",
        "type": "conference",
        "start_date": "2024-03-15T09:00:00Z",
        "end_date": "2024-03-16T18:00:00Z",
        "is_online": false,
        "venue_city": "Barcelona",
        "is_free": false,
        "price": 12000,
        "image_url": "https://ejemplo.com/evento.jpg",
        "current_attendees": 156,
        "max_attendees": 300,
        "is_featured": false,
        "organization_name": "Otra Organización",
        "tags": ["desarrollo", "AI"]
      }
    ],
    "registered_events": [
      // Array similar a favorite_events para eventos registrados
    ],
    "active_sessions": [
      {
        "id": "sess_123456789",
        "token_id": "token_abcdef123",
        "device_info": "Windows - Google Chrome",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
        "created_at": "2024-01-15T09:00:00Z",
        "last_used_at": "2024-01-15T10:30:00Z",
        "expires_at": "2024-02-15T09:00:00Z",
        "is_current": true
      }
    ],
    "can_edit": true,
    "can_manage": false,
    "last_login_at": "2024-01-15T09:30:00Z",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-15T09:30:00Z"
  }
}
```

---

### 2. Obtener Capacidades del Usuario

**GET** `/user/capabilities`

Obtiene los permisos y capacidades del usuario actual.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Capacidades del usuario obtenidas",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "role": "organizer",
    "permissions": [
      {
        "resource": "events",
        "action": "create",
        "display": "Crear eventos"
      },
      {
        "resource": "events",
        "action": "read",
        "display": "Ver eventos"
      },
      {
        "resource": "events",
        "action": "update",
        "display": "Actualizar eventos"
      },
      {
        "resource": "organizations",
        "action": "read",
        "display": "Ver organizaciones"
      },
      {
        "resource": "organizations",
        "action": "update",
        "display": "Actualizar organización"
      }
    ],
    "capabilities": {
      "can_create_events": true,
      "can_manage_organization": true,
      "can_view_analytics": true,
      "can_moderate_content": false,
      "can_manage_users": false,
      "can_access_admin_panel": false,
      "can_verify_organizations": false,
      "can_feature_events": false
    },
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
  }
}
```

---

### 3. Verificar Acceso a Recurso

**POST** `/user/check-access`

Verifica si el usuario tiene acceso a un recurso específico.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "resource": "events",
  "action": "create",
  "resource_id": "123e4567-e89b-12d3-a456-426614174000" // Opcional
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Acceso verificado",
  "data": {
    "has_access": true,
    "resource": "events",
    "action": "create",
    "reason": "Usuario tiene permisos de organización"
  }
}
```

#### Response Denied (200)

```json
{
  "success": true,
  "message": "Acceso verificado",
  "data": {
    "has_access": false,
    "resource": "events",
    "action": "delete",
    "reason": "Solo el propietario o admin puede eliminar eventos"
  }
}
```

---

### 4. Obtener Acciones Disponibles

**GET** `/user/available-actions`

Obtiene las acciones que el usuario puede realizar según su rol.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Query Parameters

```
?resource=events     // Filtrar por recurso específico
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Acciones disponibles obtenidas",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "role": "organizer",
    "actions": {
      "events": ["create", "read", "update", "delete", "publish"],
      "organizations": ["read", "update"],
      "users": ["read_own_profile", "update_own_profile"]
    }
  }
}
```

---

### 5. Obtener Sesiones del Usuario

**GET** `/user/sessions`

Lista todas las sesiones activas del usuario actual.

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
      }
    ],
    "total_sessions": 3,
    "active_sessions": 3
  }
}
```

---

### 6. Revocar Sesión

**DELETE** `/user/sessions/{sessionId}`

Revoca una sesión específica del usuario.

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

### 7. Obtener Información de Roles

**GET** `/user/roles`

Obtiene información detallada sobre los roles del sistema.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Información de roles obtenida",
  "data": {
    "current_role": "organizer",
    "available_roles": [
      {
        "role": "user",
        "display": "Usuario",
        "description": "Usuario básico con acceso limitado",
        "permissions": [
          "read_public_events",
          "register_to_events",
          "manage_own_profile"
        ]
      },
      {
        "role": "organizer",
        "display": "Organizador",
        "description": "Puede crear y gestionar eventos de su organización",
        "permissions": [
          "create_events",
          "manage_organization_events",
          "view_organization_analytics",
          "manage_own_profile"
        ]
      },
      {
        "role": "admin",
        "display": "Administrador",
        "description": "Acceso completo al sistema",
        "permissions": [
          "manage_all_content",
          "verify_organizations",
          "manage_users",
          "access_system_settings"
        ]
      }
    ]
  }
}
```

---

## Gestión de Usuarios (Admin)

### 8. Listar Todos los Usuarios

**GET** `/users`

Obtiene una lista paginada de todos los usuarios. Solo admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Query Parameters

```
?page=1                    // Página
&limit=20                  // Elementos por página
&search=juan               // Búsqueda por nombre o email
&role=organizer            // Filtrar por rol: user, organizer, admin
&is_active=true            // Solo usuarios activos
&is_verified=true          // Solo usuarios verificados
&organization_id=456...    // Filtrar por organización
&city=Madrid              // Filtrar por ciudad
&country=España           // Filtrar por país
&order_by=created_at      // Ordenar por: created_at, email, first_name, last_name, last_login_at
&order_dir=desc           // Dirección: asc, desc
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Usuarios obtenidos exitosamente",
  "data": {
    "users": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "email": "usuario@ejemplo.com",
        "first_name": "Juan",
        "last_name": "Pérez",
        "full_name": "Juan Pérez",
        "company": "Empresa SA",
        "position": "Desarrollador",
        "role": "organizer",
        "is_active": true,
        "is_verified": true,
        "city": "Madrid",
        "country": "España",
        "timezone": "Europe/Madrid",
        "language": "es",
        "newsletter_enabled": true,
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
        "can_manage": true,
        "last_login_at": "2024-01-15T09:30:00Z",
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-15T09:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 456,
      "total_pages": 23,
      "has_prev": false,
      "has_next": true
    },
    "filters": {
      "search": "juan",
      "order_by": "created_at",
      "order_dir": "desc",
      "filters": {
        "role": "organizer",
        "is_active": true
      }
    }
  }
}
```

---

### 9. Obtener Usuario por ID

**GET** `/users/{id}`

Obtiene los detalles de un usuario específico. Propio perfil o admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Usuario obtenido exitosamente",
  "data": {
    // UserDetailResponse completo (similar al endpoint de perfil)
  }
}
```

---

### 10. Actualizar Usuario

**PUT** `/users/{id}`

Actualiza un usuario. Propio perfil o admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "first_name": "Juan Carlos",
  "last_name": "Pérez García",
  "company": "Nueva Empresa",
  "position": "Senior Developer",
  "bio": "Nueva biografía actualizada",
  "website": "https://nuevo-sitio.com",
  "linkedin": "https://linkedin.com/in/juan-carlos",
  "twitter": "https://twitter.com/juan_carlos",
  "city": "Barcelona",
  "country": "España",
  "latitude": 41.3851,
  "longitude": 2.1734,
  "timezone": "Europe/Madrid",
  "language": "es",
  "newsletter_enabled": false
}
```

#### Campos Adicionales para Admin

```json
{
  "role": "organizer", // Cambiar rol
  "is_active": true, // Activar/desactivar
  "is_verified": true, // Verificar usuario
  "organization_id": "456e7890-e12b..." // Asignar organización
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Usuario actualizado exitosamente",
  "data": {
    // UserResponse completo actualizado
  }
}
```

---

### 11. Cambiar Rol de Usuario

**PUT** `/users/{id}/role`

Cambia el rol de un usuario. Solo admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "role": "organizer",
  "reason": "Usuario ha demostrado capacidad para organizar eventos"
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Rol de usuario actualizado",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "old_role": "user",
    "new_role": "organizer",
    "changed_by": "admin@cybesphere.com",
    "changed_at": "2024-01-15T10:00:00Z",
    "reason": "Usuario ha demostrado capacidad para organizar eventos"
  }
}
```

---

### 12. Activar Usuario

**POST** `/users/{id}/activate`

Activa un usuario desactivado. Solo admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "is_active": true,
  "reason": "Usuario ha cumplido con los requisitos"
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Usuario activado exitosamente",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "is_active": true,
    "changed_by": "admin@cybesphere.com",
    "changed_at": "2024-01-15T10:00:00Z",
    "reason": "Usuario ha cumplido con los requisitos"
  }
}
```

---

### 13. Desactivar Usuario

**POST** `/users/{id}/deactivate`

Desactiva un usuario. Solo admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "is_active": false,
  "reason": "Violación de términos de servicio"
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Usuario desactivado exitosamente",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "is_active": false,
    "changed_by": "admin@cybesphere.com",
    "changed_at": "2024-01-15T10:00:00Z",
    "reason": "Violación de términos de servicio"
  }
}
```

---

### 14. Obtener Sesiones de Usuario

**GET** `/users/{id}/sessions`

Obtiene las sesiones de un usuario específico. Propio perfil o admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Sesiones del usuario obtenidas",
  "data": {
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "sessions": [
      {
        "id": "sess_123456789",
        "token_id": "token_abcdef123",
        "device_info": "Windows - Google Chrome",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2024-01-15T09:00:00Z",
        "last_used_at": "2024-01-15T10:30:00Z",
        "expires_at": "2024-02-15T09:00:00Z",
        "is_current": false
      }
    ],
    "total_sessions": 2,
    "active_sessions": 2
  }
}
```

---

## Endpoints Públicos de Perfil

### 15. Ver Perfil Público

**GET** `/public/users/{id}/profile`

Obtiene el perfil público de un usuario.

#### Response Success (200)

```json
{
  "success": true,
  "message": "Perfil público obtenido",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "first_name": "Juan",
    "last_name": "Pérez",
    "full_name": "Juan Pérez",
    "company": "Empresa SA",
    "position": "Desarrollador Senior",
    "bio": "Desarrollador apasionado por la ciberseguridad",
    "city": "Madrid",
    "country": "España",
    "linkedin": "https://linkedin.com/in/juan-perez",
    "twitter": "https://twitter.com/juan_perez",
    "website": "https://juan-perez.dev",
    "public_events": [
      {
        "id": "789e0123-e45f-67g8-h901-234567890123",
        "slug": "evento-publico",
        "title": "Evento Público",
        "short_desc": "Evento organizado por Juan",
        "type": "workshop",
        "start_date": "2024-03-15T09:00:00Z",
        "end_date": "2024-03-15T17:00:00Z",
        "is_online": true,
        "venue_city": null,
        "is_free": true,
        "price": null,
        "image_url": "https://ejemplo.com/evento.jpg",
        "current_attendees": 45,
        "max_attendees": 100,
        "is_featured": false,
        "organization_name": "Empresa SA",
        "tags": ["workshop", "desarrollo"]
      }
    ]
  }
}
```

---

## Códigos de Error Específicos

### 400 - Bad Request

- `validation_error`: Error en validación de campos
- `invalid_email_format`: Formato de email inválido
- `invalid_coordinates`: Coordenadas geográficas inválidas
- `invalid_timezone`: Zona horaria no válida
- `invalid_language`: Idioma no soportado (debe ser 'es' o 'en')
- `invalid_role`: Rol no válido
- `password_too_weak`: Contraseña muy débil

### 403 - Forbidden

- `insufficient_permissions`: No tienes permisos para esta acción
- `cannot_modify_own_role`: No puedes cambiar tu propio rol
- `cannot_deactivate_self`: No puedes desactivarte a ti mismo
- `account_suspended`: Tu cuenta está suspendida
- `organization_required`: Debes pertenecer a una organización

### 404 - Not Found

- `user_not_found`: Usuario no encontrado
- `session_not_found`: Sesión no encontrada

### 409 - Conflict

- `email_already_exists`: El email ya está en uso
- `user_already_in_organization`: El usuario ya pertenece a una organización
- `cannot_delete_last_admin`: No se puede eliminar el último administrador

---

## Roles y Permisos

### Roles del Sistema

- `user`: Usuario básico
- `organizer`: Organizador de eventos
- `admin`: Administrador del sistema

### Permisos por Rol

#### Usuario (user)

- Ver eventos públicos
- Registrarse a eventos
- Gestionar su propio perfil
- Agregar eventos a favoritos

#### Organizador (organizer)

- Todos los permisos de usuario
- Crear y gestionar eventos de su organización
- Ver analytics de su organización
- Gestionar miembros de su organización

#### Administrador (admin)

- Todos los permisos del sistema
- Gestionar usuarios y roles
- Verificar organizaciones
- Moderar contenido
- Acceder a panel de administración

## Notas Importantes

1. **Privacidad**: Los usuarios solo pueden ver información completa de su propio perfil
2. **Administradores**: Pueden ver y editar cualquier perfil de usuario
3. **Verificación**: Los usuarios verificados tienen acceso a más funcionalidades
4. **Organizaciones**: Los usuarios pueden pertenecer a una sola organización
5. **Sesiones**: Los usuarios pueden tener múltiples sesiones activas
6. **Datos Sensibles**: Email y información de contacto solo es visible para el propietario y admin
7. **Geolocalización**: Las coordenadas son opcionales y privadas
8. **Configuraciones**: Timezone y language afectan la experiencia del usuario
9. **Estadísticas**: Solo son visibles para el propietario del perfil
10. **Perfil Público**: Los usuarios pueden tener un perfil público limitado
