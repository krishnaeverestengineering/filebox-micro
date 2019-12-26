package fs

import (
	"context"
	"net/http"
)

func DecodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return CreateUserRequest{}, nil
}

func EncodeCreateUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(CreateUserResponse)
	if res.Ok {
		_, err := w.Write([]byte("hello"))
		return err
	}
	return nil
}

func DecodeCreateFolderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return CreateFolderRequest{}, nil
}

func EncodeCreateFolderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(CreateFolderResponse)
	if res.Ok {
		_, err := w.Write([]byte("hello"))
		return err
	}
	return nil
}
