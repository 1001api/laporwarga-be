package categories

import (
	"context"
	db "hubku/lapor_warga_be_v2/internal/database/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoriesRepository interface {
	CreateCategory(arg db.CreateCategoryParams) (uuid.UUID, error)
	CheckCategoryExist(arg db.CheckCategoryExistParams) (bool, error)
	GetCategories() ([]db.GetCategoriesRow, error)
	GetCategoryById(id uuid.UUID) (db.GetCategoryByIdRow, error)
	GetCategoryBySlug(slug string) (db.GetCategoryBySlugRow, error)
	SearchCategories(searchTerm, orderBy, sortOrder string) ([]db.SearchCategoriesRow, error)
	ToggleCategoryActiveStatus(id uuid.UUID) (db.ToggleCategoryActiveStatusRow, error)
	UpdateCategory(arg db.UpdateCategoryParams) (uuid.UUID, error)
	DeleteCategory(id uuid.UUID) (uuid.UUID, error)
}

type repository struct {
	db *db.Queries
}

func NewCategoriesRepository(pool *pgxpool.Pool) CategoriesRepository {
	return &repository{db: db.New(pool)}
}

func (r *repository) CreateCategory(arg db.CreateCategoryParams) (uuid.UUID, error) {
	return r.db.CreateCategory(context.Background(), arg)
}

func (r *repository) CheckCategoryExist(arg db.CheckCategoryExistParams) (bool, error) {
	return r.db.CheckCategoryExist(context.Background(), arg)
}

func (r *repository) GetCategories() ([]db.GetCategoriesRow, error) {
	return r.db.GetCategories(context.Background())
}

func (r *repository) GetCategoryById(id uuid.UUID) (db.GetCategoryByIdRow, error) {
	return r.db.GetCategoryById(context.Background(), id)
}

func (r *repository) GetCategoryBySlug(slug string) (db.GetCategoryBySlugRow, error) {
	return r.db.GetCategoryBySlug(context.Background(), slug)
}

func (r *repository) SearchCategories(searchTerm, orderBy, sortOrder string) ([]db.SearchCategoriesRow, error) {
	// whitelist sorted by column
	col, ok := map[string]string{
		"name":       "name",
		"slug":       "slug",
		"icon":       "icon",
		"color":      "color",
		"is_active":  "is_active",
		"sort_order": "sort_order",
	}[orderBy]
	if !ok {
		orderBy = "name"
	}

	// whitelist sorted order
	order, ok := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}[sortOrder]
	if !ok {
		sortOrder = "ASC"
	}

	return r.db.SearchCategories(context.Background(), db.SearchCategoriesParams{
		SearchTerm: searchTerm,
		SortBy:     col,
		SortOrder:  order,
	})
}

func (r *repository) ToggleCategoryActiveStatus(id uuid.UUID) (db.ToggleCategoryActiveStatusRow, error) {
	return r.db.ToggleCategoryActiveStatus(context.Background(), id)
}

func (r *repository) UpdateCategory(arg db.UpdateCategoryParams) (uuid.UUID, error) {
	return r.db.UpdateCategory(context.Background(), arg)
}

func (r *repository) DeleteCategory(id uuid.UUID) (uuid.UUID, error) {
	return r.db.DeleteCategory(context.Background(), id)
}
