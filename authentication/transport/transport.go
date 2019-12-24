package transport

import (
	"Filebox-Micro/authentication/model"
	"Filebox-Micro/authentication/service"
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/jinzhu/gorm"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
}

func MakeEndPoints(s service.Service) Endpoints {
	return Endpoints{
		CreateUser: makeCreateUserEndPoint(s),
		GetUser:    makeGetUserEndPoint(s),
	}
}

func makeCreateUserEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.CreateUserRequest)
		var err error
		var ok bool
		ok, user, err := s.GetUser(ctx, req.UserID)
		if gorm.IsRecordNotFoundError(err) {
			ok, err = s.CreateUser(ctx, req.UserID, req.Name, req.Email)
		}
		token, time, err := s.GetSessionToken(req.UserID, []byte("my_sceret"))
		if err != nil {
			return model.CreateUserResponse{Ok: false}, err
		}
		fmt.Println(token, user)
		return model.CreateUserResponse{Ok: ok, Token: token, ExprireAt: time}, err
	}
}

func makeGetUserEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.GetUserRequest)
		ok, user, err := s.GetUser(ctx, req.Id)
		return model.GetUserResponse{Ok: ok, Data: user}, err
	}
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(model.CreateUserResponse)
	if res.Ok {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    res.Token,
			Expires:  res.ExprireAt,
			HttpOnly: true,
		})
		_, err := w.Write([]byte("hello"))
		return err
	}
	return nil
}

func DecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		return nil, fmt.Errorf("invalid userid")
	}
	return model.CreateUserRequest{UserID: userId, Name: "krishna", Email: "test"}, nil
}
