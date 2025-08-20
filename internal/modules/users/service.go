package users

import (
	"errors"
	db "hubku/lapor_warga_be_v2/internal/database/generated"

	"github.com/google/uuid"
)

type UserService interface {
	GetUsers(arg db.GetUsersParams) ([]db.GetUsersRow, error)
	UpdateUserLastLogin(id uuid.UUID) error
	IncrementFailedLogins(id uuid.UUID) error
	CheckUserExists(arg db.CheckUserExistsParams) (bool, error)
	CreateUser(params db.CreateUserParams) (db.CreateUserRow, error)
	UpdateUser(targetID uuid.UUID, req UpdateUserParams) error
	DeleteUser(id uuid.UUID) error
	RestoreUser(id uuid.UUID) error
	SearchUser(query string, page, limit int32) ([]db.SearchUserRow, error)
	GetUserByIdentifier(identifier string) (db.GetUserByIdentifierRow, error)
	GetUserByID(id uuid.UUID) (db.GetUserByIDRow, error)
}

type service struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &service{repo: repo}
}

func (s *service) GetUsers(arg db.GetUsersParams) ([]db.GetUsersRow, error) {
	return s.repo.GetUsers(arg)
}

func (s *service) UpdateUserLastLogin(id uuid.UUID) error {
	return s.repo.UpdateLastLogin(id)
}

func (s *service) IncrementFailedLogins(id uuid.UUID) error {
	return s.repo.IncrementFailedLogins(id)
}

func (s *service) CheckUserExists(arg db.CheckUserExistsParams) (bool, error) {
	return s.repo.CheckUserExists(arg)
}

func (s *service) CreateUser(params db.CreateUserParams) (db.CreateUserRow, error) {
	exists, err := s.CheckUserExists(db.CheckUserExistsParams{
		Email:    params.Email,
		Username: params.Username,
	})
	if err != nil {
		return db.CreateUserRow{}, err
	}

	if exists {
		return db.CreateUserRow{}, errors.New("username or email already exists")
	}

	return s.repo.CreateUser(params)
}

func (s *service) UpdateUser(targetID uuid.UUID, req UpdateUserParams) error {
	// check if user exists
	exists, err := s.CheckUserExists(db.CheckUserExistsParams{
		Email: req.Email,
	})
	if err != nil {
		return err
	}

	if exists {
		return errors.New("username or email already exists")
	}

	return s.repo.UpdateUser(db.UpdateUserParams{
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Role:        req.Role,
		ID:          targetID,
	})
}

func (s *service) DeleteUser(id uuid.UUID) error {
	return s.repo.DeleteUser(id)
}

func (s *service) RestoreUser(id uuid.UUID) error {
	return s.repo.RestoreUser(id)
}

func (s *service) SearchUser(query string, page, limit int32) ([]db.SearchUserRow, error) {
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

	return results, nil

}

func (s *service) GetUserByIdentifier(identifier string) (db.GetUserByIdentifierRow, error) {
	uuid, _ := uuid.Parse(identifier)
	return s.repo.GetUserByIdentifier(db.GetUserByIdentifierParams{
		ID:       uuid,
		Email:    identifier,
		Username: identifier,
	})
}

func (s *service) GetUserByID(id uuid.UUID) (db.GetUserByIDRow, error) {
	return s.repo.GetUserByID(id)
}
