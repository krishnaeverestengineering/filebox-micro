package service

import (
	"Filebox-Micro/authentication/model"
	"Filebox-Micro/authentication/repository"
	"context"

	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
)

type Service interface {
	CreateUser(ctx context.Context, userId string, name string, email string) (bool, error)
	GetUser(ctx context.Context, userId string) (bool, error)
}

type LoginService struct {
	repo   repository.Repository
	logger log.Logger
}

func NewService(repo repository.Repository, logger log.Logger) Service {
	return &LoginService{
		repo:   repo,
		logger: logger,
	}
}

func (s LoginService) CreateUser(ctx context.Context, userId string, name string, email string) (bool, error) {
	logger := log.With(s.logger, "method", "CreateUser")
	user := model.User{
		UId:      userId,
		Name:     name,
		Root_dir: "test",
		//Email:  email,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return false, err
	}
	logger.Log("create user", userId)
	return true, nil
}

func (s LoginService) GetUser(ctx context.Context, userId string) (bool, error) {
	logger := log.With(s.logger, "method", "GetUser")
	if _, err := s.repo.GetUser(ctx, userId); err != nil {
		level.Error(logger).Log("err", err)
		return false, err
	}
	logger.Log("get user", userId)
	return true, nil
}
