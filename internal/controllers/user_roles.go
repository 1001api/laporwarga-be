package controllers

import (
	db "hubku/lapor_warga_be_v2/internal/database/generated"
	userroles "hubku/lapor_warga_be_v2/internal/modules/user_roles"
	"hubku/lapor_warga_be_v2/pkg"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserRolesController struct {
	userRolesSvc userroles.UserRolesService
	validate     *validator.Validate
}

func NewUserRolesController(userRolesSvc userroles.UserRolesService, validate *validator.Validate) *UserRolesController {
	return &UserRolesController{
		userRolesSvc: userRolesSvc,
		validate:     validate,
	}
}

func (c *UserRolesController) ListAllRoles(ctx *fiber.Ctx) error {
	startTime := time.Now()

	result, err := c.userRolesSvc.ListAllRoles()
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

	return ctx.JSON(fiber.Map{
		"data": result,
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *UserRolesController) CreateRole(ctx *fiber.Ctx) error {
	startTime := time.Now()

	var req userroles.CreateRoleRequest

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

	if err := pkg.ValidateInput(req, c.validate); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err,
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	createdRole, err := c.userRolesSvc.CreateRole(db.CreateRoleParams{
		Name:        req.Name,
		Description: req.Description,
	})
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

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": createdRole,
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}
