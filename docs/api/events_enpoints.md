# API de Eventos - CybESphere

Este documento describe todos los endpoints relacionados con la gestión de eventos.

## Base URL

```
/api/v1
```

## Endpoints Públicos

### 1. Listar Eventos Públicos

**GET** `/public/events`

Obtiene una lista paginada de eventos públicos.

#### Query Parameters

```
?page=1                          // Página (por defecto: 1)
&limit=20                        // Elementos por página (por defecto: 20, máx: 100)
&search=cyberseguridad          // Búsqueda en título y descripción
&type=conference                // Tipo de evento
&category=security              // Categoría
&level=beginner                 // Nivel
&city=Madrid                    // Ciudad
&country=España                 // País
&is_online=true                 // Solo eventos online
&is_free=true                   // Solo eventos gratuitos
&min_price=0                    // Precio mínimo
&max_price=100                  // Precio máximo
&start_date_from=2024-01-01     // Fecha inicio desde
&start_date_to=2024-12-31       // Fecha inicio hasta
&is_featured=true               // Solo eventos destacados
&tags=AI,Machine Learning       // Tags (separados por coma)
&order_by=start_date            // Ordenar por: start_date, created_at, title, views_count
&order_dir=asc                  // Dirección: asc, desc
```

#### Valores válidos para `type`

- `conference`
- `workshop`
- `meetup`
- `webinar`
- `training`
- `competition`
- `other`

#### Valores válidos para `level`

- `beginner`
- `intermediate`
- `advanced`

#### Response Success (200)

```json
{
  "success": true,
  "message": "Eventos obtenidos exitosamente",
  "data": {
    "events": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "slug": "conferencia-cyberseguridad-madrid-2024",
        "title": "Conferencia de Ciberseguridad Madrid 2024",
        "description": "La conferencia líder en ciberseguridad de España...",
        "short_desc": "Conferencia líder en ciberseguridad con expertos internacionales",
        "type": "conference",
        "category": "Ciberseguridad",
        "level": "intermediate",
        "start_date": "2024-03-15T09:00:00Z",
        "end_date": "2024-03-16T18:00:00Z",
        "timezone": "Europe/Madrid",
        "duration": 1980,
        "is_online": false,
        "venue_name": "Palacio de Congresos",
        "venue_address": "Calle de la Castellana, 123",
        "venue_city": "Madrid",
        "venue_country": "España",
        "latitude": 40.4168,
        "longitude": -3.7038,
        "online_url": null,
        "streaming_url": null,
        "max_attendees": 500,
        "current_attendees": 287,
        "available_spots": 213,
        "is_free": false,
        "price": 15000,
        "currency": "EUR",
        "registration_url": "https://evento.com/registro",
        "status": "published",
        "is_public": true,
        "is_featured": true,
        "views_count": 1543,
        "image_url": "https://ejemplo.com/evento-imagen.jpg",
        "banner_url": "https://ejemplo.com/evento-banner.jpg",
        "tags": ["ciberseguridad", "hacking", "AI"],
        "registration_open": true,
        "registration_start_date": "2024-01-01T00:00:00Z",
        "registration_end_date": "2024-03-10T23:59:59Z",
        "meta_title": "Conferencia Ciberseguridad Madrid 2024 - CybESphere",
        "meta_description": "Únete a la conferencia más importante de ciberseguridad...",
        "organization": {
          "id": "456e7890-e12b-34d5-b678-901234567890",
          "slug": "cybesphere-org",
          "name": "CybESphere Organization",
          "logo_url": "https://ejemplo.com/org-logo.png",
          "is_verified": true,
          "events_count": 25,
          "city": "Madrid",
          "country": "España"
        },
        "is_upcoming": true,
        "is_past": false,
        "is_ongoing": false,
        "is_favorite": false,
        "is_registered": false,
        "can_edit": false,
        "can_manage": false,
        "created_at": "2024-01-01T10:00:00Z",
        "updated_at": "2024-01-15T14:30:00Z",
        "published_at": "2024-01-10T10:00:00Z",
        "canceled_at": null,
        "completed_at": null
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 156,
      "total_pages": 8,
      "has_prev": false,
      "has_next": true,
      "prev_page": null,
      "next_page": 2
    },
    "filters": {
      "search": "cyberseguridad",
      "order_by": "start_date",
      "order_dir": "asc",
      "filters": {
        "type": "conference",
        "city": "Madrid"
      }
    }
  }
}
```

---

### 2. Obtener Evento por ID

**GET** `/public/events/{id}`

Obtiene los detalles completos de un evento específico.

#### Path Parameters

- `id`: UUID del evento

#### Response Success (200)

```json
{
  "success": true,
  "message": "Evento obtenido exitosamente",
  "data": {
    // Incluye todos los campos del EventResponse anterior PLUS:
    "requirements": "Conocimientos básicos de programación",
    "agenda": "09:00 - Registro\n10:00 - Conferencia inaugural\n...",
    "contact_email": "contacto@evento.com",
    "contact_phone": "+34 600 123 456",
    "statistics": null // Solo visible para organizadores
  }
}
```

---

### 3. Eventos Destacados

**GET** `/public/events/featured`

Obtiene los eventos destacados por la plataforma.

#### Query Parameters

```
?limit=10      // Número de eventos (por defecto: 10, máx: 50)
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Eventos destacados obtenidos",
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "slug": "evento-destacado",
      "title": "Evento Destacado",
      "short_desc": "Descripción breve del evento",
      "type": "conference",
      "start_date": "2024-03-15T09:00:00Z",
      "end_date": "2024-03-16T18:00:00Z",
      "is_online": false,
      "venue_city": "Madrid",
      "is_free": false,
      "price": 15000,
      "image_url": "https://ejemplo.com/imagen.jpg",
      "current_attendees": 287,
      "max_attendees": 500,
      "is_featured": true,
      "organization_name": "CybESphere Organization",
      "tags": ["ciberseguridad", "AI"]
    }
  ]
}
```

---

### 4. Próximos Eventos

**GET** `/public/events/upcoming`

Obtiene eventos próximos organizados por períodos de tiempo.

#### Response Success (200)

```json
{
  "success": true,
  "message": "Próximos eventos obtenidos",
  "data": {
    "today": [
      // Array de EventSummaryResponse
    ],
    "this_week": [
      // Array de EventSummaryResponse
    ],
    "this_month": [
      // Array de EventSummaryResponse
    ],
    "later": [
      // Array de EventSummaryResponse
    ]
  }
}
```

---

## Endpoints Protegidos

### 5. Crear Evento

**POST** `/events`

Crea un nuevo evento. Requiere permisos de organización.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
Content-Type: application/json
```

#### Request Body

```json
{
  "title": "Mi Nuevo Evento",
  "description": "Descripción completa del evento con al menos 10 caracteres",
  "short_desc": "Descripción breve del evento",
  "type": "workshop",
  "category": "Desarrollo",
  "level": "beginner",
  "start_date": "2024-04-15T10:00:00Z",
  "end_date": "2024-04-15T17:00:00Z",
  "timezone": "Europe/Madrid",
  "is_online": false,
  "venue_address": "Calle Ejemplo, 123",
  "venue_name": "Centro de Conferencias",
  "venue_city": "Madrid",
  "venue_country": "España",
  "latitude": 40.4168,
  "longitude": -3.7038,
  "online_url": null,
  "streaming_url": null,
  "max_attendees": 100,
  "is_free": true,
  "price": null,
  "currency": "EUR",
  "registration_url": null,
  "image_url": "https://ejemplo.com/evento.jpg",
  "banner_url": "https://ejemplo.com/banner.jpg",
  "tags": ["desarrollo", "workshop", "programación"],
  "requirements": "Laptop con IDE instalado",
  "agenda": "10:00 - Introducción\n12:00 - Práctica\n15:00 - Q&A",
  "registration_start_date": "2024-03-01T00:00:00Z",
  "registration_end_date": "2024-04-10T23:59:59Z",
  "contact_email": "organizador@ejemplo.com",
  "contact_phone": "+34 600 123 456",
  "meta_title": "Mi Nuevo Evento - Workshop de Desarrollo",
  "meta_description": "Workshop completo de desarrollo para principiantes",
  "organization_id": "456e7890-e12b-34d5-b678-901234567890"
}
```

#### Campos Opcionales para Admin

```json
{
  "status": "published", // draft, published (solo admin)
  "is_public": true, // Visibilidad pública (solo admin)
  "is_featured": false // Evento destacado (solo admin)
}
```

#### Response Success (201)

```json
{
  "success": true,
  "message": "Evento creado exitosamente",
  "data": {
    // EventResponse completo del evento creado
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
    "title": ["El título es requerido"],
    "start_date": ["La fecha de inicio debe ser posterior a hoy"],
    "venue_address": ["La dirección es requerida para eventos presenciales"]
  }
}
```

---

### 6. Actualizar Evento

**PUT** `/events/{id}`

Actualiza un evento existente. Requiere ser propietario o admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

Cualquier campo del CreateEventRequest puede ser actualizado usando valores opcionales:

```json
{
  "title": "Título Actualizado",
  "description": "Nueva descripción",
  "start_date": "2024-04-20T10:00:00Z",
  "max_attendees": 150,
  "tags": ["nuevo", "tag"]
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Evento actualizado exitosamente",
  "data": {
    // EventResponse completo del evento actualizado
  }
}
```

---

### 7. Eliminar Evento

**DELETE** `/events/{id}`

Elimina un evento. Solo organizadores de la organización propietaria o admin.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Evento eliminado exitosamente",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "success": true,
    "message": "Evento eliminado exitosamente"
  }
}
```

---

### 8. Publicar Evento

**POST** `/events/{id}/publish`

Cambia el estado del evento a "published". Requiere permisos de publicación.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "send_notifications": true
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Evento publicado exitosamente",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "old_status": "draft",
    "new_status": "published",
    "changed_at": "2024-01-15T10:00:00Z",
    "changed_by": "usuario@ejemplo.com",
    "message": "Evento publicado exitosamente"
  }
}
```

---

### 9. Cancelar Evento

**POST** `/events/{id}/cancel`

Cancela un evento con razón obligatoria.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Request Body

```json
{
  "reason": "Motivo de cancelación que debe tener al menos 10 caracteres",
  "send_notifications": true
}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Evento cancelado exitosamente",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "old_status": "published",
    "new_status": "canceled",
    "changed_at": "2024-01-15T10:00:00Z",
    "changed_by": "usuario@ejemplo.com",
    "message": "Evento cancelado exitosamente"
  }
}
```

---

### 10. Agregar a Favoritos

**POST** `/events/{id}/favorite`

Agrega el evento a los favoritos del usuario.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Evento agregado a favoritos",
  "data": {
    "event_id": "123e4567-e89b-12d3-a456-426614174000",
    "event_title": "Conferencia de Ciberseguridad",
    "event_slug": "conferencia-cyberseguridad-madrid-2024",
    "start_date": "2024-03-15T09:00:00Z",
    "is_free": false,
    "image_url": "https://ejemplo.com/evento.jpg",
    "favorited_at": "2024-01-15T10:00:00Z"
  }
}
```

---

### 11. Remover de Favoritos

**DELETE** `/events/{id}/favorite`

Remueve el evento de los favoritos del usuario.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Evento removido de favoritos",
  "data": {
    "event_id": "123e4567-e89b-12d3-a456-426614174000",
    "removed_at": "2024-01-15T10:00:00Z"
  }
}
```

---

### 12. Eventos por Organización

**GET** `/events/organization/{orgId}`

Obtiene todos los eventos de una organización específica.

**Headers requeridos:**

```
Authorization: Bearer {access_token}
```

#### Query Parameters

```
?page=1                    // Página
&limit=20                  // Elementos por página
&status=published          // Filtrar por estado
&order_by=start_date       // Ordenar por
&order_dir=desc            // Dirección
```

#### Response Success (200)

```json
{
  "success": true,
  "message": "Eventos de la organización obtenidos",
  "data": {
    "events": [
      // Array de EventResponse
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 45,
      "total_pages": 3,
      "has_prev": false,
      "has_next": true
    },
    "organization": {
      "id": "456e7890-e12b-34d5-b678-901234567890",
      "name": "CybESphere Organization",
      "slug": "cybesphere-org"
    }
  }
}
```

---

## Códigos de Error Específicos

### 400 - Bad Request

- `validation_error`: Error en validación de campos
- `invalid_date_range`: Rango de fechas inválido (end_date debe ser posterior a start_date)
- `invalid_venue_data`: Datos de ubicación requeridos para eventos presenciales
- `invalid_online_url`: URL requerida para eventos online
- `invalid_price_data`: Precio requerido para eventos de pago
- `invalid_tags`: Máximo 10 tags permitidos
- `past_event_date`: No se pueden crear eventos en fechas pasadas

### 403 - Forbidden

- `insufficient_permissions`: No tienes permisos para esta acción
- `organization_required`: Debes pertenecer a una organización para crear eventos
- `not_event_owner`: Solo el propietario del evento puede realizar esta acción
- `event_already_published`: El evento ya está publicado
- `event_already_canceled`: El evento ya está cancelado
- `registration_period_active`: No se puede modificar evento durante período de registro activo

### 404 - Not Found

- `event_not_found`: Evento no encontrado
- `organization_not_found`: Organización no encontrada

### 409 - Conflict

- `event_slug_exists`: Ya existe un evento con ese slug
- `max_events_reached`: Has alcanzado el límite máximo de eventos

---

## Estados de Eventos

- `draft`: Borrador (no visible públicamente)
- `published`: Publicado (visible públicamente)
- `canceled`: Cancelado
- `completed`: Completado

## Tipos de Eventos

- `conference`: Conferencia
- `workshop`: Taller/Workshop
- `meetup`: Meetup/Encuentro
- `webinar`: Webinar
- `training`: Entrenamiento/Curso
- `competition`: Competición
- `other`: Otro

## Niveles de Eventos

- `beginner`: Principiante
- `intermediate`: Intermedio
- `advanced`: Avanzado

## Notas Importantes

1. **Precios**: Se almacenan en céntimos (15000 = 150.00 EUR)
2. **Fechas**: Todas en formato ISO 8601 UTC
3. **Coordenadas**: Latitud/Longitud en formato decimal
4. **Tags**: Máximo 10 tags por evento, cada uno máximo 50 caracteres
5. **Imágenes**: URLs válidas, se recomienda HTTPS
6. **Eventos Online**: Requieren `online_url` si `is_online` es true
7. **Eventos Presenciales**: Requieren `venue_address` si `is_online` es false
8. **Registro**: Si no se especifica URL de registro, se usa el sistema interno
9. **Capacidad**: `null` en `max_attendees` significa capacidad ilimitada
10. **Organizaciones**: Los eventos siempre deben estar asociados a una organización

## Permisos Requeridos

- **Crear Evento**: Usuario con organización activa
- **Actualizar Evento**: Miembro de la organización propietaria o admin
- **Eliminar Evento**: Miembro de la organización propietaria o admin
- **Publicar Evento**: Miembro de la organización propietaria o admin
- **Cancelar Evento**: Miembro de la organización propietaria o admin
- **Destacar Evento**: Solo admin
- **Ver Estadísticas**: Miembro de la organización propietaria o admin
