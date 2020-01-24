package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUserSuccessRespose struct {
	Status          string `json:"status"`
	IsAuthenticated bool   `json:"is_authenticated"`
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(CreateUserResponse)
	if res.Ok {
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(CreateUserSuccessRespose{
			Status:          "ok",
			IsAuthenticated: res.IsAuthenticated,
		})
		return err
	}
	return nil
}

func DecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	userId := r.Header.Get("UserID")
	if userId == "" {
		return nil, fmt.Errorf("invalid userid")
	}
	return CreateUserRequest{UserID: userId, Name: "krishna", Email: "test"}, nil
}
