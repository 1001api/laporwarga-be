package userroles

import (
	"context"

	db "hubku/lapor_warga_be_v2/internal/database/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRolesRepository interface {
	CreateUserRole(arg db.CreateUserRoleParams) error
	GetUserRoleByUserID(userID uuid.UUID) (db.UserRole, error)
	UpdateUserRole(arg db.UpdateUserRoleParams) (db.UserRole, error)
	DeleteUserRole(arg db.DeleteUserRoleParams) error
}

type repository struct {
	queries *db.Queries
}

func NewUserRolesRepository(pool *pgxpool.Pool) UserRolesRepository {
	return &repository{
		queries: db.New(pool),
	}
}

func (r *repository) CreateUserRole(arg db.CreateUserRoleParams) error {
	return r.queries.CreateUserRole(context.Background(), arg)
}

func (r *repository) GetUserRoleByUserID(userID uuid.UUID) (db.UserRole, error) {
	return r.queries.GetUserRoleByUserID(context.Background(), userID)
}

func (r *repository) UpdateUserRole(arg db.UpdateUserRoleParams) (db.UserRole, error) {
	return r.queries.UpdateUserRole(context.Background(), arg)
}

func (r *repository) DeleteUserRole(arg db.DeleteUserRoleParams) error {
	return r.queries.DeleteUserRole(context.Background(), arg)
}
