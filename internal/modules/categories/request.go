package categories

type CreateCategoryRequest struct {
	Name      string `json:"name" form:"name" validate:"required,min=3,max=100"`
	Slug      string `json:"slug" form:"slug" validate:"required,min=3,max=100"`
	Icon      string `json:"icon" form:"icon"`
	Color     string `json:"color" form:"color"`
	IsActive  bool   `json:"is_active" form:"is_active"`
	SortOrder int    `json:"sort_order" form:"sort_order"`
}

type UpdateCategoryRequest struct {
	Name      string `json:"name" form:"name" validate:"required,min=3,max=100"`
	Slug      string `json:"slug" form:"slug" validate:"required,min=3,max=100"`
	Icon      string `json:"icon" form:"icon"`
	Color     string `json:"color" form:"color"`
	IsActive  bool   `json:"is_active" form:"is_active"`
	SortOrder int    `json:"sort_order" form:"sort_order"`
}

type SearchCategoryRequest struct {
	SearchTerm string `json:"search_term"`
	SortBy     string `json:"sort_by"`
	SortOrder  string `json:"sort_order"`
}
