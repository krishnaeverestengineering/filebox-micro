package fs

import (
	"Filebox-Micro/fs-service/utils"
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type (
	CreateFolderPostBody struct {
		FolderName string `json:"name"`
		ParenntId  string `json:"pid"`
	}
	CreateUserRequest struct {
		UserID string
	}
	CreateUserResponse struct {
		Ok bool `json:"ok"`
	}
	CreateFolderRequest struct {
		UserID   string `json:"userId"`
		Name     string `json:"name"`
		ParentID string `json:"pid"`
	}
	CreateFolderResponse struct {
		Ok    bool       `json:"ok"`
		Files []UserFile `json:"files"`
	}

	ListDirectoryRequest struct {
		UserId   string `json:"userId"`
		FolderId string `json:"pid"`
	}
	ListDirectoryResponse struct {
		Ok    bool       `json:"ok"`
		Files []UserFile `json:"files"`
	}
)

type File struct {
	FName string `json:"name"`
	FType string `json:"type"`
	FTime int64  `json:"time"`
	FSize int64  `json:"size"`
	FPath string `json:"path,omitempty"`
}

type Endpoints struct {
	CreateUser    endpoint.Endpoint
	CreateFolder  endpoint.Endpoint
	ListDirectory endpoint.Endpoint
}

func MakeEndPoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:    makeCreateUserEndPoint(s),
		CreateFolder:  makeCreateFolderEndPoint(s),
		ListDirectory: makeListDirectoryEndPoint(s),
	}
}

func makeCreateUserEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("Createuser")
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
		files, err := s.CreateFolder(ctx, UserFile{
			UserID:   req.UserID,
			FileName: req.Name,
			ParentId: req.ParentID,
			FileID:   utils.NewHMAC(),
		})
		if err != nil {
			return CreateFolderResponse{Ok: false, Files: []UserFile{}}, err
		}
		return CreateFolderResponse{Ok: true, Files: files.([]UserFile)}, nil
	}
}

func makeListDirectoryEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListDirectoryRequest)
		files, err := s.ListDirectoryFiles(ctx, req.FolderId, req.UserId)
		if err != nil {
			return nil, err
		}
		return ListDirectoryResponse{Ok: true, Files: files}, nil
	}
}
