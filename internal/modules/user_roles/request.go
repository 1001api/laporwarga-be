package userroles

type CreateRoleRequest struct {
	Name        string `json:"name" form:"name" validate:"required,min=3,max=20"`
	Description string `json:"description" form:"description"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name" form:"name" validate:"required,min=3,max=20"`
	Description string `json:"description" form:"description"`
}

type AssignRoleRequest struct {
	RoleName string `json:"role_name" form:"role_name" validate:"required"`
}
