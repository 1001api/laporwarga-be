package users

import (
	"context"
	db "hubku/lapor_warga_be_v2/internal/database/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	IncrementFailedLogins(id uuid.UUID) error
	CheckUserExists(arg db.CheckUserExistsParams) (bool, error)
	CreateUser(req db.CreateUserParams) (uuid.UUID, error)
	GetUsers(arg db.GetUsersParams) ([]db.GetUsersRow, error)
	UpdateUser(arg db.UpdateUserParams) error
	DeleteUser(id uuid.UUID) error
	RestoreUser(id uuid.UUID) error
	UpdateLastLogin(id uuid.UUID) error
	SearchUser(arg db.SearchUserParams) ([]db.SearchUserRow, error)
	GetUserByIdentifier(arg db.GetUserByIdentifierParams) (db.GetUserByIdentifierRow, error)
	GetUserByID(id uuid.UUID) (db.GetUserByIDRow, error)
}

type repository struct {
	queries *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &repository{
		queries: db.New(pool),
	}
}

func (r *repository) GetUsers(arg db.GetUsersParams) ([]db.GetUsersRow, error) {
	return r.queries.GetUsers(context.Background(), arg)
}

func (r *repository) CreateUser(arg db.CreateUserParams) (uuid.UUID, error) {
	return r.queries.CreateUser(context.Background(), arg)
}

func (r *repository) UpdateUser(arg db.UpdateUserParams) error {
	return r.queries.UpdateUser(context.Background(), arg)
}

func (r *repository) UpdateLastLogin(id uuid.UUID) error {
	return r.queries.UpdateLastLogin(context.Background(), id)
}

func (r *repository) IncrementFailedLogins(id uuid.UUID) error {
	return r.queries.IncrementFailedLoginCount(context.Background(), id)
}

func (r *repository) CheckUserExists(arg db.CheckUserExistsParams) (bool, error) {
	return r.queries.CheckUserExists(context.Background(), arg)
}

func (r *repository) DeleteUser(id uuid.UUID) error {
	return r.queries.DeleteUser(context.Background(), id)
}

func (r *repository) RestoreUser(id uuid.UUID) error {
	return r.queries.RestoreUser(context.Background(), id)
}

func (r *repository) SearchUser(arg db.SearchUserParams) ([]db.SearchUserRow, error) {
	return r.queries.SearchUser(context.Background(), arg)
}

func (r *repository) GetUserByIdentifier(arg db.GetUserByIdentifierParams) (db.GetUserByIdentifierRow, error) {
	return r.queries.GetUserByIdentifier(context.Background(), arg)
}

func (r *repository) GetUserByID(id uuid.UUID) (db.GetUserByIDRow, error) {
	return r.queries.GetUserByID(context.Background(), id)
}
