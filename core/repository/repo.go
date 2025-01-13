package repository

type Repositories struct {
	RoleRepo     PermissionRepository
	UserRepo     UserRepository
	OrgRepo      OrganizationRepository
	SettingsRepo SettingsRepository
}
