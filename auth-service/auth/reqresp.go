package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CreateUserSuccessRespose struct {
	Status          string `json:"status"`
	IsAuthenticated bool   `json:"is_authenticated"`
	Path            string `json:"path"`
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(CreateUserResponse)
	if res.Ok {
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(CreateUserSuccessRespose{
			Status:          "ok",
			IsAuthenticated: res.IsAuthenticated,
			Path:            "/",
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

func DecodeGetTokenRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var d TokenRequestBody
	er := json.NewDecoder(r.Body).Decode(&d)
	defer r.Body.Close()
	if er != nil {
		return nil, fmt.Errorf("invalid userid")
	}
	b, err := ioutil.ReadFile("./token.json")
	if err != nil {
		return nil, err
	}
	var t TokenInfo
	er = json.Unmarshal(b, &t)
	if er != nil {
		return nil, er
	}
	t.AccessToken.Sub = d.UserID
	t.RefreshToken.Sub = d.UserID
	return t, nil
}

func EncodeGetTokenResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	req := response.(TokenInfo)
	err := json.NewEncoder(w).Encode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	return err
}
