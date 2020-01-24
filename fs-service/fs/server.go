package fs

import (
	"context"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

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
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
