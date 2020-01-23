package models

import "github.com/dgrijalva/jwt-go"

type AuthTokenClaims struct {
	UserId string
	jwt.StandardClaims
}
