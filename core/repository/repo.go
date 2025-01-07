package repository

type Repositories struct {
	RoleRepo     RoleRepository
	UserRepo     UserRepository
	OrgRepo      OrganizationRepository
	SettingsRepo SettingsRepository
}
