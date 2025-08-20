package users

type CreateUserRequest struct {
	Username    string `json:"username" form:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" form:"email" validate:"required,email"`
	Fullname    string `json:"fullname" form:"fullname" validate:"required,min=2,max=100"`
	Password    string `json:"password" form:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"omitempty,min=5,max=50"`
	Role        string `json:"role" form:"role" validate:"omitempty,oneof=citizen admin superadmin"`
}

type UpdateUserParams struct {
	Username    string `json:"username" form:"username" validate:"omitempty,min=3,max=50"`
	Email       string `json:"email" form:"email" validate:"omitempty,email"`
	Fullname    string `json:"fullname" form:"fullname" validate:"omitempty,min=2,max=100"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"omitempty,min=5,max=50"`
	Role        string `json:"role" form:"role" validate:"omitempty,oneof=citizen admin superadmin"`
}
