package server

import (
	"Filebox-Micro/authentication/endpoint"
	"Filebox-Micro/authentication/transport"
	"context"

	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/gin-gonic/gin"
)

//NewHTTPServer returns router configuration
func NewHTTPServer(ctx context.Context, endpoints endpoint.Endpoints) *gin.Engine {
	r := gin.Default()
	r.Use(commonMiddleware())
	r.GET("/ServiceLogin", func(c *gin.Context) {
		kithttp.NewServer(
			endpoints.CreateUser,
			transport.DecodeRequest,
			transport.EncodeResponse,
		).ServeHTTP(c.Writer, c.Request)
	})
	return r
}

func commonMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Context-Type", "application/json")
	}
}
