# API de Organizaciones - CybESphere

Este documento describe todos los endpoints relacionados con la gestión de organizaciones.

## Base URL

```
/api/v1
```

## Endpoints Públicos

### 1. Listar Organizaciones Públicas

**GET** `/public/organizations`

Obtiene una lista paginada de organizaciones activas y verificadas.

#### Query Parameters

```
?page=1                    // Página (por defecto: 1)
&limit=20                  // Elementos por página (por defecto: 20, máx: 100)
&search=cybesphere         // Búsqueda en nombre y descripción
&city=Madrid               // Filtrar por ciudad
&country=España            // Filtrar por país
&is_verified=true          // Solo organizaciones verificadas
&order_by=name             // Ordenar por: name, created_at, events_count, city
&order_dir=asc             // Dirección: asc, desc
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Organizaciones obtenidas exitosamente",
  "data": {
    "organizations": [
      {
        "id": "456e7890-e12b-34d5-b678-901234567890",
        "slug": "cybesphere-organization",
        "name": "CybESphere Organization",
        "description": "Organización líder en eventos de ciberseguridad y tecnología",
        "website": "https://cybesphere.org",
        "city": "Madrid",
        "country": "España",
        "logo_url": "https://ejemplo.com/logo.png",
        "banner_url": "https://ejemplo.com/banner.jpg",
        "primary_color": "#1a73e8",
        "secondary_color": "#34a853",
        "social_media": {
          "linkedin": "https://linkedin.com/company/cybesphere",
          "twitter": "https://twitter.com/cybesphere",
          "facebook": "https://facebook.com/cybesphere",
          "instagram": "https://instagram.com/cybesphere",
          "youtube": "https://youtube.com/c/cybesphere"
        },
        "status": "active",
        "is_verified": true,
        "verified_at": "2024-01-10T10:00:00Z",
        "events_count": 25,
        "members_count": 8,
        "upcoming_events": 5,
        "is_member": false,
        "can_edit": false,
        "can_manage": false,
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-15T14:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 45,
      "total_pages": 3,
      "has_prev": false,
      "has_next": true,
      "prev_page": null,
      "next_page": 2
    },
    "filters": {
      "search": "cybesphere",
      "order_by": "name",
      "order_dir": "asc",
      "filters": {
        "city": "Madrid",
        "is_verified": true
      }
    }
  }
}
```

---

### 2. Obtener Organización por ID

**GET** `/public/organizations/{id}`

Obtiene los detalles completos de una organización específica.

#### Path Parameters

- `id`: UUID de la organización

#### Response Success (200)

```json
{
  "success": true,
  "message": "Organización obtenida exitosamente",
  "data": {
    "id": "456e7890-e12b-34d5-b678-901234567890",
    "slug": "cybesphere-organization",
    "name": "CybESphere Organization",
    "description": "Organización líder en eventos de ciberseguridad y tecnología...",
    "website": "https://cybesphere.org",
    "city": "Madrid",
    "country": "España",
    "latitude": 40.4168,
    "longitude": -3.7038,
    "logo_url": "https://ejemplo.com/logo.png",
    "banner_url": "https://ejemplo.com/banner.jpg",
    "primary_color": "#1a73e8",
    "secondary_color": "#34a853",
    "social_media": {
      "linkedin": "https://linkedin.com/company/cybesphere",
      "twitter": "https://twitter.com/cybesphere",
      "facebook": "https://facebook.com/cybesphere",
      "instagram": "https://instagram.com/cybesphere",
      "youtube": "https://youtube.com/c/cybesphere"
    },
    "status": "active",
    "is_verified": true,
    "verified_at": "2024-01-10T10:00:00Z",
    "events_count": 25,
    "members_count": 8,
    "upcoming_events": 5,
    "is_member": false,
    "can_edit": false,
    "can_manage": false,
    "recent_events": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "slug": "evento-reciente",
        "title": "Evento Reciente",
        "short_desc": "Descripción del evento reciente",
        "type": "conference",
        "start_date": "2024-03-15T09:00:00Z",
        "end_date": "2024-03-16T18:00:00Z",
        "is_online": false,
        "venue_city": "Madrid",
        "is_free": false,
        "price": 15000,
        "image_url": "https://ejemplo.com/evento.jpg",
        "current_attendees": 287,
        "max_attendees": 500,
        "is_featured": true,
        "organization_name": "CybESphere Organization",
        "tags": ["ciberseguridad", "tecnología"]
      }
    ],
    "featured_members": [
      {
        "id": "789e0123-e45f-67g8-h901-234567890123",
        "email": "miembro@cybesphere.org",
        "first_name": "Ana",
        "last_name": "García",
        "full_name": "Ana García",
        "company": "CybESphere Organization",
        "position": "CTO",
        "role": "organizer",
        "is_verified": true,
        "organization_name": "CybESphere Organization"
      }
    ],
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-15T14:30:00Z"
  }
}
```

---

### 3. Organizaciones Activas

**GET** `/public/organizations/active`

Obtiene solo las organizaciones activas y verificadas.

#### Query Parameters

```
?limit=50              // Número de organizaciones (por defecto: 50)
&city=Madrid          // Filtrar por ciudad
&country=España       // Filtrar por país
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Organizaciones activas obtenidas",
  "data": [
    {
      "id": "456e7890-e12b-34d5-b678-901234567890",
      "slug": "cybesphere-organization",
      "name": "CybESphere Organization",
      "logo_url": "https://ejemplo.com/logo.png",
      "is_verified": true,
      "events_count": 25,
      "city": "Madrid",
      "country": "España"
    }
  ]
}
```

---

## Endpoints Protegidos

### 4. Crear Organización

**POST** `/organizations`

Crea una nueva organización. Requiere usuario verificado.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
Content-Type: application/json
```

#### Request Body

```json
{
  "name": "Mi Organización",
  "description": "Descripción detallada de la organización con al menos 10 caracteres",
  "website": "https://mi-organizacion.com",
  "email": "contacto@mi-organizacion.com",
  "phone": "+34 600 123 456",
  "address": "Calle Ejemplo, 123",
  "city": "Madrid",
  "country": "España",
  "postal_code": "28001",
  "latitude": 40.4168,
  "longitude": -3.7038,
  "logo_url": "https://mi-organizacion.com/logo.png",
  "banner_url": "https://mi-organizacion.com/banner.jpg",
  "primary_color": "#1a73e8",
  "secondary_color": "#34a853",
  "linkedin": "https://linkedin.com/company/mi-org",
  "twitter": "https://twitter.com/mi_org",
  "facebook": "https://facebook.com/mi-org",
  "instagram": "https://instagram.com/mi_org",
  "youtube": "https://youtube.com/c/mi-org",
  "tax_id": "B12345678",
  "legal_name": "Mi Organización S.L.",
  "registration_docs": "https://docs.com/registro.pdf"
}
```

#### Campos Opcionales para Admin

```json
{
  "status": "active", // pending, active, suspended, inactive
  "is_verified": true, // Verificación manual
  "max_events": 100, // Límite de eventos
  "can_create_events": true // Permiso para crear eventos
}
```

#### Response Success (201)

```json
{
  "success": true,
  "message": "Organización creada exitosamente",
  "data": {
    // OrganizationResponse completo de la organización creada
    "id": "456e7890-e12b-34d5-b678-901234567890",
    "slug": "mi-organizacion",
    "name": "Mi Organización",
    "status": "pending", // Requiere verificación por defecto
    "is_verified": false
    // ... resto de campos
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
    "name": ["El nombre es requerido"],
    "email": ["El email ya está en uso por otra organización"],
    "primary_color": ["El color debe estar en formato hexadecimal"]
  }
}
```

---

### 5. Actualizar Organización

**PUT** `/organizations/{id}`

Actualiza una organización existente. Requiere ser miembro o admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

Cualquier campo del CreateOrganizationRequest puede ser actualizado:

```json
{
  "name": "Nombre Actualizado",
  "description": "Nueva descripción",
  "website": "https://nuevo-sitio.com",
  "primary_color": "#ff5722"
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Organización actualizada exitosamente",
  "data": {
    // OrganizationResponse completo actualizado
  }
}
```

---

### 6. Eliminar Organización

**DELETE** `/organizations/{id}`

Elimina una organización. Solo admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Organización eliminada exitosamente",
  "data": {
    "id": "456e7890-e12b-34d5-b678-901234567890",
    "success": true,
    "message": "Organización eliminada exitosamente"
  }
}
```

---

### 7. Verificar Organización

**POST** `/organizations/{id}/verify`

Verifica una organización. Solo admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "notes": "Documentación revisada y verificada"
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Organización verificada exitosamente",
  "data": {
    "id": "456e7890-e12b-34d5-b678-901234567890",
    "name": "Mi Organización",
    "is_verified": true,
    "verified_at": "2024-01-15T10:00:00Z",
    "verified_by": "admin@cybesphere.com",
    "status": "active"
  }
}
```

---

### 8. Obtener Miembros de Organización

**GET** `/organizations/{id}/members`

Obtiene la lista de miembros de una organización. Solo miembros de la organización o admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Query Parameters

```
?page=1              // Página
&limit=20            // Elementos por página
&search=juan         // Búsqueda por nombre o email
&role=organizer      // Filtrar por rol
&is_active=true      // Solo miembros activos
&order_by=joined_at  // Ordenar por: joined_at, name, email
&order_dir=desc      // Dirección
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Miembros de la organización obtenidos",
  "data": {
    "organization_id": "456e7890-e12b-34d5-b678-901234567890",
    "organization_name": "CybESphere Organization",
    "members": [
      {
        "id": "789e0123-e45f-67g8-h901-234567890123",
        "email": "miembro@cybesphere.org",
        "first_name": "Ana",
        "last_name": "García",
        "role": "organizer",
        "position": "CTO",
        "joined_at": "2024-01-01T10:00:00Z",
        "is_active": true
      },
      {
        "id": "987f6543-e21d-43c5-b876-543210987654",
        "email": "otro@cybesphere.org",
        "first_name": "Carlos",
        "last_name": "López",
        "role": "user",
        "position": "Desarrollador",
        "joined_at": "2024-01-05T14:00:00Z",
        "is_active": true
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 8,
      "total_pages": 1,
      "has_prev": false,
      "has_next": false
    }
  }
}
```

---

## Endpoints de Administración

### 9. Verificación Masiva de Organizaciones

**POST** `/admin/organizations/bulk-verify`

Verifica múltiples organizaciones a la vez. Solo admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "organization_ids": [
    "456e7890-e12b-34d5-b678-901234567890",
    "789e0123-e45f-67g8-h901-234567890123"
  ]
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Organizaciones verificadas exitosamente",
  "data": {
    "requested": 2,
    "verified": 2,
    "timestamp": "2024-01-15T10:00:00Z"
  }
}
```

---

## Estados de Organizaciones

- `pending`: Pendiente de verificación
- `active`: Activa y verificada
- `suspended`: Suspendida temporalmente
- `inactive`: Inactiva

## Códigos de Error Específicos

### 400 - Bad Request

- `validation_error`: Error en validación de campos
- `invalid_color_format`: Color debe estar en formato hexadecimal (#RRGGBB)
- `invalid_url_format`: URL no válida
- `duplicate_email`: El email ya está en uso por otra organización
- `duplicate_name`: El nombre ya está en uso
- `invalid_coordinates`: Coordenadas geográficas inválidas
- `invalid_phone_format`: Formato de teléfono inválido

### 403 - Forbidden

- `insufficient_permissions`: No tienes permisos para esta acción
- `not_organization_member`: No eres miembro de esta organización
- `organization_suspended`: La organización está suspendida
- `verification_required`: La organización requiere verificación

### 404 - Not Found

- `organization_not_found`: Organización no encontrada

### 409 - Conflict

- `organization_slug_exists`: Ya existe una organización con ese slug
- `email_already_verified`: La organización ya está verificada
- `organization_has_events`: No se puede eliminar una organización con eventos activos

---

## Permisos y Roles

### Permisos de Organización

- **Ver Pública**: Cualquier usuario
- **Ver Privada**: Miembros de la organización o admin
- **Crear**: Usuario verificado
- **Actualizar**: Miembro de la organización (rol organizer) o admin
- **Eliminar**: Solo admin
- **Verificar**: Solo admin
- **Ver Miembros**: Miembros de la organización o admin

### Roles en Organización

- `user`: Usuario básico de la organización
- `organizer`: Puede gestionar eventos y la organización
- `admin`: Administrador del sistema (todos los permisos)

## Notas Importantes

1. **Verificación**: Las organizaciones nuevas requieren verificación manual por admin
2. **Slug**: Se genera automáticamente desde el nombre
3. **Colores**: Deben estar en formato hexadecimal (#RRGGBB)
4. **Coordenadas**: Latitud entre -90 y 90, Longitud entre -180 y 180
5. **URLs**: Todas las URLs deben ser válidas y preferiblemente HTTPS
6. **Límites**: Una organización puede tener un límite de eventos (configurable por admin)
7. **Suspensión**: Organizaciones suspendidas no pueden crear nuevos eventos
8. **Eliminación**: Solo se pueden eliminar organizaciones sin eventos activos
9. **Miembros**: Los usuarios se asocian a organizaciones mediante el campo `organization_id`
10. **Documentación**: Los documentos de registro son URLs a archivos externos

## Campos Sensibles

Algunos campos solo son visibles para:

- **Información de contacto** (email, phone, address): Miembros y admin
- **Documentación legal** (tax_id, legal_name, registration_docs): Solo admin
- **Estadísticas detalladas**: Miembros y admin
- **Configuración administrativa** (max_events, can_create_events): Solo admin
