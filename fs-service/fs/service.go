package fs

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/go-kit/kit/log"
)

type Service interface {
	CreateUser(ctx context.Context, userId string) (bool, error)
	CreateFolder(ctx context.Context, data UserFile) (interface{}, error)
	ListDirectoryFiles(ctx context.Context, id string, userID string) ([]UserFile, error)
	DeleteFileOrDirectory(ctx context.Context, fid string, fname string, userID string) error
	//OpenFile(ctx context.Context, fid string, userID string) error
	GetFileContent(ctx context.Context, fid string, userID string) (interface{}, string, error)
	EditTextFileContent(ctx context.Context, fid string, content string, userID string) error
}

type FSService struct {
	repo   Repository
	logger log.Logger
}

func NewService(repo Repository, logger log.Logger) Service {
	return &FSService{
		repo:   repo,
		logger: logger,
	}
}

func (s *FSService) CreateUser(ctx context.Context, userId string) (bool, error) {
	return s.repo.CreateUser(ctx, userId)
}

func (s *FSService) CreateFolder(ctx context.Context, data UserFile) (interface{}, error) {
	path := s.repo.GetFullPath(nil, data.ParentId, "/", data.UserID)
	if path == "" {
		return nil, fmt.Errorf("path is empty", nil)
	}
	fmt.Println(path)
	err := CreateFolder(path, data)
	if err != nil {
		return nil, err
	}
	files, e := s.repo.CreateFolder(ctx, data)
	if e != nil {
		return nil, e
	}
	return files, nil
}

func (s *FSService) ListDirectoryFiles(ctx context.Context, id string, userID string) ([]UserFile, error) {
	return s.repo.ListDirectoryFiles(ctx, id, userID)
}

func (s *FSService) DeleteFileOrDirectory(ctx context.Context, fid string, fname string, userID string) error {
	path := s.repo.GetFullPath(nil, fid, "/", userID)
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	err := DeleteFolder(path, fname, userID)
	if err != nil {
		return err
	}
	s.repo.DeleteFileOrFolder(ctx, fid, userID)
	return nil
}

func (s *FSService) GetFileContent(ctx context.Context, fid string, userID string) (interface{}, string, error) {
	path := s.repo.GetFullPath(ctx, fid, "/", userID)
	if path == "" {
		return nil, "", fmt.Errorf("path is empty")
	}
	file, err := s.repo.GetDocument(ctx, fid, userID)
	if err != nil {
		return nil, "", err
	}
	uFile := file.(UserFile)
	content, err := ReadFile(filepath.Join(path+uFile.Extension), userID)
	return file, content.(string), err
}

func (s *FSService) EditTextFileContent(ctx context.Context, fid string, content string, userID string) error {
	path := s.repo.GetFullPath(ctx, fid, "/", userID)
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	file, err := s.repo.GetDocument(ctx, fid, userID)
	if err != nil {
		return err
	}
	uFile := file.(UserFile)
	return EditTextFile(filepath.Join(path+uFile.Extension), content, userID)
}
