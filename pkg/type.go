package pkg

type JSONB map[string]interface{}

type RoleType string
type LogType string

const (
	RoleCitizen  RoleType = "citizen"
	RoleOfficial RoleType = "official"
	RoleAdmin    RoleType = "admin"

	// Log Type
	LogTypeLogin   LogType = "login"
	LogTypeCreate  LogType = "create"
	LogTypeUpdate  LogType = "update"
	LogTypeDelete  LogType = "delete"
	LogTypeAssign  LogType = "assign"
	LogTypeRestore LogType = "restore"

	// Log Entiry
	LogEntityUsers LogType = "users"
	LogEntityRoles LogType = "roles"

	// JWT
	AccessTokenName  = "__hk_asid"
	RefreshTokenName = "__hk_rsid"
	MobileKeyName    = "X-Request-Tag"
)
