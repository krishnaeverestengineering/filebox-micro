package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/jinzhu/gorm"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
}

type CreateUserSuccessRespose struct {
	Status          string `json:"status"`
	IsAuthenticated bool   `json:"is_authenticated"`
}

func MakeEndPoints(s Service) Endpoints {
	return Endpoints{
		CreateUser: MakeCreateUserEndPoint(s),
		GetUser:    MakeGetUserEndPoint(s),
	}
}

func MakeCreateUserEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		ok, _, err := s.GetUser(ctx, req.UserID)
		if gorm.IsRecordNotFoundError(err) {
			ok, err = s.CreateUser(ctx, req.UserID, req.Name, req.Email)
		}
		if err != nil {
			return nil, err
		}
		token, time, err := s.GetSessionToken(req.UserID, []byte("my_sceret"))
		if err != nil {
			return CreateUserResponse{Ok: false, IsAuthenticated: false}, err
		}
		return CreateUserResponse{Ok: ok, Token: token, ExprireAt: time, IsAuthenticated: true}, err
	}
}

func MakeGetUserEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		fmt.Println("world")
		ok, user, err := s.GetUser(ctx, req.Id)
		return GetUserResponse{Ok: ok, Data: user}, err
	}
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(CreateUserResponse)
	if res.Ok {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    res.Token,
			Expires:  res.ExprireAt,
			HttpOnly: true,
			Path:     "/",
		})
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
	userId := r.URL.Query().Get("userId")
	fmt.Println(r.Cookie("token"))
	if userId == "" {
		return nil, fmt.Errorf("invalid userid")
	}
	return CreateUserRequest{UserID: userId, Name: "krishna", Email: "test"}, nil
}
