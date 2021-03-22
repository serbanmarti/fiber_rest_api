package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	prefixed "github.com/t-tomalak/logrus-prefixed-formatter"

	"github.com/serbanmarti/fiber_rest_api/server"
)

func main() {
	// Configure base logging
	logrus.SetFormatter(&prefixed.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02T15:04:05.000000000Z07:00",
		FullTimestamp:   true,
		ForceFormatting: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	// Instantiate the Fiber REST API server and DB connection
	app, serverPort, db := server.InitServer()

	// Configure graceful shutdown of the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		<-quit
		logrus.Infof("Gracefully shutting down the server")

		if err := app.Shutdown(); err != nil {
			logrus.Fatalf("Failed to gracefully shut down the server: %s", err)
		}
		if err := db.Disconnect(context.TODO()); err != nil {
			logrus.Fatalf("Failed to gracefully disconnect from the DB: %s", err)
		}
	}()

	// Start the server
	if err := app.Listen(fmt.Sprintf(":%d", serverPort)); err != nil {
		logrus.Fatalf("Failed to start the server: %s", err)
	}

	logrus.Infof("Server shut down complete")
}
