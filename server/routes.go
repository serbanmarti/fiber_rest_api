package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/serbanmarti/fiber_rest_api/server/handler"
)

func assignRoutesAndHandlers(app *fiber.App, h *handler.Handler) {
	// Index
	app.Get("/", h.Index)

	// Users processes management
	app.Post("/login", h.Login)

	// Restricted
	app.Get("/restricted", h.Restricted)
}
