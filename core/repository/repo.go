package repository

type Repositories struct {
	UserRepo       UserRepository
	SettingsRepo   SettingsRepository
	PermissionRepo PermissionRepository
	OrgRepo        OrganizationRepository
}
