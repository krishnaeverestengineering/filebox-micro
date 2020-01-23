package fs

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
)

type Service interface {
	CreateUser(ctx context.Context, userId string) (bool, error)
	CreateFolder(ctx context.Context, data UserFile) error
	ListDirectoryFiles(ctx context.Context, id string, userID string) ([]UserFile, error)
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

func (s *FSService) CreateFolder(ctx context.Context, data UserFile) error {
	path := s.repo.GetFullPath(nil, data.ParentId, data.RootId, data.UserID)
	if path == "" {
		return fmt.Errorf("path is empty", nil)
	}
	err := CreateFolder(path, data.FileName)
	if err != nil {
		return err
	}
	e := s.repo.CreateFolder(ctx, data)
	if e != nil {
		return e
	}

	fmt.Println(path)
	return nil
}

func (s *FSService) ListDirectoryFiles(ctx context.Context, id string, userID string) ([]UserFile, error) {
	return s.repo.ListDirectoryFiles(ctx, id, userID)
}
