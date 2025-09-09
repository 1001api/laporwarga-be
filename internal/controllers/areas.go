package controllers

import (
	"hubku/lapor_warga_be_v2/internal/modules/areas"
	"hubku/lapor_warga_be_v2/pkg"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type AreasController struct {
	service   areas.AreaService
	validator *validator.Validate
}

func NewAreasController(s areas.AreaService, v *validator.Validate) *AreasController {
	return &AreasController{service: s, validator: v}
}

func (c *AreasController) CreateArea(ctx *fiber.Ctx) error {
	startTime := time.Now()

	currentUserID := ctx.Locals("user_id")
	currentUserUUID, err := uuid.Parse(cast.ToString(currentUserID))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "unauthenticated",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	var req areas.CreateAreaRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid json body",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	if err := pkg.ValidateInput(req, c.validator); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err,
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	createdID, err := c.service.CreateArea(currentUserUUID, req)
	if err != nil {
		if strings.Contains(err.Error(), "area already exist") {
			return ctx.Status(fiber.StatusConflict).JSON(
				fiber.Map{
					"error": err.Error(),
					"meta": fiber.Map{
						"duration": time.Since(startTime).String(),
					},
				},
			)
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		fiber.Map{
			"data": createdID.String(),
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *AreasController) GetAreas(ctx *fiber.Ctx) error {
	startTime := time.Now()

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 20)
	tolerance := ctx.Query("tolerance", "simple")

	areas, err := c.service.GetAreas(page, limit, pkg.AreaTolerance(tolerance))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		fiber.Map{
			"data": areas,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *AreasController) GetAreaBoundary(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id := ctx.Params("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid uuid",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	area, err := c.service.GetAreaBoundary(uid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		fiber.Map{
			"data": area,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *AreasController) ToggleAreaActiveStatus(ctx *fiber.Ctx) error {
	startTime := time.Now()

	currentUserID := ctx.Locals("user_id")
	currentUserUUID, err := uuid.Parse(cast.ToString(currentUserID))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error": "unauthenticated",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	id := ctx.Params("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "invalid uuid",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	res, err := c.service.ToggleAreaActiveStatus(currentUserUUID, uid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return ctx.JSON(
		fiber.Map{
			"data": res,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}
