package users

import (
	"errors"
	"fmt"
	db "hubku/lapor_warga_be_v2/internal/database/generated"
	userroles "hubku/lapor_warga_be_v2/internal/modules/user_roles"
	"hubku/lapor_warga_be_v2/pkg"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type UserService interface {
	InitializeRootUser() error
	GetUsers(arg db.GetUsersParams) ([]UserProfileResponse, error)
	UpdateUserLastLogin(id uuid.UUID) error
	IncrementFailedLogins(id uuid.UUID) error
	CheckUserExists(email, username string) (bool, error)
	CreateUser(params CreateUserRequest) (uuid.UUID, error)
	UpdateUser(targetID uuid.UUID, updatedBy uuid.UUID, req UpdateUserRequest) error
	DeleteUser(id uuid.UUID) error
	RestoreUser(id uuid.UUID) error
	SearchUser(query string, page, limit int32) ([]UserProfileResponse, error)
	GetUserByIdentifier(identifier string) (db.GetUserByIdentifierRow, error)
	GetUserByID(id uuid.UUID) (UserProfileResponse, error)
}

type service struct {
	enckey   []byte
	repo     UserRepository
	roleRepo userroles.UserRolesRepository
}

func NewUserService(repo UserRepository, roleRepo userroles.UserRolesRepository, encKey string) UserService {
	return &service{
		repo:     repo,
		roleRepo: roleRepo,
		enckey:   []byte(encKey),
	}
}

func (s *service) InitializeRootUser() error {
	username := viper.GetString("ROOT_USERNAME")
	passwordHash, _ := pkg.HashPassword(viper.GetString("ROOT_PASSWORD"))
	email := viper.GetString("ROOT_EMAIL")
	fullname := viper.GetString("ROOT_FULLNAME")
	phoneNumber := viper.GetString("ROOT_PHONE")

	// check if root user exists
	exists, err := s.CheckUserExists(email, username)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		log.Println("Failed to check root user:", err)
		return err
	}

	// skip if exist
	if exists {
		return nil
	}

	if username == "" || passwordHash == "" || email == "" || fullname == "" || phoneNumber == "" {
		return errors.New("ROOT_USERNAME, ROOT_PASSWORD, ROOT_EMAIL, ROOT_FULLNAME, ROOT_PHONE_NUMBER are required")
	}

	_, err = s.CreateUser(CreateUserRequest{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
		FullName:     fullname,
		PhoneNumber:  phoneNumber,
		Role:         string(pkg.RoleAdmin),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetUsers(arg db.GetUsersParams) ([]UserProfileResponse, error) {
	result, err := s.repo.GetUsers(arg)
	if err != nil {
		return nil, err
	}

	var res []UserProfileResponse

	for _, user := range result {
		decryptedEmail, err := pkg.Decrypt(user.Email, s.enckey)
		if err != nil {
			return nil, err
		}

		decryptedFullName, err := pkg.Decrypt(user.Fullname, s.enckey)
		if err != nil {
			return nil, err
		}

		decryptedPhone, err := pkg.Decrypt(user.Phone, s.enckey)
		if err != nil {
			return nil, err
		}

		user := UserProfileResponse{
			ID:               user.ID,
			Username:         user.Username,
			Email:            string(decryptedEmail),
			Fullname:         string(decryptedFullName),
			Phone:            string(decryptedPhone),
			Role:             user.Role.String,
			CredibilityScore: int(user.CredibilityScore.Int16),
			Status:           user.Status.String,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
			IsEmailVerified:  user.IsEmailVerified.Bool,
			IsPhoneVerified:  user.IsPhoneVerified.Bool,
			LastLoginAt:      user.LastLoginAt,
		}

		res = append(res, user)
	}

	return res, nil
}

func (s *service) UpdateUserLastLogin(id uuid.UUID) error {
	return s.repo.UpdateLastLogin(id)
}

func (s *service) IncrementFailedLogins(id uuid.UUID) error {
	return s.repo.IncrementFailedLogins(id)
}

func (s *service) CheckUserExists(email, username string) (bool, error) {
	hashedEmail := pkg.HashValue(email)

	return s.repo.CheckUserExists(db.CheckUserExistsParams{
		EmailHash: hashedEmail,
		Username:  username,
	})
}

func (s *service) CreateUser(params CreateUserRequest) (uuid.UUID, error) {
	exists, err := s.CheckUserExists(params.Email, params.Username)
	if err != nil {
		return uuid.UUID{}, err
	}

	if exists {
		return uuid.UUID{}, errors.New("username or email already exists")
	}

	fullnameHash := pkg.HashValue(params.FullName)
	emailHash := pkg.HashValue(params.Email)
	phoneHash := pkg.HashValue(params.PhoneNumber)
	emailEnc, err := pkg.Encrypt([]byte(params.Email), s.enckey)
	if err != nil {
		return uuid.UUID{}, err
	}
	fullnameEnc, err := pkg.Encrypt([]byte(params.FullName), s.enckey)
	if err != nil {
		return uuid.UUID{}, err
	}
	phoneEnc, err := pkg.Encrypt([]byte(params.PhoneNumber), s.enckey)
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := s.repo.CreateUser(db.CreateUserParams{
		EmailHash:    emailHash,
		EmailEnc:     emailEnc,
		FullnameHash: fullnameHash,
		FullnameEnc:  fullnameEnc,
		Username:     params.Username,
		PasswordHash: params.PasswordHash,
		PhoneHash:    phoneHash,
		PhoneEnc:     phoneEnc,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	fmt.Println("created user:", userID)

	if err := s.roleRepo.CreateUserRole(db.CreateUserRoleParams{
		UserID:    userID,
		RoleType:  params.Role,
		CreatedBy: userID,
	}); err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (s *service) UpdateUser(targetID uuid.UUID, updatedBy uuid.UUID, req UpdateUserRequest) error {
	exists, err := s.CheckUserExists(req.Email, req.Username)
	if err != nil {
		log.Println("Failed to check user exists:", err)
		return err
	}

	fmt.Println("exists:", exists, req.PhoneNumber)

	if exists {
		return errors.New("username or email already exists")
	}

	fullnameHash := pkg.HashValue(req.FullName)
	emailHash := pkg.HashValue(req.Email)
	phoneHash := pkg.HashValue(req.PhoneNumber)
	emailEnc, _ := pkg.Encrypt([]byte(req.Email), s.enckey)
	fullnameEnc, _ := pkg.Encrypt([]byte(req.FullName), s.enckey)
	phoneEnc, _ := pkg.Encrypt([]byte(req.PhoneNumber), s.enckey)

	return s.repo.UpdateUser(db.UpdateUserParams{
		Username:     req.Username,
		EmailHash:    emailHash,
		EmailEnc:     emailEnc,
		FullnameHash: fullnameHash,
		FullnameEnc:  fullnameEnc,
		PhoneHash:    phoneHash,
		PhoneEnc:     phoneEnc,
		Status:       string(req.Status),
		ID:           targetID,
		UpdatedBy:    updatedBy,
	})
}

func (s *service) DeleteUser(id uuid.UUID) error {
	return s.repo.DeleteUser(id)
}

func (s *service) RestoreUser(id uuid.UUID) error {
	return s.repo.RestoreUser(id)
}

func (s *service) SearchUser(query string, page, limit int32) ([]UserProfileResponse, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	results, err := s.repo.SearchUser(db.SearchUserParams{
		Query:       query,
		OffsetCount: (page - 1) * limit,
		LimitCount:  limit,
	})
	if err != nil {
		return nil, err
	}

	var res []UserProfileResponse

	for _, user := range results {
		decryptedEmail, err := pkg.Decrypt(user.Email, s.enckey)
		if err != nil {
			return nil, err
		}

		decryptedFullName, err := pkg.Decrypt(user.Fullname, s.enckey)
		if err != nil {
			return nil, err
		}

		decryptedPhone, err := pkg.Decrypt(user.Phone, s.enckey)
		if err != nil {
			return nil, err
		}

		user := UserProfileResponse{
			ID:               user.ID,
			Username:         user.Username,
			Email:            string(decryptedEmail),
			Fullname:         string(decryptedFullName),
			Phone:            string(decryptedPhone),
			Role:             user.Role.String,
			CredibilityScore: int(user.CredibilityScore.Int16),
			Status:           user.Status.String,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
			IsEmailVerified:  user.IsEmailVerified.Bool,
			IsPhoneVerified:  user.IsPhoneVerified.Bool,
			LastLoginAt:      user.LastLoginAt,
		}

		res = append(res, user)
	}

	return res, nil

}

func (s *service) GetUserByIdentifier(identifier string) (db.GetUserByIdentifierRow, error) {
	parsedUUID, _ := uuid.Parse(identifier)

	user, err := s.repo.GetUserByIdentifier(db.GetUserByIdentifierParams{
		ID:        parsedUUID,
		EmailHash: pkg.HashValue(identifier),
		Username:  identifier,
	})

	if err != nil {
		return db.GetUserByIdentifierRow{}, err
	}

	decryptedEmail, err := pkg.Decrypt(user.Email, s.enckey)
	if err != nil {
		return db.GetUserByIdentifierRow{}, err
	}
	user.Email = decryptedEmail

	decryptedFullName, err := pkg.Decrypt(user.Fullname, s.enckey)
	if err != nil {
		return db.GetUserByIdentifierRow{}, err
	}
	user.Fullname = decryptedFullName

	decryptedPhone, err := pkg.Decrypt(user.Phone, s.enckey)
	if err != nil {
		return db.GetUserByIdentifierRow{}, err
	}
	user.Phone = decryptedPhone

	return user, nil
}

func (s *service) GetUserByID(id uuid.UUID) (UserProfileResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return UserProfileResponse{}, err
	}

	decryptedEmail, err := pkg.Decrypt(user.Email, s.enckey)
	if err != nil {
		return UserProfileResponse{}, err
	}
	user.Email = decryptedEmail

	decryptedFullName, err := pkg.Decrypt(user.Fullname, s.enckey)
	if err != nil {
		return UserProfileResponse{}, err
	}
	user.Fullname = decryptedFullName

	decryptedPhone, err := pkg.Decrypt(user.Phone, s.enckey)
	if err != nil {
		return UserProfileResponse{}, err
	}
	user.Phone = decryptedPhone

	return UserProfileResponse{
		ID:               user.ID,
		Username:         user.Username,
		Email:            string(decryptedEmail),
		Fullname:         string(decryptedFullName),
		Phone:            string(decryptedPhone),
		Role:             user.Role.String,
		CredibilityScore: int(user.CredibilityScore.Int16),
		Status:           user.Status.String,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		IsEmailVerified:  user.IsEmailVerified.Bool,
		IsPhoneVerified:  user.IsPhoneVerified.Bool,
		LastLoginAt:      user.LastLoginAt,
	}, nil
}
