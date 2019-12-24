package main

import (
	"Filebox-Micro/authentication/config"
	"Filebox-Micro/authentication/repository"
	"Filebox-Micro/authentication/server"
	"Filebox-Micro/authentication/service"
	"Filebox-Micro/authentication/transport"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

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
		repo, err := repository.New(config, nil)
		if err != nil {
			os.Exit(-1)
		}
		srv = service.NewService(repo, logger)
	}

	endpoints := transport.MakeEndPoints(srv)

	errs := make(chan error)
	go func() {
		handler := server.NewHTTPServer(ctx, endpoints)
		fmt.Println("Listening :8080")
		errs <- http.ListenAndServe("127.0.0.1:8080", handler)
	}()
	level.Error(logger).Log("exit", <-errs)
}
