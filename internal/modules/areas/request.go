package areas

type CreateAreaRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	AreaType    string `json:"area_type" validate:"required"`
	AreaCode    string `json:"area_code" validate:"required"`
	GeoJSON     string `json:"geojson" validate:"required"`
}
