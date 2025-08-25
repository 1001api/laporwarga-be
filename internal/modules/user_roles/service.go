package userroles

import (
	"errors"
	db "hubku/lapor_warga_be_v2/internal/database/generated"

	"github.com/google/uuid"
)

type UserRolesService interface {
	CheckRoleExists(name string) (bool, error)
	CreateRole(arg db.CreateRoleParams) (uuid.UUID, error)
	AssignRoleToUser(arg db.AssignRoleToUserParams) error
	RemoveUserRole(userID uuid.UUID) error
	GetRoleByName(name string) (db.Role, error)
	GetRoleByID(id uuid.UUID) (db.Role, error)
	HasRole(arg db.HasRoleParams) (bool, error)
	ListAllRoles() ([]db.Role, error)
	UpdateRole(targetID uuid.UUID, req UpdateRoleRequest) error
	RemoveRole(id uuid.UUID) error
}

type service struct {
	repo UserRolesRepository
}

func NewUserRolesService(repo UserRolesRepository) UserRolesService {
	return &service{
		repo: repo,
	}
}

func (s *service) CheckRoleExists(name string) (bool, error) {
	return s.repo.CheckRoleExists(name)
}

func (s *service) CreateRole(arg db.CreateRoleParams) (uuid.UUID, error) {
	return s.repo.CreateRole(arg)
}

func (s *service) AssignRoleToUser(arg db.AssignRoleToUserParams) error {
	if arg.RoleName == "" {
		return errors.New("role name is required")
	}

	// check if rolename exist
	exists, err := s.CheckRoleExists(arg.RoleName)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("role name does not exist")
	}

	return s.repo.AssignRoleToUser(arg)
}

func (s *service) RemoveUserRole(userID uuid.UUID) error {
	return s.repo.RemoveUserRole(userID)
}

func (s *service) GetRoleByName(name string) (db.Role, error) {
	return s.repo.GetRoleByName(name)
}

func (s *service) GetRoleByID(id uuid.UUID) (db.Role, error) {
	return s.repo.GetRoleByID(id)
}

func (s *service) HasRole(arg db.HasRoleParams) (bool, error) {
	return s.repo.HasRole(arg)
}

func (s *service) ListAllRoles() ([]db.Role, error) {
	return s.repo.ListAllRoles()
}

func (s *service) UpdateRole(targetID uuid.UUID, req UpdateRoleRequest) error {
	// check if role exist
	exists, err := s.GetRoleByID(targetID)
	if err != nil {
		return err
	}

	if exists.ID == uuid.Nil {
		return errors.New("role does not exist")
	}

	// check if rolename already exist
	existsName, err := s.CheckRoleExists(req.Name)
	if err != nil {
		return err
	}

	if exists.Name != req.Name && existsName {
		return errors.New("role name already exist")
	}

	return s.repo.UpdateRole(db.UpdateRoleParams{
		ID:          targetID,
		Name:        req.Name,
		Description: req.Description,
	})
}

func (s *service) RemoveRole(id uuid.UUID) error {
	return s.repo.RemoveRole(id)
}
