DROP TRIGGER IF EXISTS set_areas_center ON areas;
DROP FUNCTION IF EXISTS calculate_center_point();
DROP INDEX IF EXISTS idx_areas_boundary;
DROP INDEX IF EXISTS idx_areas_area_code;
DROP TABLE IF EXISTS areas;