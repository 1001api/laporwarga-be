package users

import (
	"errors"
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
	DeleteUser(id uuid.UUID, deletedBy uuid.UUID) error
	RestoreUser(id uuid.UUID) error
	SearchUser(query string, page, limit int32) ([]UserProfileResponse, error)
	GetUserByIdentifier(identifier string) (db.GetUserByIdentifierRow, error)
	GetUserByID(id uuid.UUID) (UserProfileResponse, error)
}

type service struct {
	enckey  []byte
	repo    UserRepository
	roleSvc userroles.UserRolesService
}

func NewUserService(repo UserRepository, roleSvc userroles.UserRolesService, encKey string) UserService {
	return &service{
		repo:    repo,
		roleSvc: roleSvc,
		enckey:  []byte(encKey),
	}
}

func (s *service) decryptFields(encEmail, encFullname, encPhone []byte) ([]byte, []byte, []byte, error) {
	email, err := pkg.Decrypt(encEmail, s.enckey)
	if err != nil {
		return nil, nil, nil, err
	}

	fullname, err := pkg.Decrypt(encFullname, s.enckey)
	if err != nil {
		return nil, nil, nil, err
	}

	phone, err := pkg.Decrypt(encPhone, s.enckey)
	if err != nil {
		return nil, nil, nil, err
	}

	return email, fullname, phone, nil
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
		de, df, dp, err := s.decryptFields(user.Email, user.Fullname, user.Phone)
		if err != nil {
			return nil, err
		}

		res = append(res, UserProfileResponse{
			ID:               user.ID,
			Username:         user.Username,
			Email:            string(de),
			Fullname:         string(df),
			Phone:            string(dp),
			Role:             user.RoleName.String,
			CredibilityScore: int(user.CredibilityScore.Int16),
			Status:           user.Status.String,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
			IsEmailVerified:  user.IsEmailVerified.Bool,
			IsPhoneVerified:  user.IsPhoneVerified.Bool,
			LastLoginAt:      user.LastLoginAt,
		})
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

	var roleID uuid.UUID

	// check if role exists
	roleExist, err := s.roleSvc.GetRoleByName(string(pkg.RoleAdmin))
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return uuid.UUID{}, err
	}

	if roleExist.ID == uuid.Nil {
		// if role does not exist, create it first
		createdID, err := s.roleSvc.CreateRole(db.CreateRoleParams{
			Name:        string(pkg.RoleAdmin),
			Description: "Default admin role",
		})
		if err != nil {
			return uuid.UUID{}, err
		}

		roleID = createdID
	} else {
		roleID = roleExist.ID
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
		RoleID:       roleID,
	})
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (s *service) UpdateUser(targetID uuid.UUID, updatedBy uuid.UUID, req UpdateUserRequest) error {
	current, err := s.GetUserByID(targetID)
	if err != nil {
		log.Println("Failed to get user for update:", err)
		return err
	}

	exists, err := s.CheckUserExists(req.Email, req.Username)
	if err != nil {
		log.Println("Failed to check user exists:", err)
		return err
	}

	if exists && (req.Email != current.Email || req.Username != current.Username) {
		return errors.New("username or email already exists")
	}

	fullnameHash := pkg.HashValue(req.FullName)
	emailHash := pkg.HashValue(req.Email)
	phoneHash := pkg.HashValue(req.PhoneNumber)

	emailEnc, err := pkg.Encrypt([]byte(req.Email), s.enckey)
	if err != nil {
		return err
	}

	fullnameEnc, err := pkg.Encrypt([]byte(req.FullName), s.enckey)
	if err != nil {
		return err
	}
	phoneEnc, err := pkg.Encrypt([]byte(req.PhoneNumber), s.enckey)
	if err != nil {
		return err
	}

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

func (s *service) DeleteUser(id uuid.UUID, deletedBy uuid.UUID) error {
	return s.repo.DeleteUser(db.DeleteUserParams{
		ID:        id,
		DeletedBy: deletedBy,
	})
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
		de, df, dp, err := s.decryptFields(user.Email, user.Fullname, user.Phone)
		if err != nil {
			return nil, err
		}

		res = append(res, UserProfileResponse{
			ID:               user.ID,
			Username:         user.Username,
			Email:            string(de),
			Fullname:         string(df),
			Phone:            string(dp),
			Role:             user.RoleName.String,
			CredibilityScore: int(user.CredibilityScore.Int16),
			Status:           user.Status.String,
			CreatedAt:        user.CreatedAt,
			UpdatedAt:        user.UpdatedAt,
			IsEmailVerified:  user.IsEmailVerified.Bool,
			IsPhoneVerified:  user.IsPhoneVerified.Bool,
			LastLoginAt:      user.LastLoginAt,
		})
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

	de, df, dp, err := s.decryptFields(user.Email, user.Fullname, user.Phone)
	if err != nil {
		return db.GetUserByIdentifierRow{}, err
	}

	user.Email = de
	user.Fullname = df
	user.Phone = dp

	return user, nil
}

func (s *service) GetUserByID(id uuid.UUID) (UserProfileResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return UserProfileResponse{}, err
	}

	de, df, dp, err := s.decryptFields(user.Email, user.Fullname, user.Phone)
	if err != nil {
		return UserProfileResponse{}, err
	}

	return UserProfileResponse{
		ID:               user.ID,
		Username:         user.Username,
		Email:            string(de),
		Fullname:         string(df),
		Phone:            string(dp),
		Role:             user.RoleName.String,
		CredibilityScore: int(user.CredibilityScore.Int16),
		Status:           user.Status.String,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		IsEmailVerified:  user.IsEmailVerified.Bool,
		IsPhoneVerified:  user.IsPhoneVerified.Bool,
		LastLoginAt:      user.LastLoginAt,
	}, nil
}
