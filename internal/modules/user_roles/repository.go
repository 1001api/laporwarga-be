package userroles

import (
	"context"

	db "hubku/lapor_warga_be_v2/internal/database/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRolesRepository interface {
	CheckRoleExists(name string) (bool, error)
	CreateRole(arg db.CreateRoleParams) (uuid.UUID, error)
	AssignRoleToUser(arg db.AssignRoleToUserParams) error
	RemoveUserRole(userID uuid.UUID) error
	GetRoleByName(name string) (db.Role, error)
	GetRoleByID(id uuid.UUID) (db.Role, error)
	HasRole(arg db.HasRoleParams) (bool, error)
	ListAllRoles() ([]db.Role, error)
	UpdateRole(arg db.UpdateRoleParams) error
	RemoveRole(id uuid.UUID) error
}

type repository struct {
	queries *db.Queries
}

func NewUserRolesRepository(pool *pgxpool.Pool) UserRolesRepository {
	return &repository{
		queries: db.New(pool),
	}
}

func (r *repository) CheckRoleExists(name string) (bool, error) {
	return r.queries.CheckRoleExists(context.Background(), name)
}

func (r *repository) CreateRole(arg db.CreateRoleParams) (uuid.UUID, error) {
	return r.queries.CreateRole(context.Background(), arg)
}

func (r *repository) AssignRoleToUser(arg db.AssignRoleToUserParams) error {
	return r.queries.AssignRoleToUser(context.Background(), arg)
}

func (r *repository) RemoveUserRole(userID uuid.UUID) error {
	return r.queries.RemoveUserRole(context.Background(), userID)
}

func (r *repository) GetRoleByName(name string) (db.Role, error) {
	return r.queries.GetRoleByName(context.Background(), name)
}

func (r *repository) GetRoleByID(id uuid.UUID) (db.Role, error) {
	return r.queries.GetRoleByID(context.Background(), id)
}

func (r *repository) HasRole(arg db.HasRoleParams) (bool, error) {
	return r.queries.HasRole(context.Background(), arg)
}

func (r *repository) ListAllRoles() ([]db.Role, error) {
	return r.queries.ListAllRoles(context.Background())
}

func (r *repository) UpdateRole(arg db.UpdateRoleParams) error {
	return r.queries.UpdateRole(context.Background(), arg)
}

func (r *repository) RemoveRole(id uuid.UUID) error {
	return r.queries.DeleteRole(context.Background(), id)
}
