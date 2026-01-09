package mappers

import (
	"time"

	"cybesphere-backend/internal/common"
	"cybesphere-backend/internal/dto"
	"cybesphere-backend/internal/models"
)

// ResponseMapper interfaz para mapeo de DTOs a models y viceversa
type ResponseMapper interface {
	// Métodos genéricos requeridos por services
	DTOToEntity(dto interface{}, userCtx *common.UserContext) (interface{}, error)
	ApplyUpdateDTO(entity interface{}, dto interface{}, userCtx *common.UserContext) (interface{}, error)
	EntityToResponse(entity interface{}, userCtx *common.UserContext) (interface{}, error)
}

// EventMapper interfaz específica para mapeo de eventos
type EventMapper interface {
	CreateEventRequestToModel(req *dto.CreateEventRequest, userCtx *common.UserContext) (*models.Event, error)
	UpdateEventRequestToModel(event *models.Event, req *dto.UpdateEventRequest, userCtx *common.UserContext) error
	EventToResponse(event *models.Event, userCtx *common.UserContext) dto.EventResponse
	EventToDetailResponse(event *models.Event, userCtx *common.UserContext) dto.EventDetailResponse
	EventToSummaryResponse(event *models.Event) dto.EventSummaryResponse
	EventsToListResponse(events []*models.Event, pagination *common.PaginationMeta, userCtx *common.UserContext) dto.EventListResponse
}

// OrganizationMapper interfaz específica para mapeo de organizaciones
type OrganizationMapper interface {
	CreateOrganizationRequestToModel(req *dto.CreateOrganizationRequest, userCtx *common.UserContext) (*models.Organization, error)
	UpdateOrganizationRequestToModel(org *models.Organization, req *dto.UpdateOrganizationRequest, userCtx *common.UserContext) error
	OrganizationToResponse(org *models.Organization, userCtx *common.UserContext) dto.OrganizationResponse
	OrganizationToDetailResponse(org *models.Organization, userCtx *common.UserContext) dto.OrganizationDetailResponse
	OrganizationToSummaryResponse(org *models.Organization) dto.OrganizationSummaryResponse
	OrganizationsToListResponse(orgs []*models.Organization, pagination *common.PaginationMeta, userCtx *common.UserContext) dto.OrganizationListResponse
}

// UserMapper interfaz específica para mapeo de usuarios
type UserMapper interface {
	CreateUserRequestToModel(req *dto.CreateUserRequest, userCtx *common.UserContext) (*models.User, error)
	UpdateUserRequestToModel(user *models.User, req *dto.UpdateUserRequest, userCtx *common.UserContext) error
	UserToResponse(user *models.User, userCtx *common.UserContext) dto.UserResponse
	UserToDetailResponse(user *models.User, userCtx *common.UserContext) dto.UserDetailResponse
	UserToSummaryResponse(user *models.User) dto.UserSummaryResponse
	UserToProfileResponse(user *models.User) dto.UserProfileResponse
	UsersToListResponse(users []*models.User, pagination *common.PaginationMeta, userCtx *common.UserContext) dto.UserListResponse
}

// AuthMapper interfaz específica para mapeo de autenticación
type AuthMapper interface {
	RegisterRequestToUser(req *dto.RegisterRequest) (*models.User, error)
	UserToAuthResponse(user *models.User, accessToken, refreshToken string, expiresIn int) dto.AuthResponse
	UserToRegisterResponse(user *models.User, accessToken, refreshToken string, expiresIn int) dto.RegisterResponse
	BuildPasswordResetResponse(success bool, expiresInMinutes int) dto.PasswordResetResponse
	BuildEmailVerificationResponse(success bool, isVerified bool) dto.EmailVerificationResponse
	BuildLogoutResponse(success bool) dto.LogoutResponse
	BuildTokenResponse(accessToken, refreshToken string, expiresAt time.Time) dto.TokenResponse
	RefreshTokensToSessionList(tokens []*models.RefreshToken, currentTokenID string) dto.SessionListResponse
}

// UnifiedMapper estructura que implementa todas las interfaces
type UnifiedMapper struct {
	// Usar implementaciones concretas en lugar de interfaces
	eventMapper EventMapperImpl
	orgMapper   OrganizationMapperImpl
	userMapper  UserMapperImpl
	authMapper  AuthMapperImpl
}

// NewUnifiedMapper crea una nueva instancia del mapper unificado
func NewUnifiedMapper() *UnifiedMapper {
	return &UnifiedMapper{
		eventMapper: NewEventMapper(),
		orgMapper:   NewOrganizationMapper(),
		userMapper:  NewUserMapper(),
		authMapper:  NewAuthMapper(),
	}
}

// =============================================================================
// IMPLEMENTACIÓN DE ResponseMapper (interfaz genérica)
// =============================================================================

// DTOToEntity implementa ResponseMapper - FIX: usar type assertion más específica
func (m *UnifiedMapper) DTOToEntity(dtoParam interface{}, userCtx *common.UserContext) (interface{}, error) {
	switch d := dtoParam.(type) {
	case *dto.CreateEventRequest:
		return m.eventMapper.CreateEventRequestToModel(d, userCtx)
	case *dto.CreateOrganizationRequest:
		return m.orgMapper.CreateOrganizationRequestToModel(d, userCtx)
	case *dto.CreateUserRequest:
		return m.userMapper.CreateUserRequestToModel(d, userCtx)
	case *dto.RegisterRequest:
		return m.authMapper.RegisterRequestToUser(d)
	}

	// Fallback para debugging - ver qué tipo estamos recibiendo realmente
	return nil, common.NewBusinessError("invalid_dto", "DTO type not supported for creation")
}

// ApplyUpdateDTO implementa ResponseMapper
func (m *UnifiedMapper) ApplyUpdateDTO(entity interface{}, dtoParam interface{}, userCtx *common.UserContext) (interface{}, error) {
	switch e := entity.(type) {
	case *models.Event:
		if updateReq, ok := dtoParam.(*dto.UpdateEventRequest); ok {
			err := m.eventMapper.UpdateEventRequestToModel(e, updateReq, userCtx)
			return e, err
		}
	case *models.Organization:
		if updateReq, ok := dtoParam.(*dto.UpdateOrganizationRequest); ok {
			err := m.orgMapper.UpdateOrganizationRequestToModel(e, updateReq, userCtx)
			return e, err
		}
	case *models.User:
		if updateReq, ok := dtoParam.(*dto.UpdateUserRequest); ok {
			err := m.userMapper.UpdateUserRequestToModel(e, updateReq, userCtx)
			return e, err
		}
	}
	return nil, common.NewBusinessError("invalid_dto", "Update DTO type not supported")
}

// EntityToResponse implementa ResponseMapper
func (m *UnifiedMapper) EntityToResponse(entity interface{}, userCtx *common.UserContext) (interface{}, error) {
	switch e := entity.(type) {
	case *models.Event:
		return m.eventMapper.EventToResponse(e, userCtx), nil
	case models.Event:
		return m.eventMapper.EventToResponse(&e, userCtx), nil
	case *models.Organization:
		return m.orgMapper.OrganizationToResponse(e, userCtx), nil
	case models.Organization:
		return m.orgMapper.OrganizationToResponse(&e, userCtx), nil
	case *models.User:
		return m.userMapper.UserToResponse(e, userCtx), nil
	case models.User:
		return m.userMapper.UserToResponse(&e, userCtx), nil
	default:
		return nil, common.NewBusinessError("invalid_entity", "Entity type not supported for response mapping")
	}
}

// =============================================================================
// IMPLEMENTACIÓN DE EventMapper
// =============================================================================

func (m *UnifiedMapper) CreateEventRequestToModel(req *dto.CreateEventRequest, userCtx *common.UserContext) (*models.Event, error) {
	return m.eventMapper.CreateEventRequestToModel(req, userCtx)
}

func (m *UnifiedMapper) UpdateEventRequestToModel(event *models.Event, req *dto.UpdateEventRequest, userCtx *common.UserContext) error {
	return m.eventMapper.UpdateEventRequestToModel(event, req, userCtx)
}

func (m *UnifiedMapper) EventToResponse(event *models.Event, userCtx *common.UserContext) dto.EventResponse {
	return m.eventMapper.EventToResponse(event, userCtx)
}

func (m *UnifiedMapper) EventToDetailResponse(event *models.Event, userCtx *common.UserContext) dto.EventDetailResponse {
	return m.eventMapper.EventToDetailResponse(event, userCtx)
}

func (m *UnifiedMapper) EventToSummaryResponse(event *models.Event) dto.EventSummaryResponse {
	return m.eventMapper.EventToSummaryResponse(event)
}

func (m *UnifiedMapper) EventsToListResponse(events []*models.Event, pagination *common.PaginationMeta, userCtx *common.UserContext) dto.EventListResponse {
	return m.eventMapper.EventsToListResponse(events, pagination, userCtx)
}

// =============================================================================
// IMPLEMENTACIÓN DE OrganizationMapper
// =============================================================================

func (m *UnifiedMapper) CreateOrganizationRequestToModel(req *dto.CreateOrganizationRequest, userCtx *common.UserContext) (*models.Organization, error) {
	return m.orgMapper.CreateOrganizationRequestToModel(req, userCtx)
}

func (m *UnifiedMapper) UpdateOrganizationRequestToModel(org *models.Organization, req *dto.UpdateOrganizationRequest, userCtx *common.UserContext) error {
	return m.orgMapper.UpdateOrganizationRequestToModel(org, req, userCtx)
}

func (m *UnifiedMapper) OrganizationToResponse(org *models.Organization, userCtx *common.UserContext) dto.OrganizationResponse {
	return m.orgMapper.OrganizationToResponse(org, userCtx)
}

func (m *UnifiedMapper) OrganizationToDetailResponse(org *models.Organization, userCtx *common.UserContext) dto.OrganizationDetailResponse {
	return m.orgMapper.OrganizationToDetailResponse(org, userCtx)
}

func (m *UnifiedMapper) OrganizationToSummaryResponse(org *models.Organization) dto.OrganizationSummaryResponse {
	return m.orgMapper.OrganizationToSummaryResponse(org)
}

func (m *UnifiedMapper) OrganizationsToListResponse(orgs []*models.Organization, pagination *common.PaginationMeta, userCtx *common.UserContext) dto.OrganizationListResponse {
	return m.orgMapper.OrganizationsToListResponse(orgs, pagination, userCtx)
}

// =============================================================================
// IMPLEMENTACIÓN DE UserMapper
// =============================================================================

func (m *UnifiedMapper) CreateUserRequestToModel(req *dto.CreateUserRequest, userCtx *common.UserContext) (*models.User, error) {
	return m.userMapper.CreateUserRequestToModel(req, userCtx)
}

func (m *UnifiedMapper) UpdateUserRequestToModel(user *models.User, req *dto.UpdateUserRequest, userCtx *common.UserContext) error {
	return m.userMapper.UpdateUserRequestToModel(user, req, userCtx)
}

func (m *UnifiedMapper) UserToResponse(user *models.User, userCtx *common.UserContext) dto.UserResponse {
	return m.userMapper.UserToResponse(user, userCtx)
}

func (m *UnifiedMapper) UserToDetailResponse(user *models.User, userCtx *common.UserContext) dto.UserDetailResponse {
	return m.userMapper.UserToDetailResponse(user, userCtx)
}

func (m *UnifiedMapper) UserToSummaryResponse(user *models.User) dto.UserSummaryResponse {
	return m.userMapper.UserToSummaryResponse(user)
}

func (m *UnifiedMapper) UserToProfileResponse(user *models.User) dto.UserProfileResponse {
	return m.userMapper.UserToProfileResponse(user)
}

func (m *UnifiedMapper) UsersToListResponse(users []*models.User, pagination *common.PaginationMeta, userCtx *common.UserContext) dto.UserListResponse {
	return m.userMapper.UsersToListResponse(users, pagination, userCtx)
}

// =============================================================================
// IMPLEMENTACIÓN DE AuthMapper
// =============================================================================

func (m *UnifiedMapper) RegisterRequestToUser(req *dto.RegisterRequest) (*models.User, error) {
	return m.authMapper.RegisterRequestToUser(req)
}

func (m *UnifiedMapper) UserToAuthResponse(user *models.User, accessToken, refreshToken string, expiresIn int) dto.AuthResponse {
	return m.authMapper.UserToAuthResponse(user, accessToken, refreshToken, expiresIn)
}

func (m *UnifiedMapper) UserToRegisterResponse(user *models.User, accessToken, refreshToken string, expiresIn int) dto.RegisterResponse {
	return m.authMapper.UserToRegisterResponse(user, accessToken, refreshToken, expiresIn)
}

func (m *UnifiedMapper) BuildPasswordResetResponse(success bool, expiresInMinutes int) dto.PasswordResetResponse {
	return m.authMapper.BuildPasswordResetResponse(success, expiresInMinutes)
}

func (m *UnifiedMapper) BuildEmailVerificationResponse(success bool, isVerified bool) dto.EmailVerificationResponse {
	return m.authMapper.BuildEmailVerificationResponse(success, isVerified)
}

func (m *UnifiedMapper) BuildLogoutResponse(success bool) dto.LogoutResponse {
	return m.authMapper.BuildLogoutResponse(success)
}

func (m *UnifiedMapper) BuildTokenResponse(accessToken, refreshToken string, expiresAt time.Time) dto.TokenResponse {
	return m.authMapper.BuildTokenResponse(accessToken, refreshToken, expiresAt)
}

func (m *UnifiedMapper) RefreshTokensToSessionList(tokens []*models.RefreshToken, currentTokenID string) dto.SessionListResponse {
	return m.authMapper.RefreshTokensToSessionList(tokens, currentTokenID)
}
