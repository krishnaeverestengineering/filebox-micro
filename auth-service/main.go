package main

import (
	ls "Filebox-Micro/auth-service/auth"
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
	var srv ls.Service
	{
		config, _ := ls.InitConfig()
		repo, err := ls.New(config, nil)
		if err != nil {
			fmt.Println("failed")
			os.Exit(-1)
		}
		srv = ls.NewService(repo, logger)
	}

	endpoints := ls.MakeEndPoints(srv)

	errs := make(chan error)
	go func() {
		handler := ls.NewHTTPServer(ctx, endpoints)
		fmt.Println("Listening :8081")
		errs <- http.ListenAndServe("127.0.0.1:8081", handler)
	}()
	level.Error(logger).Log("exit", <-errs)
}
