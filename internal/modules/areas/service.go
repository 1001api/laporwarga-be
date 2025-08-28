package areas

import (
	"encoding/json"
	"errors"
	db "hubku/lapor_warga_be_v2/internal/database/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AreaService interface {
	CreateArea(arg CreateAreaRequest) (uuid.UUID, error)
}

type service struct {
	repo AreaRepository
}

func NewAreaService(repo AreaRepository) AreaService {
	return &service{repo: repo}
}

func (s *service) CreateArea(arg CreateAreaRequest) (uuid.UUID, error) {
	// validate geojson
	var geojson map[string]interface{}

	if err := json.Unmarshal([]byte(arg.GeoJSON), &geojson); err != nil {
		return uuid.Nil, err
	}

	if geojson["type"] != "Polygon" && (geojson["type"] != "Feature" || geojson["geometry"].(map[string]interface{})["type"] != "Polygon") {
		return uuid.Nil, errors.New("geojson type must be Polygon")
	}

	return s.repo.CreateArea(db.CreateAreaParams{
		Name: arg.Name,
		Description: pgtype.Text{
			String: arg.Description,
			Valid:  arg.Description != "",
		},
		AreaType: arg.AreaType,
		AreaCode: arg.AreaCode,
		Boundary: arg.GeoJSON,
	})
}
