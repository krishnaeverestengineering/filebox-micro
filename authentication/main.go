package main

import (
	"Filebox-Micro/authentication/config"
	"Filebox-Micro/authentication/endpoint"
	"Filebox-Micro/authentication/repository"
	"Filebox-Micro/authentication/server"
	"Filebox-Micro/authentication/service"
	"context"
	"filebox/login-service/login"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var db *login.Database
var jwtKey = []byte("dumy_key")

func main() {
	ctx := context.Background()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service: ", "authentication",
			"time: ", log.DefaultTimestampUTC,
			"caller: ", log.DefaultCaller)
	}
	var srv service.Service
	{
		config, _ := config.InitConfig()
		repo := repository.New(config, nil)
		srv = service.NewService(repo, logger)
	}

	endpoints := endpoint.MakeEndPoints(srv)

	errs := make(chan error)
	go func() {
		handler := server.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe("127.0.0.1:8080", handler)
	}()
	level.Error(logger).Log("exit", <-errs)
}
