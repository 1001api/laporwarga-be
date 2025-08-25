package controllers

import (
	"hubku/lapor_warga_be_v2/internal/modules/auditlogs"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LogsController struct {
	logsService auditlogs.LogsService
}

func NewLogsController(s auditlogs.LogsService) *LogsController {
	return &LogsController{
		logsService: s,
	}
}

func (c *LogsController) ListLogs(ctx *fiber.Ctx) error {
	startTime := time.Now()

	result, err := c.logsService.GetLogs()
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
