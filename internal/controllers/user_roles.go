package controllers

import (
	db "hubku/lapor_warga_be_v2/internal/database/generated"
	userroles "hubku/lapor_warga_be_v2/internal/modules/user_roles"
	"hubku/lapor_warga_be_v2/pkg"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/cast"
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

	createdRole, err := c.userRolesSvc.CreateRole(currentUserUUID, db.CreateRoleParams{
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

func (c *UserRolesController) AssignRole(ctx *fiber.Ctx) error {
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
	targetUUID, err := uuid.Parse(id)
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

	var req userroles.AssignRoleRequest

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

	if err := c.userRolesSvc.AssignRoleToUser(currentUserUUID, db.AssignRoleToUserParams{
		UserID:   targetUUID,
		RoleName: req.RoleName,
	}); err != nil {
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
		"data": "success",
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *UserRolesController) GetRoleByID(ctx *fiber.Ctx) error {
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

	result, err := c.userRolesSvc.GetRoleByID(uid)
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

func (c *UserRolesController) GetRoleByName(ctx *fiber.Ctx) error {
	startTime := time.Now()

	name := ctx.Params("name")
	if name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "name is required",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	result, err := c.userRolesSvc.GetRoleByName(name)
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

func (c *UserRolesController) UpdateRole(ctx *fiber.Ctx) error {
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

	var req userroles.UpdateRoleRequest

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

	if err := c.userRolesSvc.UpdateRole(currentUserUUID, uid, req); err != nil {
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
		"data": "success",
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *UserRolesController) RemoveRole(ctx *fiber.Ctx) error {
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

	if err := c.userRolesSvc.RemoveRole(currentUserUUID, uid); err != nil {
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
		"data": "success",
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}
