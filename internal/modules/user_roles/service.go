package userroles

import (
	"errors"
	db "hubku/lapor_warga_be_v2/internal/database/generated"
	"hubku/lapor_warga_be_v2/internal/modules/auditlogs"
	"hubku/lapor_warga_be_v2/pkg"

	"github.com/google/uuid"
)

type UserRolesService interface {
	CheckRoleExists(name string) (bool, error)
	CreateRole(currentUserID uuid.UUID, arg db.CreateRoleParams) (uuid.UUID, error)
	AssignRoleToUser(currentUserID uuid.UUID, arg db.AssignRoleToUserParams) error
	RemoveUserRole(currentUserID uuid.UUID, userID uuid.UUID) error
	GetRoleByName(name string) (db.Role, error)
	GetRoleByID(id uuid.UUID) (db.Role, error)
	HasRole(arg db.HasRoleParams) (bool, error)
	ListAllRoles() ([]db.Role, error)
	UpdateRole(currentUserID uuid.UUID, targetID uuid.UUID, req UpdateRoleRequest) error
	RemoveRole(currentUserID uuid.UUID, id uuid.UUID) error
}

type service struct {
	repo       UserRolesRepository
	logService auditlogs.LogsService
}

func NewUserRolesService(repo UserRolesRepository, logService auditlogs.LogsService) UserRolesService {
	return &service{
		repo:       repo,
		logService: logService,
	}
}

func (s *service) CheckRoleExists(name string) (bool, error) {
	return s.repo.CheckRoleExists(name)
}

func (s *service) CreateRole(currentUserID uuid.UUID, arg db.CreateRoleParams) (uuid.UUID, error) {
	createdID, err := s.repo.CreateRole(arg)
	if err != nil {
		return uuid.UUID{}, err
	}

	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityRoles),
			Action:      string(pkg.LogTypeCreate),
			EntityID:    createdID,
			PerformedBy: currentUserID,
		})
	}()

	return createdID, nil
}

func (s *service) AssignRoleToUser(currentUserID uuid.UUID, arg db.AssignRoleToUserParams) error {
	if arg.RoleName == "" {
		return errors.New("role name is required")
	}

	// check if rolename exist
	exists, err := s.GetRoleByName(arg.RoleName)
	if err != nil {
		return err
	}

	if exists.ID == uuid.Nil {
		return errors.New("role name does not exist")
	}

	err = s.repo.AssignRoleToUser(arg)
	if err != nil {
		return err
	}

	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityRoles),
			Action:      string(pkg.LogTypeAssign),
			EntityID:    exists.ID,
			PerformedBy: currentUserID,
		})
	}()

	return nil
}

func (s *service) RemoveUserRole(currentUserID uuid.UUID, userID uuid.UUID) error {
	err := s.repo.RemoveUserRole(userID)
	if err != nil {
		return err
	}

	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityRoles),
			Action:      string(pkg.LogTypeDelete),
			EntityID:    userID,
			PerformedBy: currentUserID,
		})
	}()

	return nil
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

func (s *service) UpdateRole(currentUserID uuid.UUID, targetID uuid.UUID, req UpdateRoleRequest) error {
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

	if err = s.repo.UpdateRole(db.UpdateRoleParams{
		ID:          targetID,
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		return err
	}

	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityRoles),
			Action:      string(pkg.LogTypeUpdate),
			EntityID:    exists.ID,
			PerformedBy: currentUserID,
		})
	}()

	return nil
}

func (s *service) RemoveRole(currentUserID uuid.UUID, id uuid.UUID) error {
	exists, err := s.GetRoleByID(id)
	if err != nil {
		return err
	}

	if exists.ID == uuid.Nil {
		return errors.New("role does not exist")
	}

	if err = s.repo.RemoveRole(id); err != nil {
		return err
	}

	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityRoles),
			Action:      string(pkg.LogTypeDelete),
			EntityID:    exists.ID,
			PerformedBy: currentUserID,
		})
	}()

	return nil
}
