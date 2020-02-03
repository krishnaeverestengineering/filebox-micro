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
		FileType   int    `json:"fileType"`
		FileFormat string `json:"fileFormat"`
	}
	CreateUserRequest struct {
		UserID string
	}
	CreateUserResponse struct {
		Ok bool `json:"ok"`
	}
	CreateFolderRequest struct {
		UserID     string   `json:"userId"`
		Name       string   `json:"name"`
		ParentID   string   `json:"pid"`
		FileType   FileType `json:"fileType"`
		FileFormat string   `json:"fileformat"`
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

	DeleteFileRequest struct {
		FileId string `json:"fid"`
		Name   string `json:"name"`
		UserID string
	}

	DeleteFileResponse struct {
		Ok bool `json:"ok"`
	}

	EditFileRequest struct {
		UserId  string `json:"userId"`
		FileId  string `json:"fid"`
		Content string `json:"content"`
	}

	EditFileResponse struct {
		Ok bool `json:"ok"`
	}

	GetFileContentRequest struct {
		UserID string
		FileID string
	}

	GetFileContentResponse struct {
		Content string   `json:"content"`
		File    UserFile `json:"file"`
		Ok      bool     `json:"ok"`
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
	CreateUser        endpoint.Endpoint
	CreateFolder      endpoint.Endpoint
	ListDirectory     endpoint.Endpoint
	DeleteFileOFolder endpoint.Endpoint
	EditFile          endpoint.Endpoint
	GetFileContent    endpoint.Endpoint
}

func MakeEndPoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:        makeCreateUserEndPoint(s),
		CreateFolder:      makeCreateFolderEndPoint(s),
		ListDirectory:     makeListDirectoryEndPoint(s),
		DeleteFileOFolder: makeDeleteFileOrFolderEndPoint(s),
		EditFile:          makeEditFileEndPoint(s),
		GetFileContent:    makeGetFileContentEndPoint(s),
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
			UserID:    req.UserID,
			FileName:  req.Name,
			ParentId:  req.ParentID,
			FileID:    utils.NewHMAC(),
			Type:      req.FileType,
			Extension: req.FileFormat,
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

func makeDeleteFileOrFolderEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteFileRequest)
		err := s.DeleteFileOrDirectory(ctx, req.FileId, req.Name, req.UserID)
		if err != nil {
			return DeleteFileResponse{Ok: false}, err
		}
		return DeleteFileResponse{Ok: true}, nil
	}
}

func makeEditFileEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(EditFileRequest)
		err := s.EditTextFileContent(ctx, req.FileId, req.Content, req.UserId)
		return EditFileResponse{Ok: err == nil}, err
	}
}

func makeGetFileContentEndPoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetFileContentRequest)
		file, c, err := s.GetFileContent(ctx, req.FileID, req.UserID)
		fmt.Println(c, err)
		return GetFileContentResponse{
			Content: c,
			File:    file.(UserFile),
		}, err
	}
}
