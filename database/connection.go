package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/serbanmarti/fiber_rest_api/database/model"
	"github.com/serbanmarti/fiber_rest_api/internal"
	"github.com/serbanmarti/fiber_rest_api/security"
)

// NewDBClientAndConnection create new MongoDB client and connection objects
func NewDBClientAndConnection(ctx context.Context, env *internal.Environ) (*mongo.Client, *mongo.Database) {
	dbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(env.DBUri))
	if err != nil {
		log.Fatalf("Failed to create DB client: %s", err)
	}

	// Check that the database connection has been established
	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Failed to establish a DB connection in a timely manner")
	}

	// Create the DB connection pipe
	dbConn := dbClient.Database(env.DBName)

	// Check for DB root account existence; create it if it does not exist
	rootExists, err := model.UserFindRoot(dbConn, env.DBRootUser)
	if err != nil {
		log.Fatalf("Failed to get root user: %s", err)
	}
	if !rootExists {
		// Generate a new random salt
		s, err := security.NewSalt()
		if err != nil {
			log.Fatalf("Failed to generate root user salt: %s", err)
		}

		// Hash and set the password field
		p := security.HashPassword(env.DBRootPass, s)

		// Create the new root user object
		u := &model.User{
			Email:    env.DBRootUser,
			Password: p,
			Salt:     s,
			Role:     "admin",
			Active:   true,
		}

		// Create the user in the DB
		err = model.UserCreateRoot(dbConn, u)
		if err != nil {
			log.Fatalf("Failed to create root user: %s", err)
		}
	}

	return dbClient, dbConn
}
