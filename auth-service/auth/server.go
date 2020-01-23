package auth

import (
	"Filebox-Micro/api-gateway/models"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func SetRoutes(r *mux.Router, endpoints Endpoints) {
	r.Use(commonMiddleware)
	r.Methods("GET").Path("/ServiceLogin").Handler(kithttp.NewServer(
		endpoints.CreateUser,
		DecodeRequest,
		EncodeResponse,
	))

	r.Methods("GET").Path("/login").HandlerFunc(ArticlesCategoryHandler)
}

func getCORS() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000",
			"http://127.0.0.1:8082"}),
		handlers.AllowedMethods([]string{"PUT", "PATCH", "GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Origin",
			"Content-Type",
			"X-Requested-With",
			"Authorization",
		}),
		handlers.ExposedHeaders([]string{"Content-Length", "Set-Cookie", "Cookie"}),
		handlers.AllowCredentials(),
		handlers.AllowedOriginValidator(func(origin string) bool {
			return origin == "http://localhost:3000"
		}),
	)
}

//NewHTTPServer returns router configuration
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.Methods("GET").Path("/ServiceLogin").Handler(kithttp.NewServer(
		endpoints.CreateUser,
		DecodeRequest,
		EncodeResponse,
	))

	r.Methods("GET").Path("/token-issuer").HandlerFunc(ArticlesCategoryHandler)
	r.Methods("GET").Path("/token").HandlerFunc(tokenHandler)
	r.Methods("GET").Path("/keys").HandlerFunc(keysHandler)
	// var http_dir *string
	// http_dir = flag.String("d", ".", "stats directory to serve from")

	// r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(*http_dir))))
	r.Methods("GET").Path("/private").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&CreateUserResponse{
			Ok: true,
		})
	})
	//r.Methods("GET").Path("/token-issuer").HandlerFunc(ArticlesCategoryHandler)
	return r
}

type TokeRes struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expires      time.Time `json:"exp"`
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadFile("./token.json")
	rawIn := json.RawMessage(string(b))
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(rawIn, &objmap)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(objmap)
}

func keysHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadFile("./keys.json")
	rawIn := json.RawMessage(string(b))
	var objmap map[string]*json.RawMessage
	err := json.Unmarshal(rawIn, &objmap)
	if err != nil {
		fmt.Println(err)
	}
	json.NewEncoder(w).Encode(objmap)
}

func GetSessionToken(uid string, secret []byte) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	cliams := &models.AuthTokenClaims{}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	_, time, _ := GetSessionToken("123", []byte("sceret"))
	json.NewEncoder(w).Encode(&TokeRes{
		AccessToken:  "hell",
		RefreshToken: "hell",
		Expires:      time,
	})
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
