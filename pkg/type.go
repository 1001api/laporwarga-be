package pkg

type JSONB map[string]interface{}

type RoleType string
type LogType string
type AreaTolerance string
type AreaToleranceValue float64

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
	LogEntityAreas LogType = "areas"

	// JWT
	AccessTokenName  = "__hk_asid"
	RefreshTokenName = "__hk_rsid"
	MobileKeyName    = "X-Request-Tag"

	// Area Tolerance
	AreaOff             AreaTolerance      = "off"
	AreaSimple          AreaTolerance      = "simple"
	AreaDetail          AreaTolerance      = "detail"
	SimpleAreaTolerance AreaToleranceValue = 0.001
	DetailAreaTolerance AreaToleranceValue = 0.0001
	OffAreaTolerance    AreaToleranceValue = -99
)
