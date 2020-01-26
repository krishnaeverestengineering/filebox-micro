package auth

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/jinzhu/gorm"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
	GetToken   endpoint.Endpoint
	GetKeys    endpoint.Endpoint
}

func MakeEndPoints(s Service) Endpoints {
	return Endpoints{
		CreateUser: makeCreateUserEndPoint(s),
		GetUser:    makeGetUserEndPoint(s),
		GetToken:   makeGetTokenEndPoint(s),
		GetKeys:    makeGetKeysEndPoint(s),
	}
}

func makeGetKeysEndPoint(s Service) endpoint.Endpoint {
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

func makeGetTokenEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if request == nil {
			return nil, fmt.Errorf("Invalid request")
		}
		req := request.(TokenInfo)
		return req, nil
	}
}

func makeCreateUserEndPoint(s Service) endpoint.Endpoint {
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

func makeGetUserEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		ok, user, err := s.GetUser(ctx, req.Id)
		return GetUserResponse{Ok: ok, Data: user}, err
	}
}
