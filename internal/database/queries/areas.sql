-- name: CreateArea :one
INSERT INTO areas (
    name,
    description,
    area_type,
    area_code,
    boundary
) VALUES (
    @name,
    @description,
    @area_type,
    @area_code,
    ST_SetSRID(ST_GeomFromGeoJSON(@boundary), 4326)
) RETURNING id;