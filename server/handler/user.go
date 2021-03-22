package handler

import (
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/serbanmarti/fiber_rest_api/database/model"
	"github.com/serbanmarti/fiber_rest_api/internal"
)

// Login a user into the system and return an authorization JWT
func (h *Handler) Login(c *fiber.Ctx) (err error) {
	// Bind request data
	u := new(model.User)
	if err = c.BodyParser(&u); err != nil {
		return
	}

	// Validate request data
	if err = h.Validator.Validate(u); err != nil {
		return
	}

	// Find the user in the DB
	if err = model.UserFind(h.DB, u); err != nil {
		return
	}

	// Check if the account is active
	if !u.Active {
		return internal.NewError(internal.ErrBENotActive, nil, 1)
	}

	// JWT

	// Set claims
	claims := &internal.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        u.ID.Hex(),
			ExpiresAt: time.Now().Add(h.JwtExp).Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Generate encoded token and send it as response
	Token, err := token.SignedString([]byte(h.JwtSecret))
	if err != nil {
		return
	}

	return HTTPSuccess(c, fiber.Map{
		"Token": Token,
	})
}
