package controllers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"hubku/lapor_warga_be_v2/internal/modules/auth"
	"hubku/lapor_warga_be_v2/pkg"
)

type AuthController struct {
	authService auth.AuthService
	validator   *validator.Validate
}

func NewAuthController(s auth.AuthService, v *validator.Validate) *AuthController {
	return &AuthController{
		authService: s,
		validator:   v,
	}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	startTime := time.Now()

	var req auth.RegisterRequest
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

	userRow, err := c.authService.Register(req)
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
		"data": userRow,
		"meta": fiber.Map{
			"duration": time.Since(startTime).String(),
		},
	})
}

func (c *AuthController) Refresh(ctx *fiber.Ctx) error {
	startTime := time.Now()

	var req auth.RefreshRequest

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

	if req.RefreshToken == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "refresh_token is required",
				"meta": fiber.Map{
					"duration": time.Since(startTime).String(),
				},
			},
		)
	}

	resp, err := c.authService.RefreshToken(req)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
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
			"data": resp,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	startTime := time.Now()

	var req auth.LoginRequest

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

	resp, err := c.authService.Login(req)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
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
			"data": resp,
			"meta": fiber.Map{
				"duration": time.Since(startTime).String(),
			},
		},
	)
}
