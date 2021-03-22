package handler

import (
	"fmt"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/serbanmarti/fiber_rest_api/internal"
)

func (h *Handler) Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(*internal.JWTClaims)

	return HTTPSuccess(c, fiber.Map{"Msg": fmt.Sprintf("Welcome %s", claims.Id)})
}
