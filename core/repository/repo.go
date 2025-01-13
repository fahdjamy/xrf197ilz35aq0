package repository

type Repositories struct {
	PermissionRepo PermissionRepository
	UserRepo       UserRepository
	OrgRepo        OrganizationRepository
	SettingsRepo   SettingsRepository
}
