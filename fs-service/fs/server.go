package fs

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func getCORS() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8081"}),
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
	//r.Use(middleware.AuthenticationMiddleware)
	r.Methods("GET").Path("/createUser").Handler(kithttp.NewServer(
		endpoints.CreateUser,
		DecodeCreateUserRequest,
		EncodeCreateUserResponse,
	))
	r.Methods("POST").Path("/createFolder").Handler(kithttp.NewServer(
		endpoints.CreateFolder,
		DecodeCreateFolderRequest,
		EncodeCreateFolderResponse,
	))

	r.Methods("GET").Path("/ls").Handler(kithttp.NewServer(
		endpoints.ListDirectory,
		DecodeListDirectoryRequest,
		EncodeListDirectoryResponse,
	))
	return getCORS()(r)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
