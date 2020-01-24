package main

import (
	"Filebox-Micro/fs-service/fs"
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
			"service: ", "filesystem",
			"time: ", log.DefaultTimestampUTC,
			"caller: ", log.DefaultCaller)
	}
	var srv fs.Service
	{
		config, err := fs.InitConfig()
		repo, err := fs.NewRepo(config.(fs.Config), nil)
		if err != nil {
			os.Exit(-1)
		}
		srv = fs.NewService(repo, logger)
	}

	endpoints := fs.MakeEndPoints(srv)

	errs := make(chan error)
	go func() {
		handler := fs.NewHTTPServer(ctx, endpoints)
		fmt.Println("Listening :8082")
		errs <- http.ListenAndServe("127.0.0.1:8082", handler)
	}()
	level.Error(logger).Log("exit", <-errs)
}
