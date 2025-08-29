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