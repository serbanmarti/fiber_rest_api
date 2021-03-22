package server

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/serbanmarti/fiber_rest_api/database"
	"github.com/serbanmarti/fiber_rest_api/internal"
	"github.com/serbanmarti/fiber_rest_api/security"
	"github.com/serbanmarti/fiber_rest_api/server/handler"
)

func InitServer() (*fiber.App, int, *mongo.Client) {
	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		ErrorHandler:          internal.ErrorHandler,
		Prefork:               false,
		DisableStartupMessage: false,
	})

	// Get configuration from environment
	env, err := internal.GetEnv()
	if err != nil {
		logrus.Fatalf("Failed to get environment variables: %s", err)
	}

	// Configure the Fiber instance
	configureFiber(app)

	// Configure security
	if !env.SkipChecks {
		security.ConfigureSecurity(app, env)
	}

	// Create the database client & connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbClient, dbConn := database.NewDBClientAndConnection(ctx, env)

	// Configure the request validator
	v, err := internal.NewValidator()
	if err != nil {
		logrus.Fatalf("Failed to create request validator: %s", err)
	}

	// Initialize the route handler
	h := &handler.Handler{
		JwtSecret: env.JwtSecret,
		JwtExp:    env.JwtExp,
		DB:        dbConn,
		Validator: v,
	}

	// Assign the routes & handlers
	assignRoutesAndHandlers(app, h)

	return app, env.ServerPort, dbClient
}
