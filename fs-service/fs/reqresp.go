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
	userId := r.Header.Get("UserID")
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
	userId := r.Header.Get("UserID")
	if userId == "" {
		return nil, fmt.Errorf("UserId not valid")
	}
	decoder := json.NewDecoder(r.Body)
	var data CreateFolderPostBody
	err := decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("data not valid")
	}
	return CreateFolderRequest{
		UserID:   userId,
		ParentID: data.ParenntId,
		Name:     data.FolderName,
	}, nil
}

func EncodeCreateFolderResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(CreateFolderResponse)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(res)
}

func DecodeListDirectoryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid := r.Header.Get("UserID")
	fid := r.URL.Query().Get("path")
	return ListDirectoryRequest{
		UserId:   uid,
		FolderId: fid,
	}, nil
}

func EncodeListDirectoryResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(ListDirectoryResponse)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(res)
}

func DecodeDeleteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	uid := r.Header.Get("UserID")
	fid := r.URL.Query().Get("path")
	return ListDirectoryRequest{
		UserId:   uid,
		FolderId: fid,
	}, nil
}

func EncodeDeleteResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(ListDirectoryResponse)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(res)
}
