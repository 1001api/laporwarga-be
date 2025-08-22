package userroles

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=20"`
	Description string `json:"description"`
}
