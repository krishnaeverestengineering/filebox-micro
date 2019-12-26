package fs

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type (
	CreateUserRequest struct {
		UserID string
	}
	CreateUserResponse struct {
		Ok bool `json:"ok"`
	}
	CreateFolderRequest struct {
		UserID   string `json:"userId"`
		FolderId string `json:"fid"`
		Name     string `json:"name"`
		ParentID string `json:"pid"`
	}
	CreateFolderResponse struct {
		Ok bool `json:"ok"`
	}
)

type Endpoints struct {
	CreateUser   endpoint.Endpoint
	CreateFolder endpoint.Endpoint
}

func MakeEndPoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:   makeCreateUserEndPoint(s),
		CreateFolder: makeCreateFolderEndPoint(s),
	}
}

func makeCreateUserEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		ok, err := s.CreateUser(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		return CreateUserResponse{Ok: ok}, nil
	}
}

func makeCreateFolderEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateFolderRequest)
		err := s.CreateFolder(ctx, UserFile{
			UserID:   req.UserID,
			FileID:   req.FolderId,
			FileName: req.Name,
			ParentId: req.ParentID,
		})
		if err != nil {
			return nil, err
		}
		return CreateFolderResponse{Ok: false}, nil
	}
}
