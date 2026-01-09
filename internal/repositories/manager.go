package repositories

// RepositoryManager centraliza todos los repositorios
type RepositoryManager struct {
	Events        *EventRepository
	Organizations *OrganizationRepository
	Users         *UserRepository
	RefreshTokens *RefreshTokenRepository
}

// NewRepositoryManager crea una nueva instancia del manager
func NewRepositoryManager() *RepositoryManager {
	return &RepositoryManager{
		Events:        NewEventRepository(),
		Organizations: NewOrganizationRepository(),
		Users:         NewUserRepository(),
		RefreshTokens: NewRefreshTokenRepository(),
	}
}
