package constants

const (
	V1            = "v1"
	API           = "api"
	SLASH         = "/"
	DASH          = "-"
	EMPTY         = ""
	EQUALS        = "="
	UNDERSCORE    = "_"
	NAME          = "name"
	EMAIL         = "email"
	OrgId         = "orgId"
	USERID        = "userId"
	PASSWORD      = "password"
	IsAnonymous   = "isAnonymous"
	FINGERPRINT   = "fingerPrint"
	PermissionId  = "permissionId"
	RedisPort     = "XRF_Q0_REDIS_PORT"
	RedisPassword = "XRF_Q0_REDIS_PASS"
	RedisAddress  = "XRF_Q0_REDIS_ADDRESS"
)

// Reserved string values

const (
	DefaultOrgName = "XRF_ILZ_DEFAULT"
)

// Error Constants

const (
	DuplicateNameDBErr = "name already exists"
	NotFoundOrgErrMsg  = "organization not found"
)

const ContentType = "Content-Type"
const ContentTypeJson = "application/json"

// Mongo Collections

const (
	UserCollection     = "user"
	SettingsCollection = "settings"
	PermissionsCol     = "permission"
	OrgCollection      = "organization"
)

// AllCollections !IMPORTANT: make sure to always add all collection names to this list
var AllCollections = [...]string{
	OrgCollection,
	PermissionsCol,
	UserCollection,
	SettingsCollection,
}
