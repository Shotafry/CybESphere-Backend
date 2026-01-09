# Documentación de Migración e Integración: Frontend (React) + Backend (Go/Gin)

**Fecha:** 9 de Enero de 2026
**Autor:** Antigravity (Assistant)
**Objetivo:** Documentar el proceso completo de integración entre el Frontend y el Backend real, detallando cambios, soluciones a problemas encontrados y estado actual.

---

## 1. Contexto

El proyecto partió con un Frontend basado enteramente en interfaces simuladas ("mocks") y una capa de servicios fake. El objetivo principal fue eliminar esta simulación y conectar la aplicación a una API REST real desarrollada en Go con Gin Framework y PostgreSQL.

---

## 2. Fase de Tipado y Estandarización (Typescript)

### 2.1. Sincronización de DTOs

Se analizaron los DTOs del Backend (`internal/dto/*.go`) para replicar su estructura exacta en el Frontend.

- **Decisión clave:** Adoptar `snake_case` en el Frontend para todas las respuestas de API, coincidiendo con el JSON del backend, en lugar de mapear manualmente a `camelCase` en cada petición (lo que sería costoso y propenso a errores).

### 2.2. Actualización de `src/types/index.ts`

- Se redefinió la interfaz `User` para incluir campos nuevos (`role`, `is_active`, `social_links` aplanados).
- Se actualizó `Event` para soportar `organization_id`, `is_online`, `price` y `currency`.
- Se crearon interfaces envoltorio para respuestas de API: `ApiResponse<T>`, `PaginationMeta`.
- **Legacy Support:** Se mantuvieron algunos campos opcionales (`FavoriteEvents` alias) para no romper componentes antiguos mientras se migraba.

---

## 3. Capa de Servicios (Networking)

### 3.1. HttpClient Robusto (`src/services/httpClient.ts`)

Se implementó un cliente Axios con patrón Singleton.

- **Request Interceptor:** Inyecta automáticamente el `Authorization: Bearer <token>` si existe en localStorage.
- **Response Interceptor:** Manejo centralizado de errores. Detecta `401 Unauthorized` e intenta refrescar el token automáticamente (flow de Refresh Token) antes de cerrar la sesión del usuario.

### 3.2. Refactorización de `apiService.ts`

- Eliminación total de `mockData` y `setTimeout` simulados.
- Implementación de métodos reales: `getEvents`, `login`, `register`, `getMe`, `createEvent`, etc.
- Mapeo de endpoints específicos de administración y organizador (`/admin/*`, `/organizer/*`).

---

## 4. Integración de Autenticación (`AuthContext`)

- Conexión de `AuthProvider` con `apiService.getMe()`.
- Eliminación de la lógica "local-only" de login. Ahora el estado de autenticación depende puramente de la validación del token contra el backend.
- Persistencia segura de tokens y manejo de sesión.

---

## 5. Corrección de UI y Compilación (`Fixes`)

### 5.1. Actualización de Componentes (MUI Grid)

Se detectó un uso deprecado de `Grid item` que causaba errores de tipo en las nuevas versiones de MUI.

- **Solución:** Migración masiva a `Grid size={{ xs: 12, md: 6 }}` en:
  - `EventCard.tsx`
  - `PanelDeUsuario.tsx`
  - `UserProfile.tsx`
  - `Page.tsx`
  - `LandingPage.tsx`
  - `SignUp.tsx`

### 5.2. Limpieza de Errores de TypeScript (`tsc --noEmit`)

- Se resolvieron inconsistencias en `Page.tsx` donde `AgendaItem.id` podía ser undefined.
- Se ajustaron los **Mocks (`mocks/db.ts`)** para que cumplieran estrictamente con las nuevas interfaces de `User` y `Event` (añadiendo campos faltantes como `timezone`, `organization_name`, `is_ongoing`). Esto permite que los tests y el modo dev sigan funcionando si se necesitan mocks puntuales.

---

## 6. Fase de Datos y Assets (Seeders)

Para mejorar la experiencia visual y eliminar "placeholders", se modificaron los seeders del Backend en Go:

- **Archivos:** `internal/seeders/organization_seeder.go` y `event_seeder.go`.
- **Cambio:** Reemplazo de URLs dummy (`https://images.example.com...`) por rutas relativas a assets reales del proyecto (`/CloudEvents-logo-1@2x.png`, etc.).
- **Lógica:** Implementación de selección aleatoria segura (`SecureRandInt`) para asignar imágenes válidas a todos los eventos generados, garantizando que no haya imágenes rotas en la demo.

---

## 7. Resolución de Bugs Críticos (Fase Final)

### 7.1. Error 500 en Home (Endpoints Públicos)

**Problema:** El usuario no logueado recibía error 500 al intentar cargar eventos, porque `apiService` llamaba a `/events` (ruta protegida).
**Solución:** Se redirigierón las llamadas "públicas" a endpoints específicos `/public/*`:

- `getEvents` -> `/public/events`
- `getEventBySlug` -> `/public/events/:slug`
- `getAllOrganizations` -> `/public/organizations`

### 7.2. CORS (Cross-Origin Resource Sharing)

**Problema:** El navegador bloqueaba las peticiones desde `localhost:5173` hacia el backend.
**Solución:** Se integró `github.com/gin-contrib/cors` en `cmd/cybesphere/main.go`.

- Configuración explícita para permitir Origen `http://localhost:5173`.
- Métodos permitidos: `GET, POST, PUT, PATCH, DELETE, OPTIONS`.
- Headers expuestos: `Content-Length`, `Authorization`.

---

### 7.3. Interceptor de Respuestas (Unwrapping)

**Problema:** El backend devuelve las respuestas envueltas en `{ success: true, data: { ... } }`, pero el código frontend (`AuthContext`, etc.) esperaba recibir el objeto de datos directamente. Esto causaba errores como `Cannot read properties of undefined (reading 'role')`.
**Solución:** Se configuró un **Response Interceptor** en `src/services/httpClient.ts`.

- Automáticamente detecta si la respuesta tiene la estructura `{ data: ..., success: true }`.
- "Desenvuelve" el objeto retornando `response.data.data` al llamante.
- Esto permite que `apiService.ts` siga limpio (`return response.data`) sin tener que acceder manualmente a `.data.data` en cada función.

### 7.4. Refinamiento de Endpoints Públicos

**Problema:** A pesar de los cambios iniciales, algunas secciones (`getOrganizationEvents`, `getEventReviews`) seguían apuntando a rutas protegidas, causando errores 500.
**Solución:** Se auditó `src/services/apiService.ts` y se movieron TODAS las peticiones de lectura pública a `/public/*`:

- `getOrganizationEvents` -> `/public/organizations/:id/events`
- `getEventReviews` -> `/public/events/:id/reviews`

---

## 8. Integraciones Avanzadas (Backend)

### 8.1. Soporte para Slugs en Eventos e Identificación Inteligente

**Problema Crítico:** Al intentar acceder al detalle de un evento mediante su slug (ej: `/events/mi-gran-evento`), el backend devolvía Error 500 (SQL Error) porque el endpoint `GetByID` intentaba forzar la búsqueda de ese string en una columna `UUID`.

**Implementación Clean Architecture:**

1.  **Repository (`event_repository.go`):** Se añadió el método `GetBySlug` para búsquedas textuales optimizadas.
2.  **Service (`event_service.go`):** Se expuso `GetEventBySlug` en la interfaz y la implementación.
3.  **Handler (`event_handler.go`):** Se creó un nuevo endpoint inteligente `GetPublicEvent`.
    - **Lógica Híbrida:** Utiliza `uuid.Parse()` para determinar si el parámetro `id` es un UUID o un Slug.
    - Si es UUID -> Llama a `GetByID`.
    - Si es Texto -> Llama a `GetEventBySlug`.
4.  **Routes (`routes.go`):** Se actualizó la ruta pública `GET /events/:id` para usar este nuevo handler versátil.

---

## Estado Actual del Proyecto

1.  **Frontend:** Compila sin errores (`0 errors`). Typescript estricto. UI alineada con datos reales.
2.  **Backend:** Seeders generan datos ricos y visualmente correctos. API segura y endpoints públicos separados correctamente.
3.  **Conexión:** Login, Registro, Home Page, Dashboard y Listados funcionan contra la base de datos real.

**Siguientes pasos sugeridos:**

- Verificar flujo completo de creación de eventos (Frontend -> Backend).
- Testear subida de imágenes real (si aplica en el futuro).
- Desplegar en entorno de Staging.
