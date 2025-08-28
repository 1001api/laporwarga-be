CREATE TABLE IF NOT EXISTS areas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    name TEXT NOT NULL,
    description TEXT,
    area_type VARCHAR(30) NOT NULL DEFAULT 'kabupaten',
    area_code VARCHAR(20) NOT NULL,
    is_active BOOL DEFAULT false,

    parent_id UUID REFERENCES areas(id) ON DELETE SET NULL, -- future hierarchy (provinsi or kecamatan)

    boundary GEOMETRY(POLYGON, 4326) NOT NULL
        CHECK (ST_IsValid(boundary) AND GeometryType(boundary) = 'POLYGON' AND ST_SRID(boundary) = 4326),
    
    center_point GEOMETRY(POINT, 4326)
        CHECK (center_point IS NULL OR (ST_IsValid(center_point) AND ST_SRID(center_point) = 4326)),
   
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- GIST Index
CREATE INDEX idx_areas_boundary ON areas USING GIST(boundary);

-- Indexes for faster lookups
CREATE INDEX idx_areas_area_code ON areas(area_code);

-- function to calculate center point
CREATE OR REPLACE FUNCTION calculate_center_point()
RETURNS TRIGGER AS $$
BEGIN
    NEW.center_point = ST_Centroid(NEW.boundary);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- attach trigger
CREATE TRIGGER set_areas_center
BEFORE INSERT OR UPDATE ON areas
FOR EACH ROW EXECUTE FUNCTION calculate_center_point();