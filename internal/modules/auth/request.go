package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username    string `json:"username" form:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" form:"email" validate:"required,email"`
	Fullname    string `json:"fullname" form:"fullname" validate:"required,min=2,max=100"`
	Password    string `json:"password" form:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number" form:"phone_number" validate:"omitempty,min=5,max=50"`
	Role        string `json:"role" form:"role" validate:"omitempty,oneof=citizen admin superadmin"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Identifier string `json:"identifier" form:"identifier"` // username or email
	Password   string `json:"password" form:"password"`
}

type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	User         *UserInfo `json:"user,omitempty"`
}

type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Fullname string    `json:"fullname"`
	Role     string    `json:"role"`
}
