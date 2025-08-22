package pkg

type JSONB map[string]interface{}

type RoleType string

const (
	RoleCitizen  RoleType = "citizen"
	RoleOfficial RoleType = "official"
	RoleAdmin    RoleType = "admin"
)
