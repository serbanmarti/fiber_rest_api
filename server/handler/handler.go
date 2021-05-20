package handler

import (
	"time"

	"github.com/serbanmarti/fiber_rest_api/internal"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Handler struct {
		JwtSecret string
		JwtExp    time.Duration
		DB        *mongo.Database
		Validator *internal.Validator
	}
)

// HTTPSuccess returns a formatted HTTP Success response
func HTTPSuccess(c *fiber.Ctx, d interface{}) error {
	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Action completed successfully",
		"data":    d,
	})
}
