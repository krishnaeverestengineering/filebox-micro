package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserName string
	jwt.StandardClaims
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.Header().Add("Location", "http://localhost:3000/login")
				w.WriteHeader(302)
				return
			}
			w.Header().Add("location", "http://localhost:3000")
			w.WriteHeader(302)
			return
		}
		tknStr := c.Value
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("my_sceret"), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.Header().Add("location", "http://localhost:3000")
				w.WriteHeader(302)
				return
			}
			fmt.Println(err)
			w.Header().Add("location", "http://localhost:3000")
			w.WriteHeader(302)
			return
		}
		if !tkn.Valid {
			w.Header().Add("location", "http://localhost:3000")
			w.WriteHeader(302)
			return
		}
		//w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.UserName)))
		next.ServeHTTP(w, r)
	})
}
