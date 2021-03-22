package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/serbanmarti/fiber_rest_api/internal"
	"github.com/serbanmarti/fiber_rest_api/security"
)

type (
	User struct {
		ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		Email    string             `bson:"email" json:"email" validate:"required,email"`
		Password string             `bson:"password" json:"password" validate:"required,password"`
		Salt     []byte             `bson:"salt" json:"-"`
		Role     string             `bson:"role" json:"role,omitempty" validate:"omitempty,role"`
		Active   bool               `bson:"active" json:"-"`
	}
)

const (
	usersCollectionName = "users"
)

// Check if a user is found based on given email and password in the DB
func UserFind(m *mongo.Database, u *User) error {
	// Save the raw password for later use
	rawPassword := u.Password

	// Create a DB connection
	db := m.Collection(usersCollectionName)

	if err := db.FindOne(context.TODO(), bson.M{"email": u.Email}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return internal.NewError(internal.ErrDBNoData, err, 1)
		}

		return internal.NewError(internal.ErrDBQuery, err, 1)
	}

	// Create the hashed password for the current user
	hashedPassword := security.HashPassword(rawPassword, u.Salt)

	// Check if hashed passwords match
	if u.Password != hashedPassword {
		return internal.NewError(internal.ErrBEInvalidPassword, nil, 1)
	}

	return nil
}

// Check if the root user account is already available in the DB
func UserFindRoot(m *mongo.Database, email string) (bool, error) {
	// Create a DB connection
	db := m.Collection(usersCollectionName)

	count, err := db.CountDocuments(context.TODO(), bson.M{"email": email})
	if err != nil && err != mongo.ErrNoDocuments {
		return false, internal.NewError(internal.ErrDBQuery, err, 1)
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// Create a new root user in the DB
func UserCreateRoot(m *mongo.Database, u *User) error {
	// Create a DB connection
	db := m.Collection(usersCollectionName)

	// Add the user to the DB
	_, err := db.InsertOne(context.TODO(), u)
	if err != nil {
		return internal.NewError(internal.ErrDBInsert, err, 1)
	}

	return nil
}
