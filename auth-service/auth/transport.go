package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/jinzhu/gorm"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
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
			return CreateUserResponse{Ok: false, IsAuthenticated: false}, err
		}
		return CreateUserResponse{Ok: ok, IsAuthenticated: true}, err
	}
}

func MakeGetUserEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		ok, user, err := s.GetUser(ctx, req.Id)
		return GetUserResponse{Ok: ok, Data: user}, err
	}
}
