package controllers

import (
	"hubku/lapor_warga_be_v2/internal/modules/categories"
	"hubku/lapor_warga_be_v2/pkg"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type CategoriesController struct {
	service   categories.CategoriesService
	validator *validator.Validate
}

func NewCategoriesController(s categories.CategoriesService, v *validator.Validate) *CategoriesController {
	return &CategoriesController{service: s, validator: v}
}

func (c *CategoriesController) CreateCategory(ctx *fiber.Ctx) error {
	startTime := time.Now()

	currentUserID := ctx.Locals("user_id")
	currentUserUUID, err := uuid.Parse(cast.ToString(currentUserID))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			pkg.ErrorResponse{
				Error: "unauthenticated",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	var req categories.CreateCategoryRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			pkg.ErrorResponse{
				Error: "invalid json body",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	if err := pkg.ValidateInput(req, c.validator); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			pkg.ErrorResponse{
				Error: err,
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	createdID, err := c.service.CreateCategory(currentUserUUID, req)
	if err != nil {
		if err.Error() == pkg.ErrExist {
			return ctx.Status(fiber.StatusConflict).JSON(
				pkg.ErrorResponse{
					Error: "category already exist",
					Meta: pkg.Meta{
						Duration: time.Since(startTime).String(),
					},
				},
			)
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(
			pkg.ErrorResponse{
				Error: "internal server error",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		pkg.SuccessResponse{
			Data: createdID.String(),
			Meta: pkg.Meta{
				Duration: time.Since(startTime).String(),
			},
		},
	)
}

func (c *CategoriesController) GetCategories(ctx *fiber.Ctx) error {
	startTime := time.Now()

	categories, err := c.service.GetCategories()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			pkg.ErrorResponse{
				Error: "internal server error",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		pkg.SuccessResponse{
			Data: categories,
			Meta: pkg.Meta{
				Duration: time.Since(startTime).String(),
			},
		},
	)
}

func (c *CategoriesController) GetCategoryById(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			pkg.ErrorResponse{
				Error: "category id is required",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	category, err := c.service.GetCategoryById(uuid.Must(uuid.Parse(id)))
	if err != nil {
		if err.Error() == pkg.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(
				pkg.ErrorResponse{
					Error: "category not found",
					Meta: pkg.Meta{
						Duration: time.Since(startTime).String(),
					},
				},
			)
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(
			pkg.ErrorResponse{
				Error: "internal server error",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		pkg.SuccessResponse{
			Data: category,
			Meta: pkg.Meta{
				Duration: time.Since(startTime).String(),
			},
		},
	)
}

func (c *CategoriesController) GetCategoryBySlug(ctx *fiber.Ctx) error {
	startTime := time.Now()

	slug := ctx.Params("slug")
	if slug == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			pkg.ErrorResponse{
				Error: "slug is required",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	category, err := c.service.GetCategoryBySlug(slug)
	if err != nil {
		if err.Error() == pkg.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(
				pkg.ErrorResponse{
					Error: "category not found",
					Meta: pkg.Meta{
						Duration: time.Since(startTime).String(),
					},
				},
			)
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(
			pkg.ErrorResponse{
				Error: "internal server error",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		pkg.SuccessResponse{
			Data: category,
			Meta: pkg.Meta{
				Duration: time.Since(startTime).String(),
			},
		},
	)
}

func (c *CategoriesController) SearchCategories(ctx *fiber.Ctx) error {
	startTime := time.Now()

	sortBy := ctx.Query("sort_by")
	sortOrder := ctx.Query("sort_order")
	searchTerm := ctx.Query("search")

	categories, err := c.service.SearchCategories(categories.SearchCategoryRequest{
		SearchTerm: searchTerm,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			pkg.ErrorResponse{
				Error: "internal server error",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		pkg.SuccessResponse{
			Data: categories,
			Meta: pkg.Meta{
				Duration: time.Since(startTime).String(),
			},
		},
	)
}

func (c *CategoriesController) ToggleCategoryActiveStatus(ctx *fiber.Ctx) error {
	startTime := time.Now()

	currentUserID := ctx.Locals("user_id")
	currentUserUUID, err := uuid.Parse(cast.ToString(currentUserID))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			pkg.ErrorResponse{
				Error: "unauthenticated",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	id := ctx.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			pkg.ErrorResponse{
				Error: "invalid category id",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	category, err := c.service.ToggleCategoryActiveStatus(currentUserUUID, uuid)
	if err != nil {
		if err.Error() == pkg.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(
				pkg.ErrorResponse{
					Error: "category not found",
					Meta: pkg.Meta{
						Duration: time.Since(startTime).String(),
					},
				},
			)
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(
			pkg.ErrorResponse{
				Error: "internal server error",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		pkg.SuccessResponse{
			Data: category,
			Meta: pkg.Meta{
				Duration: time.Since(startTime).String(),
			},
		},
	)
}

func (c *CategoriesController) UpdateCategory(ctx *fiber.Ctx) error {
	startTime := time.Now()

	currentUserID := ctx.Locals("user_id")
	currentUserUUID, err := uuid.Parse(cast.ToString(currentUserID))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			pkg.ErrorResponse{
				Error: "unauthenticated",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	id := ctx.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			pkg.ErrorResponse{
				Error: "invalid category id",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	var req categories.UpdateCategoryRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			pkg.ErrorResponse{
				Error: "invalid json body",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	if err := pkg.ValidateInput(req, c.validator); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			pkg.ErrorResponse{
				Error: err,
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	category, err := c.service.UpdateCategory(currentUserUUID, uuid, req)
	if err != nil {
		if err.Error() == pkg.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(
				pkg.ErrorResponse{
					Error: "category not found",
					Meta: pkg.Meta{
						Duration: time.Since(startTime).String(),
					},
				},
			)
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(
			pkg.ErrorResponse{
				Error: "internal server error",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		pkg.SuccessResponse{
			Data: category,
			Meta: pkg.Meta{
				Duration: time.Since(startTime).String(),
			},
		},
	)
}

func (c *CategoriesController) DeleteCategory(ctx *fiber.Ctx) error {
	startTime := time.Now()

	currentUserID := cast.ToString(ctx.Locals("user_id"))
	currentUserUUID, err := uuid.Parse(currentUserID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			pkg.ErrorResponse{
				Error: "unauthorized",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	id := ctx.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			pkg.ErrorResponse{
				Error: "invalid category id",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	category, err := c.service.DeleteCategory(currentUserUUID, uuid)
	if err != nil {
		if err.Error() == pkg.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(
				pkg.ErrorResponse{
					Error: "category not found",
					Meta: pkg.Meta{
						Duration: time.Since(startTime).String(),
					},
				},
			)
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(
			pkg.ErrorResponse{
				Error: "internal server error",
				Meta: pkg.Meta{
					Duration: time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		pkg.SuccessResponse{
			Data: category,
			Meta: pkg.Meta{
				Duration: time.Since(startTime).String(),
			},
		},
	)
}
