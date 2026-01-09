package mappers

import (
	"strings"
	"time"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/models"
)

// OrganizationMapperImpl implementación del mapper de organizaciones
type OrganizationMapperImpl struct{}

// NewOrganizationMapper crea nueva instancia del mapper
func NewOrganizationMapper() OrganizationMapperImpl {
	return OrganizationMapperImpl{}
}

// CreateOrganizationRequestToModel convierte CreateOrganizationRequest a modelo Organization
func (m OrganizationMapperImpl) CreateOrganizationRequestToModel(req *dto.CreateOrganizationRequest, userCtx *common.UserContext) (*models.Organization, error) {
	organization := &models.Organization{
		// Información básica
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Website:     strings.TrimSpace(req.Website),

		// Información de contacto
		Email:      strings.ToLower(strings.TrimSpace(req.Email)),
		Phone:      strings.TrimSpace(req.Phone),
		Address:    strings.TrimSpace(req.Address),
		City:       strings.TrimSpace(req.City),
		Country:    strings.TrimSpace(req.Country),
		PostalCode: strings.TrimSpace(req.PostalCode),

		// Geolocalización
		Latitude:  req.Latitude,
		Longitude: req.Longitude,

		// Branding
		LogoURL:        strings.TrimSpace(req.LogoURL),
		BannerURL:      strings.TrimSpace(req.BannerURL),
		PrimaryColor:   strings.TrimSpace(req.PrimaryColor),
		SecondaryColor: strings.TrimSpace(req.SecondaryColor),

		// Redes sociales
		LinkedIn:  strings.TrimSpace(req.LinkedIn),
		Twitter:   strings.TrimSpace(req.Twitter),
		Facebook:  strings.TrimSpace(req.Facebook),
		Instagram: strings.TrimSpace(req.Instagram),
		YouTube:   strings.TrimSpace(req.YouTube),

		// Documentación para verificación
		TaxID:            strings.TrimSpace(req.TaxID),
		LegalName:        strings.TrimSpace(req.LegalName),
		RegistrationDocs: strings.TrimSpace(req.RegistrationDocs),

		// Estado inicial
		Status:          models.OrgStatusPending, // Requiere verificación por defecto
		IsVerified:      false,
		CanCreateEvents: true, // Habilitado por defecto
	}

	// Solo admin puede establecer ciertos campos en creación
	if userCtx != nil && userCtx.IsAdmin() {
		if req.Status != "" {
			organization.Status = models.OrganizationStatus(req.Status)
		}
		if req.IsVerified != nil {
			organization.IsVerified = *req.IsVerified
		}
		if req.MaxEvents != nil {
			organization.MaxEvents = req.MaxEvents
		}
		if req.CanCreateEvents != nil {
			organization.CanCreateEvents = *req.CanCreateEvents
		}
	}

	// Validar colores si están presentes
	if err := organization.SetBranding(organization.PrimaryColor, organization.SecondaryColor); err != nil {
		return nil, common.NewBusinessError("invalid_branding", "Formato de color inválido")
	}

	return organization, nil
}

// UpdateOrganizationRequestToModel aplica cambios de UpdateOrganizationRequest
func (m OrganizationMapperImpl) UpdateOrganizationRequestToModel(org *models.Organization, req *dto.UpdateOrganizationRequest, userCtx *common.UserContext) error {
	// Campos básicos que miembros de la organización pueden cambiar
	if req.Name != nil {
		org.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		org.Description = strings.TrimSpace(*req.Description)
	}
	if req.Website != nil {
		org.Website = strings.TrimSpace(*req.Website)
	}

	// Información de contacto
	if req.Email != nil {
		org.Email = strings.ToLower(strings.TrimSpace(*req.Email))
	}
	if req.Phone != nil {
		org.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Address != nil {
		org.Address = strings.TrimSpace(*req.Address)
	}
	if req.City != nil {
		org.City = strings.TrimSpace(*req.City)
	}
	if req.Country != nil {
		org.Country = strings.TrimSpace(*req.Country)
	}
	if req.PostalCode != nil {
		org.PostalCode = strings.TrimSpace(*req.PostalCode)
	}

	// Geolocalización
	if req.Latitude != nil {
		org.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		org.Longitude = req.Longitude
	}

	// Branding
	if req.LogoURL != nil {
		org.LogoURL = strings.TrimSpace(*req.LogoURL)
	}
	if req.BannerURL != nil {
		org.BannerURL = strings.TrimSpace(*req.BannerURL)
	}

	// Colores de branding
	if req.PrimaryColor != nil || req.SecondaryColor != nil {
		primaryColor := org.PrimaryColor
		secondaryColor := org.SecondaryColor

		if req.PrimaryColor != nil {
			primaryColor = strings.TrimSpace(*req.PrimaryColor)
		}
		if req.SecondaryColor != nil {
			secondaryColor = strings.TrimSpace(*req.SecondaryColor)
		}

		if err := org.SetBranding(primaryColor, secondaryColor); err != nil {
			return common.NewBusinessError("invalid_branding", "Formato de color inválido")
		}
	}

	// Redes sociales
	if req.LinkedIn != nil {
		org.LinkedIn = strings.TrimSpace(*req.LinkedIn)
	}
	if req.Twitter != nil {
		org.Twitter = strings.TrimSpace(*req.Twitter)
	}
	if req.Facebook != nil {
		org.Facebook = strings.TrimSpace(*req.Facebook)
	}
	if req.Instagram != nil {
		org.Instagram = strings.TrimSpace(*req.Instagram)
	}
	if req.YouTube != nil {
		org.YouTube = strings.TrimSpace(*req.YouTube)
	}

	// Campos que solo admin puede cambiar
	if userCtx != nil && userCtx.IsAdmin() {
		if req.Status != nil {
			org.Status = models.OrganizationStatus(*req.Status)
		}
		if req.IsVerified != nil {
			org.IsVerified = *req.IsVerified
			if *req.IsVerified {
				// Si se está verificando, actualizar campos relacionados
				now := time.Now().Format(time.RFC3339)
				org.VerifiedAt = &now
				verifierID := userCtx.ID
				org.VerifiedBy = &verifierID
			}
		}
		if req.MaxEvents != nil {
			org.MaxEvents = req.MaxEvents
		}
		if req.CanCreateEvents != nil {
			org.CanCreateEvents = *req.CanCreateEvents
		}
	}

	return nil
}

// OrganizationToResponse convierte Organization model a OrganizationResponse
func (m OrganizationMapperImpl) OrganizationToResponse(org *models.Organization, userCtx *common.UserContext) dto.OrganizationResponse {
	response := dto.OrganizationResponse{
		// Identificación
		ID:   org.ID.String(),
		Slug: org.Slug,

		// Información básica (siempre visible)
		Name:        org.Name,
		Description: org.Description,
		Website:     org.Website,
		City:        org.City,
		Country:     org.Country,

		// Branding
		LogoURL:        org.LogoURL,
		BannerURL:      org.BannerURL,
		PrimaryColor:   org.PrimaryColor,
		SecondaryColor: org.SecondaryColor,

		// Estado y verificación
		Status:     string(org.Status),
		IsVerified: org.IsVerified,
		VerifiedAt: m.formatVerifiedAt(org.VerifiedAt),

		// Estadísticas públicas
		EventsCount: org.EventsCount,

		// Información del usuario
		IsMember:  m.isMember(org, userCtx),
		CanEdit:   m.canEdit(org, userCtx),
		CanManage: m.canManage(org, userCtx),

		// Timestamps
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}

	// Información de contacto - solo para admin y miembros
	if m.canViewContactInfo(org, userCtx) {
		response.Email = org.Email
		response.Phone = org.Phone
		response.Address = org.Address
		response.PostalCode = org.PostalCode
	}

	// Geolocalización
	if m.canViewLocation(org, userCtx) {
		response.Latitude = org.Latitude
		response.Longitude = org.Longitude
	}

	// Redes sociales si están configuradas
	if m.hasSocialMedia(org) {
		response.SocialMedia = &dto.SocialMediaLinks{
			LinkedIn:  org.LinkedIn,
			Twitter:   org.Twitter,
			Facebook:  org.Facebook,
			Instagram: org.Instagram,
			YouTube:   org.YouTube,
		}
	}

	// Configuración administrativa - solo para admin
	if userCtx != nil && userCtx.IsAdmin() {
		response.MaxEvents = org.MaxEvents
		response.CanCreateEvents = org.CanCreateEvents
	}

	return response
}

// OrganizationToDetailResponse convierte Organization a respuesta detallada
func (m OrganizationMapperImpl) OrganizationToDetailResponse(org *models.Organization, userCtx *common.UserContext) dto.OrganizationDetailResponse {
	baseResponse := m.OrganizationToResponse(org, userCtx)

	detailResponse := dto.OrganizationDetailResponse{
		OrganizationResponse: baseResponse,
	}

	// Información adicional de verificación - solo admin
	if userCtx != nil && userCtx.IsAdmin() {
		detailResponse.TaxID = org.TaxID
		detailResponse.LegalName = org.LegalName
		detailResponse.RegistrationDocs = org.RegistrationDocs
		detailResponse.VerifiedBy = m.getVerifiedByName(org.VerifiedBy)
	}

	// Estadísticas detalladas - para miembros y admin
	if m.canViewStatistics(org, userCtx) {
		detailResponse.Statistics = m.buildOrganizationStatistics(org)
	}

	// Eventos recientes - para miembros y admin
	if m.canViewEvents(org, userCtx) {
		detailResponse.RecentEvents = m.getRecentEvents(org)
	}

	// Miembros destacados - información pública limitada
	detailResponse.FeaturedMembers = m.getFeaturedMembers(org)

	return detailResponse
}

// OrganizationToSummaryResponse convierte a respuesta resumida
func (m OrganizationMapperImpl) OrganizationToSummaryResponse(org *models.Organization) dto.OrganizationSummaryResponse {
	return dto.OrganizationSummaryResponse{
		ID:          org.ID.String(),
		Slug:        org.Slug,
		Name:        org.Name,
		LogoURL:     org.LogoURL,
		IsVerified:  org.IsVerified,
		EventsCount: org.EventsCount,
		City:        org.City,
		Country:     org.Country,
	}
}

// OrganizationsToListResponse convierte lista con paginación
func (m OrganizationMapperImpl) OrganizationsToListResponse(orgs []*models.Organization, pagination *common.PaginationMeta, userCtx *common.UserContext) dto.OrganizationListResponse {
	organizationResponses := make([]dto.OrganizationResponse, 0, len(orgs))

	for _, org := range orgs {
		organizationResponses = append(organizationResponses, m.OrganizationToResponse(org, userCtx))
	}

	return dto.OrganizationListResponse{
		Organizations: organizationResponses,
		Pagination:    *pagination,
		Filters:       dto.AppliedFilters{}, // TODO: implementar filtros aplicados
	}
}

// =============================================================================
// MÉTODOS HELPER PRIVADOS
// =============================================================================

// isMember verifica si el usuario es miembro de la organización
func (m OrganizationMapperImpl) isMember(org *models.Organization, userCtx *common.UserContext) bool {
	if userCtx == nil {
		return false
	}

	return userCtx.OrganizationID != nil && *userCtx.OrganizationID == org.ID.String()
}

// canEdit verifica si puede editar la organización
func (m OrganizationMapperImpl) canEdit(org *models.Organization, userCtx *common.UserContext) bool {
	if userCtx == nil {
		return false
	}

	if userCtx.IsAdmin() {
		return true
	}

	return userCtx.IsOrganizer() && m.isMember(org, userCtx)
}

// canManage verifica permisos de gestión completa
func (m OrganizationMapperImpl) canManage(org *models.Organization, userCtx *common.UserContext) bool {
	return m.canEdit(org, userCtx) // Por ahora son iguales
}

// canViewContactInfo verifica si puede ver información de contacto
func (m OrganizationMapperImpl) canViewContactInfo(org *models.Organization, userCtx *common.UserContext) bool {
	return m.canManage(org, userCtx)
}

// canViewLocation verifica si puede ver geolocalización
func (m OrganizationMapperImpl) canViewLocation(org *models.Organization, userCtx *common.UserContext) bool {
	// Geolocalización es pública para organizaciones activas
	return org.Status == models.OrgStatusActive
}

// canViewStatistics verifica si puede ver estadísticas
func (m OrganizationMapperImpl) canViewStatistics(org *models.Organization, userCtx *common.UserContext) bool {
	return m.canManage(org, userCtx)
}

// canViewEvents verifica si puede ver eventos detallados
func (m OrganizationMapperImpl) canViewEvents(org *models.Organization, userCtx *common.UserContext) bool {
	return m.canManage(org, userCtx)
}

// hasSocialMedia verifica si tiene redes sociales configuradas
func (m OrganizationMapperImpl) hasSocialMedia(org *models.Organization) bool {
	return org.LinkedIn != "" || org.Twitter != "" || org.Facebook != "" ||
		org.Instagram != "" || org.YouTube != ""
}

// formatVerifiedAt formatea fecha de verificación
func (m OrganizationMapperImpl) formatVerifiedAt(verifiedAt *string) *time.Time {
	if verifiedAt == nil || *verifiedAt == "" {
		return nil
	}

	if t, err := time.Parse(time.RFC3339, *verifiedAt); err == nil {
		return &t
	}

	return nil
}

// getVerifiedByName obtiene nombre del verificador (placeholder)
func (m OrganizationMapperImpl) getVerifiedByName(verifiedBy *string) string {
	if verifiedBy == nil {
		return ""
	}
	// En implementación real, harías lookup del nombre del admin
	return *verifiedBy
}

// buildOrganizationStatistics construye estadísticas de la organización
func (m OrganizationMapperImpl) buildOrganizationStatistics(org *models.Organization) *dto.OrganizationStatistics {
	// Implementación placeholder - en producción consultarías datos reales
	return &dto.OrganizationStatistics{
		TotalEvents:      org.EventsCount,
		PublishedEvents:  org.EventsCount, // Placeholder
		CompletedEvents:  0,               // Calcular eventos completados
		CanceledEvents:   0,               // Calcular eventos cancelados
		TotalAttendees:   0,               // Sumar asistentes de todos los eventos
		AverageAttendees: 0.0,             // Calcular promedio
		EventsByType:     make(map[string]int),
		EventsByMonth:    []dto.MonthlyCount{},
		TopCities:        []dto.CityCount{},
	}
}

// getRecentEvents obtiene eventos recientes de la organización
func (m OrganizationMapperImpl) getRecentEvents(org *models.Organization) []dto.EventSummaryResponse {
	// Placeholder - en implementación real harías query a la base de datos
	return []dto.EventSummaryResponse{}
}

// getFeaturedMembers obtiene miembros destacados públicos
func (m OrganizationMapperImpl) getFeaturedMembers(org *models.Organization) []dto.UserSummaryResponse {
	// Placeholder - en implementación real harías query de usuarios públicos
	return []dto.UserSummaryResponse{}
}
