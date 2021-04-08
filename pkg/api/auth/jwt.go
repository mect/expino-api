package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	Name string `json:"name"`
	ID   uint   `json:"id"`
	// TODO: roles
	jwt.StandardClaims
}
