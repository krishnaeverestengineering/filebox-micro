package transport

import (
	"Filebox-Micro/authentication/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func DecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	userId := r.URL.Query().Get("userId")
	//var req model.CreateUserRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	return nil, err
	// }
	// return req, nil
	if userId == "" {
		return nil, fmt.Errorf("invalid userid")
	}
	return model.CreateUserRequest{UserID: userId, Name: "krishna", Email: "test"}, nil
}
