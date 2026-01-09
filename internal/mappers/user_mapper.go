package mappers

import (
	"strings"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/models"
	"cybesphere-backend/pkg/database"
)

// UserMapperImpl implementación del mapper de usuarios
type UserMapperImpl struct{}

// NewUserMapper crea nueva instancia del mapper
func NewUserMapper() UserMapperImpl {
	return UserMapperImpl{}
}

// CreateUserRequestToModel convierte CreateUserRequest a modelo User
func (m UserMapperImpl) CreateUserRequestToModel(req *dto.CreateUserRequest, userCtx *common.UserContext) (*models.User, error) {
	user := &models.User{
		// Información básica
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Password:  req.Password, // Se hasheará automáticamente en BeforeCreate
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),

		// Perfil profesional
		Company:  strings.TrimSpace(req.Company),
		Position: strings.TrimSpace(req.Position),
		Bio:      strings.TrimSpace(req.Bio),
		Website:  strings.TrimSpace(req.Website),
		LinkedIn: strings.TrimSpace(req.LinkedIn),
		Twitter:  strings.TrimSpace(req.Twitter),

		// Ubicación
		City:      strings.TrimSpace(req.City),
		Country:   strings.TrimSpace(req.Country),
		Latitude:  req.Latitude,
		Longitude: req.Longitude,

		// Configuraciones
		Timezone:          strings.TrimSpace(req.Timezone),
		Language:          strings.TrimSpace(req.Language),
		NewsletterEnabled: req.NewsletterEnabled,

		// Estado inicial
		Role:       models.RoleUser, // Por defecto
		IsActive:   true,            // Activo por defecto
		IsVerified: false,           // Requiere verificación
	}

	// Aplicar valores por defecto
	if user.Timezone == "" {
		user.Timezone = "Europe/Madrid"
	}
	if user.Language == "" {
		user.Language = "es"
	}

	// Solo admin puede establecer ciertos campos en creación
	if userCtx != nil && userCtx.IsAdmin() {
		if req.Role != "" {
			user.Role = models.UserRole(req.Role)
		}
		if req.IsActive != nil {
			user.IsActive = *req.IsActive
		}
		if req.IsVerified != nil {
			user.IsVerified = *req.IsVerified
		}
		if req.OrganizationID != nil {
			user.OrganizationID = req.OrganizationID
		}
	}

	return user, nil
}

// UpdateUserRequestToModel aplica cambios de UpdateUserRequest a usuario existente
func (m UserMapperImpl) UpdateUserRequestToModel(user *models.User, req *dto.UpdateUserRequest, userCtx *common.UserContext) error {
	// Campos básicos que el usuario puede cambiar
	if req.FirstName != nil {
		user.FirstName = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil {
		user.LastName = strings.TrimSpace(*req.LastName)
	}

	// Perfil profesional
	if req.Company != nil {
		user.Company = strings.TrimSpace(*req.Company)
	}
	if req.Position != nil {
		user.Position = strings.TrimSpace(*req.Position)
	}
	if req.Bio != nil {
		user.Bio = strings.TrimSpace(*req.Bio)
	}
	if req.Website != nil {
		user.Website = strings.TrimSpace(*req.Website)
	}
	if req.LinkedIn != nil {
		user.LinkedIn = strings.TrimSpace(*req.LinkedIn)
	}
	if req.Twitter != nil {
		user.Twitter = strings.TrimSpace(*req.Twitter)
	}

	// Ubicación
	if req.City != nil {
		user.City = strings.TrimSpace(*req.City)
	}
	if req.Country != nil {
		user.Country = strings.TrimSpace(*req.Country)
	}
	if req.Latitude != nil {
		user.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		user.Longitude = req.Longitude
	}

	// Configuraciones
	if req.Timezone != nil {
		user.Timezone = strings.TrimSpace(*req.Timezone)
	}
	if req.Language != nil {
		user.Language = strings.TrimSpace(*req.Language)
	}
	if req.NewsletterEnabled != nil {
		user.NewsletterEnabled = *req.NewsletterEnabled
	}

	// Campos que solo admin puede cambiar
	if userCtx != nil && userCtx.IsAdmin() {
		if req.Role != nil {
			user.Role = models.UserRole(*req.Role)
		}
		if req.IsActive != nil {
			user.IsActive = *req.IsActive
		}
		if req.IsVerified != nil {
			user.IsVerified = *req.IsVerified
		}
		if req.OrganizationID != nil {
			user.OrganizationID = req.OrganizationID
		}
	}

	return nil
}

// UserToResponse convierte User model a UserResponse
func (m UserMapperImpl) UserToResponse(user *models.User, userCtx *common.UserContext) dto.UserResponse {
	response := dto.UserResponse{
		// Identificación
		ID: user.ID.String(),

		// Información básica (siempre visible para propietario/admin)
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.GetFullName(),

		// Estado del sistema
		Role:       string(user.Role),
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,

		// Configuraciones
		Timezone:          user.Timezone,
		Language:          user.Language,
		NewsletterEnabled: user.NewsletterEnabled,

		// Información del usuario actual
		CanEdit:   m.canEdit(user, userCtx),
		CanManage: m.canManage(user, userCtx),

		// Timestamps
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	// Información que se muestra según permisos
	isOwner := m.isOwner(user, userCtx)
	isAdmin := userCtx != nil && userCtx.IsAdmin()

	// Perfil profesional - visible para propietario y admin
	if isOwner || isAdmin {
		response.Company = user.Company
		response.Position = user.Position
		response.Bio = user.Bio
		response.Website = user.Website
		response.LinkedIn = user.LinkedIn
		response.Twitter = user.Twitter
		response.City = user.City
		response.Country = user.Country
		response.Latitude = user.Latitude
		response.Longitude = user.Longitude
	}

	// Organización si está asignada
	if user.Organization != nil {
		response.Organization = &dto.OrganizationSummaryResponse{
			ID:          user.Organization.ID.String(),
			Slug:        user.Organization.Slug,
			Name:        user.Organization.Name,
			LogoURL:     user.Organization.LogoURL,
			IsVerified:  user.Organization.IsVerified,
			EventsCount: user.Organization.EventsCount,
			City:        user.Organization.City,
			Country:     user.Organization.Country,
		}
	}

	// Estadísticas para propietario y admin
	if (isOwner || isAdmin) && m.shouldIncludeStatistics(user, userCtx) {
		response.Statistics = m.buildUserStatistics(user)
	}

	return response
}

// UserToDetailResponse convierte User a respuesta detallada
func (m UserMapperImpl) UserToDetailResponse(user *models.User, userCtx *common.UserContext) dto.UserDetailResponse {
	baseResponse := m.UserToResponse(user, userCtx)

	detailResponse := dto.UserDetailResponse{
		UserResponse: baseResponse,
	}

	// Información adicional solo para propietario
	if m.isOwner(user, userCtx) {
		detailResponse.FavoriteEvents = m.getFavoriteEvents(user)
		detailResponse.RegisteredEvents = m.getRegisteredEvents(user)
		detailResponse.ActiveSessions = m.getActiveSessions(user, userCtx)
	}

	return detailResponse
}

// UserToSummaryResponse convierte a respuesta resumida
func (m UserMapperImpl) UserToSummaryResponse(user *models.User) dto.UserSummaryResponse {
	response := dto.UserSummaryResponse{
		ID:         user.ID.String(),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		FullName:   user.GetFullName(),
		Role:       string(user.Role),
		IsVerified: user.IsVerified,
	}

	// Información profesional básica (pública)
	if user.Company != "" {
		response.Company = user.Company
	}
	if user.Position != "" {
		response.Position = user.Position
	}

	// Email solo para admin - se filtrará en el contexto que lo use
	response.Email = user.Email

	// Organización si existe
	if user.Organization != nil {
		response.Organization = user.Organization.Name
	}

	return response
}

// UserToProfileResponse convierte a perfil público
func (m UserMapperImpl) UserToProfileResponse(user *models.User) dto.UserProfileResponse {
	response := dto.UserProfileResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.GetFullName(),
		Company:   user.Company,
		Position:  user.Position,
		Bio:       user.Bio,
		City:      user.City,
		Country:   user.Country,
		LinkedIn:  user.LinkedIn,
		Twitter:   user.Twitter,
		Website:   user.Website,
	}

	// Eventos públicos si el usuario los ha hecho públicos
	response.PublicEvents = m.getPublicEvents(user)

	return response
}

// UsersToListResponse convierte lista de usuarios con paginación
func (m UserMapperImpl) UsersToListResponse(users []*models.User, pagination *common.PaginationMeta, userCtx *common.UserContext) dto.UserListResponse {
	userResponses := make([]dto.UserResponse, 0, len(users))

	for _, user := range users {
		userResponses = append(userResponses, m.UserToResponse(user, userCtx))
	}

	return dto.UserListResponse{
		Users:      userResponses,
		Pagination: *pagination,
		Filters:    dto.AppliedFilters{}, // TODO: implementar filtros aplicados
	}
}

// =============================================================================
// MÉTODOS HELPER PRIVADOS
// =============================================================================

// isOwner verifica si el usuario actual es el propietario del perfil
func (m UserMapperImpl) isOwner(user *models.User, userCtx *common.UserContext) bool {
	return userCtx != nil && userCtx.ID == user.ID.String()
}

// canEdit verifica si puede editar el usuario
func (m UserMapperImpl) canEdit(user *models.User, userCtx *common.UserContext) bool {
	if userCtx == nil {
		return false
	}

	// Admin puede editar cualquiera, usuario solo su propio perfil
	return userCtx.IsAdmin() || m.isOwner(user, userCtx)
}

// canManage verifica permisos de gestión completa (activar/desactivar, roles)
func (m UserMapperImpl) canManage(user *models.User, userCtx *common.UserContext) bool {
	return userCtx != nil && userCtx.IsAdmin()
}

// shouldIncludeStatistics determina si incluir estadísticas
func (m UserMapperImpl) shouldIncludeStatistics(user *models.User, userCtx *common.UserContext) bool {
	// Solo incluir para propietario o admin
	return m.isOwner(user, userCtx) || (userCtx != nil && userCtx.IsAdmin())
}

// buildUserStatistics construye estadísticas del usuario - IMPLEMENTACIÓN MEJORADA
func (m UserMapperImpl) buildUserStatistics(user *models.User) *dto.UserStatistics {
	db := database.GetDB()

	stats := &dto.UserStatistics{
		MemberSince: user.CreatedAt,
		LastActive:  user.UpdatedAt,
	}

	// Última actividad basada en último login si disponible
	if user.LastLoginAt != nil {
		stats.LastActive = *user.LastLoginAt
	}

	// Contar eventos favoritos usando la relación many-to-many
	var favoriteCount int64
	db.Model(&models.User{}).Where("id = ?", user.ID).
		Joins("JOIN user_favorite_events ON users.id = user_favorite_events.user_id").
		Count(&favoriteCount)
	stats.FavoriteEvents = int(favoriteCount)

	// TODO: Implementar conteo de eventos asistidos cuando tengamos esa funcionalidad
	stats.EventsAttended = 0

	// Si el usuario es organizador o admin, contar eventos organizados de su organización
	if (user.Role == models.RoleOrganizer || user.Role == models.RoleAdmin) && user.OrganizationID != nil {
		var eventCount int64
		db.Model(&models.Event{}).
			Where("organization_id = ? AND status = ?", *user.OrganizationID, models.EventStatusPublished).
			Count(&eventCount)
		stats.EventsOrganized = int(eventCount)
	}

	// TODO: Implementar total de conexiones si existe esa funcionalidad
	stats.TotalConnections = 0

	return stats
}

// getFavoriteEvents obtiene eventos favoritos del usuario - IMPLEMENTACIÓN REAL
func (m UserMapperImpl) getFavoriteEvents(user *models.User) []dto.EventSummaryResponse {
	db := database.GetDB()
	var events []models.Event

	// Obtener eventos favoritos con información de organización
	err := db.Preload("Organization").
		Joins("JOIN user_favorite_events ON events.id = user_favorite_events.event_id").
		Where("user_favorite_events.user_id = ? AND events.status = ?",
				user.ID, models.EventStatusPublished).
		Limit(10). // Limitar a 10 eventos más recientes
		Order("user_favorite_events.created_at DESC").
		Find(&events).Error

	if err != nil {
		return []dto.EventSummaryResponse{}
	}

	// Convertir a respuestas resumidas
	summaries := make([]dto.EventSummaryResponse, 0, len(events))
	for _, event := range events {
		summary := dto.EventSummaryResponse{
			ID:               event.ID.String(),
			Slug:             event.Slug,
			Title:            event.Title,
			ShortDesc:        event.ShortDesc,
			Type:             string(event.Type),
			StartDate:        event.StartDate,
			EndDate:          event.EndDate,
			IsOnline:         event.IsOnline,
			VenueCity:        event.VenueCity,
			IsFree:           event.IsFree,
			Price:            event.Price,
			ImageURL:         event.ImageURL,
			CurrentAttendees: event.CurrentAttendees,
			MaxAttendees:     event.MaxAttendees,
			IsFeatured:       event.IsFeatured,
			Tags:             event.GetTags(),
		}

		if event.Organization != nil {
			summary.Organization = event.Organization.Name
		}

		summaries = append(summaries, summary)
	}

	return summaries
}

// getRegisteredEvents obtiene eventos registrados del usuario
func (m UserMapperImpl) getRegisteredEvents(user *models.User) []dto.EventSummaryResponse {
	// TODO: Esta funcionalidad requiere un modelo de registro de eventos
	// Por ahora retornamos vacío, pero la estructura está lista
	// Cuando implementes el modelo EventRegistration, puedes completar esta función
	return []dto.EventSummaryResponse{}
}

// getActiveSessions obtiene sesiones activas (solo para el propietario) - IMPLEMENTACIÓN REAL
// getActiveSessions obtiene sesiones activas (solo para el propietario) - IMPLEMENTACIÓN REAL
func (m UserMapperImpl) getActiveSessions(user *models.User, userCtx *common.UserContext) []dto.SessionResponse {
	if !m.isOwner(user, userCtx) {
		return nil
	}

	db := database.GetDB()
	var tokens []models.RefreshToken // Changed from []*models.RefreshToken

	// Obtener tokens activos del usuario
	err := db.Where("user_id = ? AND is_revoked = ? AND expires_at > ?",
		user.ID, false, db.NowFunc()).
		Order("last_used_at DESC, created_at DESC").
		Find(&tokens).Error

	if err != nil {
		return []dto.SessionResponse{}
	}

	// Convert to slice of pointers for the auth mapper
	tokenPtrs := make([]*models.RefreshToken, len(tokens))
	for i := range tokens {
		tokenPtrs[i] = &tokens[i]
	}

	// Convertir a SessionResponse usando la lógica del auth_mapper
	authMapper := NewAuthMapper()
	sessionList := authMapper.RefreshTokensToSessionList(tokenPtrs, "") // Fixed: pass tokenPtrs instead of &tokens

	return sessionList.Sessions
}

// getPublicEvents obtiene eventos públicos del usuario - IMPLEMENTACIÓN BÁSICA
func (m UserMapperImpl) getPublicEvents(user *models.User) []dto.EventSummaryResponse {
	// Si el usuario es organizador, mostrar algunos eventos de su organización
	if user.Role == models.RoleOrganizer && user.OrganizationID != nil {
		db := database.GetDB()
		var events []models.Event

		err := db.Preload("Organization").
			Where("organization_id = ? AND status = ? AND is_public = ?",
					*user.OrganizationID, models.EventStatusPublished, true).
			Limit(5). // Limitar a 5 eventos
			Order("start_date DESC").
			Find(&events).Error

		if err != nil {
			return []dto.EventSummaryResponse{}
		}

		// Convertir a respuestas resumidas
		summaries := make([]dto.EventSummaryResponse, 0, len(events))
		for _, event := range events {
			summary := dto.EventSummaryResponse{
				ID:               event.ID.String(),
				Slug:             event.Slug,
				Title:            event.Title,
				ShortDesc:        event.ShortDesc,
				Type:             string(event.Type),
				StartDate:        event.StartDate,
				EndDate:          event.EndDate,
				IsOnline:         event.IsOnline,
				VenueCity:        event.VenueCity,
				IsFree:           event.IsFree,
				Price:            event.Price,
				ImageURL:         event.ImageURL,
				CurrentAttendees: event.CurrentAttendees,
				MaxAttendees:     event.MaxAttendees,
				IsFeatured:       event.IsFeatured,
				Tags:             event.GetTags(),
			}

			if event.Organization != nil {
				summary.Organization = event.Organization.Name
			}

			summaries = append(summaries, summary)
		}

		return summaries
	}

	// Para usuarios regulares, no mostrar eventos públicos por defecto
	return []dto.EventSummaryResponse{}
}
