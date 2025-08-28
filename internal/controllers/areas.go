package controllers

import (
	"hubku/lapor_warga_be_v2/internal/modules/areas"
	"hubku/lapor_warga_be_v2/pkg"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

	createdID, err := c.service.CreateArea(req)
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
			"data": createdID.String(),
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}
