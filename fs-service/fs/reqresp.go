package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUserRespose struct {
	Status string
}

func DecodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		return nil, fmt.Errorf("invalid userid")
	}
	return CreateUserRequest{UserID: userId}, nil
}

func EncodeCreateUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(CreateUserResponse)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if res.Ok {
		err := encoder.Encode(CreateUserRespose{
			Status: "OK",
		})
		return err
	}
	err := encoder.Encode(CreateUserRespose{
		Status: "NOT OK",
	})
	return err
}

func DecodeCreateFolderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	c, err := DecodeCookie(r)
	fmt.Println(c, err)
	if err != nil {
		return nil, err
	}
	return CreateFolderRequest{
		UserID: c,
	}, nil
}

func EncodeCreateFolderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(CreateFolderResponse)
	if res.Ok {
		_, err := w.Write([]byte("hello"))
		return err
	}
	return nil
}

func DecodeListDirectoryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	c, err := DecodeCookie(r)
	fmt.Println(c, err)
	if err != nil {
		//return nil, err
	}
	fid := r.URL.Query().Get("folder")
	return ListDirectoryRequest{
		UserId:   c,
		FolderId: fid,
	}, nil
}

func EncodeListDirectoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(ListDirectoryResponse)
	fmt.Println(res.Files)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(res)
}
