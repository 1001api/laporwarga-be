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
	ID               uuid.UUID          `db:"id" json:"id"`
	Email            string             `db:"email" json:"email"`
	Fullname         string             `db:"fullname" json:"fullname"`
	Phone            string             `db:"phone" json:"phone"`
	Username         string             `db:"username" json:"username"`
	Role             string             `db:"role" json:"role"`
	CredibilityScore int                `db:"credibility_score" json:"credibility_score"`
	Status           string             `db:"status" json:"status"`
	IsEmailVerified  bool               `db:"is_email_verified" json:"is_email_verified"`
	IsPhoneVerified  bool               `db:"is_phone_verified" json:"is_phone_verified"`
	LastLoginAt      pgtype.Timestamptz `db:"last_login_at" json:"last_login_at"`
	CreatedAt        pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
}
