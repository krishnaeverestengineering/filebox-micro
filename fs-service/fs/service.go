package fs

import (
	"context"

	"github.com/go-kit/kit/log"
)

type Service interface {
	CreateUser(ctx context.Context, userId string) (bool, error)
	CreateFolder(ctx context.Context, data UserFile) (interface{}, error)
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

func (s *FSService) CreateFolder(ctx context.Context, data UserFile) (interface{}, error) {
	// path := s.repo.GetFullPath(nil, data.ParentId, data.ParentId, data.UserID)
	// if path == "" {
	// 	return fmt.Errorf("path is empty", nil)
	// }
	// err := CreateFolder(data.UserID, data)
	// if err != nil {
	// 	return err
	// }
	files, e := s.repo.CreateFolder(ctx, data)
	if e != nil {
		return nil, e
	}
	return files, nil
}

func (s *FSService) ListDirectoryFiles(ctx context.Context, id string, userID string) ([]UserFile, error) {
	return s.repo.ListDirectoryFiles(ctx, id, userID)
}
