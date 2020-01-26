package main

import (
	"Filebox-Micro/api-gateway/jujumux"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	cel "github.com/devopsfaith/krakend-cel"
	muxcors "github.com/devopsfaith/krakend-cors/mux"
	jose "github.com/devopsfaith/krakend-jose"
	muxjose "github.com/devopsfaith/krakend-jose/mux"
	_ "github.com/devopsfaith/krakend-martian"
	martian "github.com/devopsfaith/krakend-martian"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/router/gorilla"
	"github.com/devopsfaith/krakend/router/mux"
	"github.com/devopsfaith/krakend/transport/http/client"
	_ "github.com/dgrijalva/jwt-go"
	g "github.com/gorilla/mux"
	_ "github.com/pismo/martians"
	_ "gopkg.in/square/go-jose.v2/jwt"
	"gopkg.in/unrolled/secure.v1"
)

type customProxyFactory struct {
	logger  logging.Logger
	factory proxy.Factory
}

// New implements the Factory interface
func (cf customProxyFactory) New(cfg *config.EndpointConfig) (p proxy.Proxy, err error) {
	p, err = cf.factory.New(cfg)
	if err == nil {
		p = proxy.NewLoggingMiddleware(cf.logger, cfg.Endpoint)(p)
	}
	return
}

func main() {
	port := flag.Int("p", 9091, "Port of the service")
	logLevel := flag.String("l", "DEBUG", "Logging level")
	debug := flag.Bool("d", true, "Enable the debug")
	configFile := flag.String("c", "config.json", "Path to the configuration filename")
	flag.Parse()

	parser := config.NewParser()
	config.RoutingPattern = config.BracketsRouterPatternBuilder
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0 {
		serviceConfig.Port = *port
	}
	fmt.Println(serviceConfig.ExtraConfig)
	logger, err := logging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{"localhost:9091"},
		SSLRedirect:           false,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IsDevelopment:         true,
	})

	cors := muxcors.New(serviceConfig.ExtraConfig)

	tokenRejecterFactory := jose.ChainedRejecterFactory([]jose.RejecterFactory{
		jose.RejecterFactoryFunc(func(l logging.Logger, cfg *config.EndpointConfig) jose.Rejecter {
			if r := cel.NewRejecter(l, cfg); r != nil {
				return r
			}
			return jose.FixedRejecter(false)
		}),
	})

	backendFactory := martian.NewBackendFactory(logger, client.DefaultHTTPRequestExecutor(client.NewHTTPClient))
	cfg := gorilla.DefaultConfig(customProxyFactory{logger, proxy.NewDefaultFactory(backendFactory, logger)}, logger)
	cfg.Middlewares = append(cfg.Middlewares, secureMiddleware, cors)
	cfg.HandlerFactory = newHandlerFactory(cfg.HandlerFactory, logger, tokenRejecterFactory)
	routerFactory := mux.NewFactory(cfg)
	routerFactory.New().Run(serviceConfig)
}

func newHandlerFactory(mh mux.HandlerFactory, logger logging.Logger, re jose.ChainedRejecterFactory) mux.HandlerFactory {
	hf := jujumux.HandlerFactory
	hf = muxjose.HandlerFactory(mh, gorillaParamsExtractor, logger, re)
	return hf
}

func gorillaParamsExtractor(r *http.Request) map[string]string {
	params := map[string]string{}
	for key, value := range g.Vars(r) {
		params[strings.Title(key)] = value
	}
	return params
}
