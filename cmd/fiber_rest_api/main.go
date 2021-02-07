package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/serbanmarti/fiber_rest_api/internal"
	"github.com/serbanmarti/fiber_rest_api/model"
	"github.com/serbanmarti/fiber_rest_api/server"
	"github.com/serbanmarti/fiber_rest_api/server/handler"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		//DisableColors: true,
		FullTimestamp: true,
	})

	// Instantiate the Fiber REST server and database connection
	app, serverPort, db := server.InitServer()

	// Restricted Routes
	app.Get("/restricted", restricted)

	// Start the server
	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", serverPort)); err != nil {
			log.Fatalf("[FATAL] Failed to start the server: %s", err)
		}
	}()

	// Graceful shutdown of the server with a timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
	logrus.Infof("Gracefully shutting down the server")
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Disconnect(context.TODO()); err != nil {

		log.Fatalf("[FATAL] Failed to disconnect from the database: %s", err)
	}
	if err := app.Shutdown(); err != nil {
		log.Fatalf("[FATAL] Failed to gracefully shut down the server: %s", err)
	}
	log.Printf("[INFO] Server shut down complete")
}

func login(c *fiber.Ctx) (err error) {
	// Parse request data
	u := new(model.User)
	if err = c.BodyParser(u); err != nil {
		return
	}

	// Throws Unauthorized error
	if u.Email != "john@g.com" || u.Password != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Set claims
	claims := &internal.JWTClaims{
		Role: "admin",
		StandardClaims: jwt.StandardClaims{
			Id:        "12345",
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return handler.HTTPSuccess(c, fiber.Map{"token": t})
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(*internal.JWTClaims)
	return handler.HTTPSuccess(c, fiber.Map{"Msg": fmt.Sprintf("Welcome %s", claims.Id)})
}
