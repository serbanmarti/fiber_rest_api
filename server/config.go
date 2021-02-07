package server

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func configureFiber(app *fiber.App) {
	// Configure panic recovery
	app.Use(recover.New())

	// Configure logging
	//app.Use(logger.New(logger.Config{
	//	TimeFormat: "2006/01/02 15:04:05",
	//	Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
	//}))

	app.Use(func(c *fiber.Ctx) (err error) {
		// Handle request, store err for logging
		chainErr := c.Next()

		// Format error if exist
		formatErr := ""
		if chainErr != nil {
			formatErr = chainErr.Error()
		}

		logrus.Infof("%d - %s - %s", c.Response().StatusCode(), c.Path(), formatErr)
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
