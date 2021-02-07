package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/serbanmarti/fiber_rest_api/server/handler"
)

func assignRoutesAndHandlers(app *fiber.App, h *handler.Handler) {
	// Index
	app.Get("/", h.Index)

	// Tester
	//e.GET("/cache_test", h.CacheTest)

	// Metrics management
	//app.Get("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Users processes management
	app.Post("/login", h.Login)
	//e.POST("/invite", h.Invite)
	//e.POST("/validate_invite", h.ValidateInvite)
	//e.POST("/signup", h.SignUp)

	//// Users management
	//users := e.Group("/users")
	//
	//users.GET("", h.UserGetAll)
	//
	//users.PUT("/:userID", h.UserUpdate)
	//
	//users.DELETE("/:userID", h.UserDelete)
	//
	//// Stats
	//stats := e.Group("/stats")
	//
	//stats.GET("", h.StatsGetData)
}
