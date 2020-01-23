package middlewares

import (
	. "Filebox-Micro/api-gateway/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func TokenValidater(r *http.Request) bool {
	c, err := r.Cookie("auth")
	if err != nil {
		return false
	}
	tknStr := c.Value
	claims := &AuthTokenClaims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("my_sceret"), nil
	})
	if err != nil || !tkn.Valid {
		return false
	}
	return true
}
