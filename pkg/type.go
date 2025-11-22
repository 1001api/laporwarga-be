package pkg

type JSONB map[string]interface{}

type RoleType string
type LogType string
type AreaTolerance string
type AreaToleranceValue float64
type JWTTokenType string

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
	LogEntityUsers      LogType = "users"
	LogEntityRoles      LogType = "roles"
	LogEntityAreas      LogType = "areas"
	LogEntityCategories LogType = "categories"
	LogEntityReports    LogType = "reports"

	// JWT
	AccessTokenName               = "__asid"
	RefreshTokenName              = "__rsid"
	MobileKeyName                 = "X-Request-Tag"
	RefreshToken     JWTTokenType = "refresh"
	AccessToken      JWTTokenType = "access"
	JWTIssuer                     = "lapor_warga"

	// Area Tolerance
	AreaOff             AreaTolerance      = "off"
	AreaSimple          AreaTolerance      = "simple"
	AreaDetail          AreaTolerance      = "detail"
	SimpleAreaTolerance AreaToleranceValue = 0.001
	DetailAreaTolerance AreaToleranceValue = 0.0001
	OffAreaTolerance    AreaToleranceValue = -99

	// Error
	ErrExist  = "exist"
	ErrNoRows = "no rows in result set"
)

type Meta struct {
	Duration string `json:"duration"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type ErrorResponse struct {
	Error interface{} `json:"error"`
	Meta  Meta        `json:"meta"`
}
