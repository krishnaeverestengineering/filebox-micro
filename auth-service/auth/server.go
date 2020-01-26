package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

//NewHTTPServer returns router configuration
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.Methods("GET").Path("/ServiceLogin").Handler(kithttp.NewServer(
		endpoints.CreateUser,
		DecodeRequest,
		EncodeResponse,
	))
	r.Methods("GET").Path("/token").HandlerFunc(tokenHandler)
	r.Methods("GET").Path("/keys").HandlerFunc(keysHandler)
	return r
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
	fmt.Println(w.Header())
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

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		//w.Header().Add("Access-Control-Allow-Origin", "http://me.filebox.com:3000")
		next.ServeHTTP(w, r)
	})
}
