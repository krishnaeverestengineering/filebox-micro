package fs

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

//NewHash creates unique id based on key provided
func NewHash(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

type Claims struct {
	UserName string
	jwt.StandardClaims
}

func DecodeCookie(r *http.Request) (string, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	tknStr := c.Value
	claims := &Claims{}
	_, e := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("my_sceret"), nil
	})
	if e != nil {
		return "", e
	}
	if claims.UserName == "null" {
		return "", fmt.Errorf("UserId null")
	}
	return claims.UserName, nil
}
