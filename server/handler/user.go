package handler

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/serbanmarti/fiber_rest_api/internal"
	"github.com/serbanmarti/fiber_rest_api/model"
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
		return internal.NewBackendError(internal.ErrBENotActive, nil, 1)
	}

	// Remove the password from memory (this is returned to the requester)
	u.Password = ""

	// Get input query parameters
	//qp := c.QueryParams()

	// Parse mobile request
	m := false
	//mRaw := qp.Get("mobile")
	//if mRaw != "" {
	//	m, err = strconv.ParseBool(mRaw)
	//	if err != nil {
	//		return internal.NewBackendError(internal.ErrBEQPInvalidMobile, nil, 1)
	//	}
	//}

	// JWT

	// If we have a mobile request, set a long expiration time
	var exp time.Duration
	if exp = h.JwtExp; m {
		exp, _ = time.ParseDuration("87600h")
	}

	// Set claims
	claims := &internal.JWTClaims{
		Role: u.Role,
		StandardClaims: jwt.StandardClaims{
			Id:        u.ID.Hex(),
			ExpiresAt: time.Now().Add(exp).Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Generate encoded token and send it as response.
	u.Token, err = token.SignedString([]byte(h.JwtSecret))
	if err != nil {
		return
	}

	return HTTPSuccess(c, u)
}
