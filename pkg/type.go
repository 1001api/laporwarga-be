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
)
