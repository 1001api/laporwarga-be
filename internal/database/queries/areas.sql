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
    ST_SetSRID(
        ST_Multi(
            ST_CollectionExtract(
                ST_MakeValid(
                    ST_GeomFromGeoJSON(@boundary::text)
                ),
                3 -- Extract only polygon/multipolygon
            )
        ), 
        4326
    )
) RETURNING id;

-- name: CheckAreaExist :one
SELECT id FROM areas WHERE name = @name OR area_code = @area_code;

-- name: GetAreas :many
SELECT 
    id,
    name,
    description,
    area_type,
    area_code,
    CASE
        WHEN @simplify_tolerance < 0 THEN NULL 
        ELSE ST_AsGeoJSON(
            ST_Simplify(boundary, @simplify_tolerance::float)
        )::jsonb
    END AS boundary,
    ST_AsGeoJSON(center_point)::jsonb AS center_point,
    is_active,
    created_at
FROM
    areas
OFFSET @offset_count LIMIT @limit_count;

-- name: GetAreaBoundary :one
SELECT
    id,
    ST_AsGeoJSON(boundary)::jsonb AS boundary,
    ST_AsGeoJSON(center_point)::jsonb AS center_point
FROM areas
WHERE id = @id;

-- name: ToggleAreaActiveStatus :one
UPDATE
    areas
SET
    is_active = NOT is_active
WHERE id = @id RETURNING id, is_active;