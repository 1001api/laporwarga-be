package categories

import (
	"fmt"
	db "hubku/lapor_warga_be_v2/internal/database/generated"
	"hubku/lapor_warga_be_v2/internal/modules/auditlogs"
	"hubku/lapor_warga_be_v2/pkg"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CategoriesService interface {
	CreateCategory(currentUserID uuid.UUID, req CreateCategoryRequest) (uuid.UUID, error)
	CheckCategoryExist(arg db.CheckCategoryExistParams) (bool, error)
	GetCategories() ([]db.GetCategoriesRow, error)
	GetCategoryById(id uuid.UUID) (db.GetCategoryByIdRow, error)
	GetCategoryBySlug(slug string) (db.GetCategoryBySlugRow, error)
	SearchCategories(req SearchCategoryRequest) ([]db.SearchCategoriesRow, error)
	ToggleCategoryActiveStatus(currentUserID uuid.UUID, id uuid.UUID) (db.ToggleCategoryActiveStatusRow, error)
	UpdateCategory(currentUserID uuid.UUID, id uuid.UUID, req UpdateCategoryRequest) (uuid.UUID, error)
	DeleteCategory(currentUserID uuid.UUID, id uuid.UUID) (uuid.UUID, error)
}

type service struct {
	repo       CategoriesRepository
	logService auditlogs.LogsService
}

func NewCategoriesService(repo CategoriesRepository, logService auditlogs.LogsService) CategoriesService {
	return &service{repo: repo, logService: logService}
}

func (s *service) CreateCategory(currentUserID uuid.UUID, req CreateCategoryRequest) (uuid.UUID, error) {
	// check category exist
	exist, err := s.CheckCategoryExist(db.CheckCategoryExistParams{
		Name: req.Name,
		Slug: req.Slug,
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	if exist {
		return uuid.UUID{}, fmt.Errorf(pkg.ErrExist)
	}

	result, err := s.repo.CreateCategory(db.CreateCategoryParams{
		Name: req.Name,
		Slug: req.Slug,
		Icon: pgtype.Text{
			String: req.Icon,
			Valid:  req.Icon != "",
		},
		Color: pgtype.Text{
			String: req.Color,
			Valid:  req.Color != "",
		},
		IsActive: pgtype.Bool{
			Valid: true,
			Bool:  req.IsActive,
		},
		SortOrder: pgtype.Int4{
			Valid: true,
			Int32: int32(req.SortOrder),
		},
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	// log create category
	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityCategories),
			Action:      string(pkg.LogTypeCreate),
			EntityID:    result,
			PerformedBy: currentUserID,
		})
	}()

	return result, nil
}

func (s *service) CheckCategoryExist(arg db.CheckCategoryExistParams) (bool, error) {
	return s.repo.CheckCategoryExist(arg)
}

func (s *service) GetCategories() ([]db.GetCategoriesRow, error) {
	return s.repo.GetCategories()
}

func (s *service) GetCategoryById(id uuid.UUID) (db.GetCategoryByIdRow, error) {
	return s.repo.GetCategoryById(id)
}

func (s *service) GetCategoryBySlug(slug string) (db.GetCategoryBySlugRow, error) {
	return s.repo.GetCategoryBySlug(slug)
}

func (s *service) SearchCategories(req SearchCategoryRequest) ([]db.SearchCategoriesRow, error) {
	req.SortBy = strings.ToLower(req.SortBy)
	req.SortOrder = strings.ToLower(req.SortOrder)
	return s.repo.SearchCategories(req.SearchTerm, req.SortBy, req.SortOrder)
}

func (s *service) ToggleCategoryActiveStatus(currentUserID, id uuid.UUID) (db.ToggleCategoryActiveStatusRow, error) {
	result, err := s.repo.ToggleCategoryActiveStatus(id)
	if err != nil {
		return db.ToggleCategoryActiveStatusRow{}, err
	}

	// log toggle category active status
	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityCategories),
			Action:      string(pkg.LogTypeUpdate),
			EntityID:    id,
			PerformedBy: currentUserID,
		})
	}()

	return result, nil
}

func (s *service) UpdateCategory(currentUserID, id uuid.UUID, req UpdateCategoryRequest) (uuid.UUID, error) {
	result, err := s.repo.UpdateCategory(db.UpdateCategoryParams{
		ID:        id,
		Name:      req.Name,
		Slug:      req.Slug,
		Icon:      req.Icon,
		Color:     req.Color,
		IsActive:  req.IsActive,
		SortOrder: int32(req.SortOrder),
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	// log update category
	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityCategories),
			Action:      string(pkg.LogTypeUpdate),
			EntityID:    id,
			PerformedBy: currentUserID,
		})
	}()

	return result, nil
}

func (s *service) DeleteCategory(currentUserID, id uuid.UUID) (uuid.UUID, error) {
	result, err := s.repo.DeleteCategory(id)
	if err != nil {
		return uuid.UUID{}, err
	}

	// log delete category
	go func() {
		s.logService.CreateLog(db.CreateAuditLogParams{
			EntityName:  string(pkg.LogEntityCategories),
			Action:      string(pkg.LogTypeDelete),
			EntityID:    id,
			PerformedBy: currentUserID,
		})
	}()

	return result, nil
}
