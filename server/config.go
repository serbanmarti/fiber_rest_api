package server

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/serbanmarti/fiber_rest_api/internal"
)

func configureFiber(app *fiber.App) {
	// Configure panic recovery
	app.Use(recover.New())

	// Configure logging
	app.Use(func(c *fiber.Ctx) (err error) {
		// Handle request, store err for logging
		chainErr := c.Next()

		// Format error if one exists
		formatErr := ""
		if chainErr != nil {
			return chainErr
		}

		internal.LogRequestResponse(c.Response().StatusCode(), c.IP(), c.Method(), c.Path(), formatErr)
		return nil
	})

	// Configure CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodHead,
			fiber.MethodOptions,
		}, ","),
		AllowCredentials: true,
	}))
}
