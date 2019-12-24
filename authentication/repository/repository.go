package repository

import (
	"Filebox-Micro/authentication/config"
	"Filebox-Micro/authentication/model"
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
)

type repo struct {
	db     *gorm.DB
	logger log.Logger
}

type Repository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id string) (model.User, error)
}

func New(config *config.Config, logger log.Logger) Repository {
	db, err := gorm.Open(config.Db.DatabaseUser, config.Db.DatabaseUri)
	if err != nil {
		return nil
	}
	return &repo{
		db:     db,
		logger: log.With(logger, "repository", "gormDB"),
	}
}

func (r *repo) CreateUser(ctx context.Context, user model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repo) GetUser(ctx context.Context, id string) (model.User, error) {
	var user model.User
	if err := r.db.Where(&model.User{UId: id}).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
