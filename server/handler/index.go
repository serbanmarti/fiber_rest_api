package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Index(c *fiber.Ctx) error {
	return c.JSON("Service alive!")
}
