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
	GetAreas(arg db.GetAreasParams) ([]db.GetAreasRow, error)
	GetAreaBoundary(id uuid.UUID) (db.GetAreaBoundaryRow, error)
	ToggleAreaActiveStatus(id uuid.UUID) (db.ToggleAreaActiveStatusRow, error)
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

func (r *repository) GetAreas(arg db.GetAreasParams) ([]db.GetAreasRow, error) {
	return r.db.GetAreas(context.Background(), arg)
}

func (r *repository) GetAreaBoundary(id uuid.UUID) (db.GetAreaBoundaryRow, error) {
	return r.db.GetAreaBoundary(context.Background(), id)
}

func (r *repository) ToggleAreaActiveStatus(id uuid.UUID) (db.ToggleAreaActiveStatusRow, error) {
	return r.db.ToggleAreaActiveStatus(context.Background(), id)
}
