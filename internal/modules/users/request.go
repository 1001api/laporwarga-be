package users

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserRequest struct {
	Username     string `json:"username" form:"username"`
	Email        string `json:"email" form:"email"`
	FullName     string `json:"full_name" form:"full_name"`
	PasswordHash string `json:"-"`
	PhoneNumber  string `json:"phone_number" form:"phone_number"`
	Role         string `json:"role" form:"role"`
}

type UpdateUserRequest struct {
	Email       string `json:"email" form:"email"`
	Username    string `json:"username" form:"username"`
	FullName    string `json:"full_name" form:"full_name"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Status      string `json:"status" form:"status"`
}

type UserProfileResponse struct {
	ID               uuid.UUID          `json:"id"`
	Email            string             `json:"email"`
	Fullname         string             `json:"fullname"`
	Phone            string             `json:"phone"`
	Username         string             `json:"username"`
	Role             string             `json:"role"`
	CredibilityScore int                `json:"credibility_score"`
	Status           string             `json:"status"`
	IsEmailVerified  bool               `json:"is_email_verified"`
	IsPhoneVerified  bool               `json:"is_phone_verified"`
	LastLoginAt      pgtype.Timestamptz `json:"last_login_at"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
}
