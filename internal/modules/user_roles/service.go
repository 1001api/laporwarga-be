package userroles

import (
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
	UpdateRole(arg db.UpdateRoleParams) error
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

func (s *service) UpdateRole(arg db.UpdateRoleParams) error {
	return s.repo.UpdateRole(arg)
}
