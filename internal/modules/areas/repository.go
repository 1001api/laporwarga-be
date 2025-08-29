package areas

import (
	"context"
	db "hubku/lapor_warga_be_v2/internal/database/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AreaRepository interface {
	CreateArea(arg db.CreateAreaParams) (uuid.UUID, error)
	CheckAreaExist(arg db.CheckAreaExistParams) (uuid.UUID, error)
}

type repository struct {
	db *db.Queries
}

func NewAreaRepository(pool *pgxpool.Pool) AreaRepository {
	return &repository{db: db.New(pool)}
}

func (r *repository) CreateArea(arg db.CreateAreaParams) (uuid.UUID, error) {
	return r.db.CreateArea(context.Background(), arg)
}

func (r *repository) CheckAreaExist(arg db.CheckAreaExistParams) (uuid.UUID, error) {
	return r.db.CheckAreaExist(context.Background(), arg)
}
