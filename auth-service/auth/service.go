package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
)

type Service interface {
	CreateUser(ctx context.Context, userId string, name string, email string) (bool, error)
	GetUser(ctx context.Context, userId string) (bool, User, error)
	GetSessionToken(userID string, secret []byte) (string, time.Time, error)
}

type LoginService struct {
	repo   Repository
	logger log.Logger
}

type Claims struct {
	UserName string
	jwt.StandardClaims
}

func NewService(repo Repository, logger log.Logger) Service {
	return &LoginService{
		repo:   repo,
		logger: logger,
	}
}

func (s LoginService) GetSessionToken(userID string, secret []byte) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	cliams := &Claims{
		UserName: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}

func (s LoginService) CreateUser(ctx context.Context, userId string, name string, email string) (bool, error) {
	logger := log.With(s.logger, "method", "CreateUser")
	user := User{
		UId:      userId,
		Name:     name,
		Root_dir: "test",
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return false, err
	}
	logger.Log("create user", userId)
	return true, nil
}

func (s LoginService) GetUser(ctx context.Context, userId string) (bool, User, error) {
	logger := log.With(s.logger, "method", "GetUser")
	user, err := s.repo.GetUser(ctx, userId)
	if err != nil {
		level.Error(logger).Log("err", err)
		return false, user, err
	}
	logger.Log("get user", userId)
	return true, user, nil
}
