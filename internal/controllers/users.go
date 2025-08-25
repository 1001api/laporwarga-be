package controllers

import (
	db "hubku/lapor_warga_be_v2/internal/database/generated"
	"hubku/lapor_warga_be_v2/internal/modules/users"
	"hubku/lapor_warga_be_v2/pkg"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type UserController struct {
	userSvc  users.UserService
	validate *validator.Validate
}

func NewUserController(s users.UserService, v *validator.Validate) *UserController {
	return &UserController{
		userSvc:  s,
		validate: v,
	}
}

func (c *UserController) GetMasterUser(ctx *fiber.Ctx) error {
	startTime := time.Now()

	// Parse pagination query parameters
	page := cast.ToInt(ctx.Query("page"))
	limit := cast.ToInt(ctx.Query("limit"))

	result, err := c.userSvc.GetUsers(db.GetUsersParams{
		OffsetCount: cast.ToInt32((page - 1) * limit),
		LimitCount:  cast.ToInt32(limit),
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

	return ctx.JSON(fiber.Map{
		"data": result,
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *UserController) SearchUser(ctx *fiber.Ctx) error {
	startTime := time.Now()

	query := ctx.Query("query")
	page := cast.ToInt32(ctx.Query("page"))
	limit := cast.ToInt32(ctx.Query("limit"))

	results, err := c.userSvc.SearchUser(query, page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}

	endTime := time.Now()

	return ctx.JSON(fiber.Map{
		"data": results,
		"meta": fiber.Map{
			"page":     page,
			"limit":    limit,
			"duration": endTime.Sub(startTime).String(),
		},
	})
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	startTime := time.Now()

	id, err := uuid.Parse(ctx.Params("id"))
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

	result, err := c.userSvc.GetUserByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"error": "user not found",
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

func (uc *UserController) GetUserByUsername(c *fiber.Ctx) error {
	startTime := time.Now()

	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "username is required",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	result, err := uc.userSvc.GetUserByIdentifier(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"error": "user not found",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	return c.JSON(fiber.Map{
		"data": result,
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *UserController) GetCurrentUser(ctx *fiber.Ctx) error {
	startTime := time.Now()

	userID := ctx.Locals("user_id")

	uid, err := uuid.Parse(cast.ToString(userID))
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

	result, err := c.userSvc.GetUserByID(uid)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(
			fiber.Map{
				"error": "user not found",
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

func (c *UserController) UpdateCurrentUser(ctx *fiber.Ctx) error {
	startTime := time.Now()

	uid, err := uuid.Parse(cast.ToString(ctx.Locals("user_id")))
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

	var req users.UpdateUserRequest

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

	if err := c.userSvc.UpdateUser(uid, uid, req); err != nil {
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

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
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

	uid, err := uuid.Parse(ctx.Params("id"))
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

	var req users.UpdateUserRequest

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

	if err := c.userSvc.UpdateUser(uid, currentUserUUID, req); err != nil {
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

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	startTime := time.Now()

	currentUser := ctx.Locals("user_id")
	currentUserUUID, err := uuid.Parse(cast.ToString(currentUser))
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

	uid, err := uuid.Parse(ctx.Params("id"))
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

	if err := c.userSvc.DeleteUser(uid, currentUserUUID); err != nil {
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

func (c *UserController) RestoreUser(ctx *fiber.Ctx) error {
	startTime := time.Now()

	currentUser := ctx.Locals("user_id")
	currentUserUUID, err := uuid.Parse(cast.ToString(currentUser))
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

	uid, err := uuid.Parse(ctx.Params("id"))
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

	if err := c.userSvc.RestoreUser(uid, currentUserUUID); err != nil {
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
