package endpoint

import (
	"Filebox-Micro/authentication/model"
	"Filebox-Micro/authentication/service"
	"context"

	"github.com/go-kit/kit/endpoint"
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
		ok, err := s.CreateUser(ctx, req.UserID, req.Name, req.Email)
		return model.CreateUserResponse{Ok: ok}, err
	}
}

func makeGetUserEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.GetUserRequest)
		ok, err := s.GetUser(ctx, req.Id)
		return model.GetUserResponse{Ok: ok}, err
	}
}
