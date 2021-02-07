package internal

import (
	"github.com/dgrijalva/jwt-go"
)

type (
	JWTClaims struct {
		Role string `json:"role"`
		jwt.StandardClaims
	}
)
