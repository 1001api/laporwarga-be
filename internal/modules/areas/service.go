package areas

import (
	"encoding/json"
	"errors"
	"fmt"
	db "hubku/lapor_warga_be_v2/internal/database/generated"
	"hubku/lapor_warga_be_v2/internal/modules/auditlogs"
	"hubku/lapor_warga_be_v2/pkg"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type AreaService interface {
	CreateArea(currentUserID uuid.UUID, arg CreateAreaRequest) (uuid.UUID, error)
	GetAreas(page, limit int, tolerance pkg.AreaTolerance) ([]db.GetAreasRow, error)
	GetAreaBoundary(id uuid.UUID) (db.GetAreaBoundaryRow, error)
	ToggleAreaActiveStatus(currentUserID uuid.UUID, id uuid.UUID) (db.ToggleAreaActiveStatusRow, error)
}

type service struct {
	logService auditlogs.LogsService
	repo       AreaRepository
}

func NewAreaService(logService auditlogs.LogsService, repo AreaRepository) AreaService {
	return &service{logService: logService, repo: repo}
}

func (s *service) CreateArea(currentUserID uuid.UUID, arg CreateAreaRequest) (uuid.UUID, error) {
	// Check if area already exist
	id, err := s.repo.CheckAreaExist(db.CheckAreaExistParams{
		Name:     arg.Name,
		AreaCode: arg.AreaCode,
	})
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		return uuid.Nil, err
	}

	if id != uuid.Nil {
		return uuid.Nil, errors.New("area already exist")
	}

	// Parse GeoJSON as a single Feature
	fc, err := geojson.UnmarshalFeatureCollection([]byte(arg.GeoJSON))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse GeoJSON: %w", err)
	}

	if len(fc.Features) == 0 {
		return uuid.Nil, errors.New("geojson must have at least one feature")
	}

	// take the first and ignore the rest
	feature := fc.Features[0]

	switch geom := feature.Geometry.(type) {
	case orb.Polygon:
		if len(geom) == 0 {
			return uuid.Nil, errors.New("polygon has no coordinates")
		}
	case orb.MultiPolygon:
		if len(geom) == 0 {
			return uuid.Nil, errors.New("multipolygon has no coordinates")
		}
	default:
		return uuid.Nil, errors.New("only Polygon and Multipolygon geojson formats are allowed")
	}

	// ring validation
	if err := validateGeometry(feature.Geometry); err != nil {
		return uuid.Nil, err
	}

	// Marshal the 2D geometry to GeoJSON (orb already dropped Z)
	geomBytes, err := marshalGeometryToJSON(feature.Geometry)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal geometry: %w", err)
	}

	result, err := s.repo.CreateArea(db.CreateAreaParams{
		Name: arg.Name,
		Description: pgtype.Text{
			String: arg.Description,
			Valid:  arg.Description != "",
		},
		AreaType: arg.AreaType,
		AreaCode: arg.AreaCode,
		Boundary: string(geomBytes),
	})

	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityAreas),
			Action:      string(pkg.LogTypeCreate),
			EntityID:    result,
			PerformedBy: currentUserID,
		})
	}()

	return result, err
}

// marshalGeometryToJSON converts orb.Geometry to 2D GeoJSON bytes efficiently
func marshalGeometryToJSON(geom orb.Geometry) ([]byte, error) {
	typeStr := geom.GeoJSONType()
	coords := getCoordinates(geom)

	return json.Marshal(struct {
		Type        string      `json:"type"`
		Coordinates interface{} `json:"coordinates"`
	}{
		Type:        typeStr,
		Coordinates: coords,
	})
}

// getCoordinates builds nested coordinate slices from orb.Geometry
func getCoordinates(geom orb.Geometry) interface{} {
	switch g := geom.(type) {
	case orb.Polygon:
		coords := make([][][]float64, len(g))
		for i, ring := range g {
			coords[i] = make([][]float64, len(ring))
			for j, p := range ring {
				coords[i][j] = []float64{p[0], p[1]}
			}
		}
		return coords
	case orb.MultiPolygon:
		coords := make([][][][]float64, len(g))
		for i, poly := range g {
			coords[i] = make([][][]float64, len(poly))
			for k, ring := range poly {
				coords[i][k] = make([][]float64, len(ring))
				for j, p := range ring {
					coords[i][k][j] = []float64{p[0], p[1]}
				}
			}
		}
		return coords
	}

	// Should not reach here due to validation
	return nil
}

// validateGeometry performs sanitization checks
func validateGeometry(geom orb.Geometry) error {
	switch g := geom.(type) {
	case orb.Polygon:
		for _, ring := range g {
			if len(ring) < 4 || !ring.Closed() {
				return errors.New("invalid polygon: must have at least 4 points and be closed")
			}
		}
	case orb.MultiPolygon:
		for _, poly := range g {
			for _, ring := range poly {
				if len(ring) < 4 || !ring.Closed() {
					return errors.New("invalid multipolygon: must have at least 4 points and be closed")
				}
			}
		}
	default:
		return errors.New("unsupported geometry type")
	}
	return nil
}

func (s *service) GetAreas(page, limit int, tolerance pkg.AreaTolerance) ([]db.GetAreasRow, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 20
	}

	var simplifyTolerance pkg.AreaToleranceValue

	switch tolerance {
	case pkg.AreaSimple:
		simplifyTolerance = pkg.SimpleAreaTolerance
	case pkg.AreaDetail:
		simplifyTolerance = pkg.DetailAreaTolerance
	case pkg.AreaOff:
		simplifyTolerance = pkg.OffAreaTolerance
	default:
		return nil, errors.New("invalid tolerance")
	}

	return s.repo.GetAreas(db.GetAreasParams{
		OffsetCount:       int32((page - 1) * limit),
		LimitCount:        int32(limit),
		SimplifyTolerance: float64(simplifyTolerance),
	})
}

func (s *service) GetAreaBoundary(id uuid.UUID) (db.GetAreaBoundaryRow, error) {
	return s.repo.GetAreaBoundary(id)
}

func (s *service) ToggleAreaActiveStatus(currentUserID uuid.UUID, id uuid.UUID) (db.ToggleAreaActiveStatusRow, error) {
	res, err := s.repo.ToggleAreaActiveStatus(id)
	if err != nil {
		return res, err
	}

	go func() {
		metadata, _ := json.Marshal(map[string]interface{}{
			"id":        res.ID,
			"is_active": res.IsActive.Bool,
		})

		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityAreas),
			Action:      string(pkg.LogTypeUpdate),
			Metadata:    json.RawMessage(metadata),
			EntityID:    id,
			PerformedBy: currentUserID,
		})
	}()

	return res, nil
}
