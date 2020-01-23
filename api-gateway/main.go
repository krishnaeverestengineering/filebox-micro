package main

import (
	"flag"
	"log"
	"os"

	jose "github.com/devopsfaith/krakend-jose"
	muxjose "github.com/devopsfaith/krakend-jose/mux"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/router/gorilla"
	"github.com/devopsfaith/krakend/router/mux"
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

func newHandlerFactory(gf mux.HandlerFactory, pe mux.ParamExtractor, rejecter jose.RejecterFactory, logger logging.Logger) mux.HandlerFactory {
	hf := muxjose.HandlerFactory(gf, pe, logger, rejecter)
	return hf
}

func main() {
	port := flag.Int("p", 9090, "Port of the service")
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

	logger, err := logging.NewLogger(*logLevel, os.Stdout, "[KRAKEND]")
	if err != nil {
		log.Fatal("ERROR:", err.Error())
	}

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{"127.0.0.1:9090"},
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
	//rf := gorilla.DefaultFactory(proxy.DefaultFactory(logger), logger)
	cfg := gorilla.DefaultConfig(customProxyFactory{logger, proxy.DefaultFactory(logger)}, logger)
	cfg.Middlewares = append(cfg.Middlewares, secureMiddleware)
	cfg.HandlerFactory = newHandlerFactory(cfg.HandlerFactory, nil, jose.ChainedRejecterFactory{}, logger)
	routerFactory := mux.NewFactory(cfg)
	routerFactory.New().Run(serviceConfig)
}

/*

func getSessionTokenHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("auth"); err == nil {
			return
		}
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		userId := r.URL.Query().Get("userId")
		if userId == "null" || userId == "" {
			encoder.Encode(TokenResponse{
				Ok: false,
			})
			return
		}
		c, expiresAt, err := GetSessionToken(userId, []byte("sceret"))
		if err != nil {
			encoder.Encode(TokenResponse{
				Ok: false,
			})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "auth",
			Value:   c,
			Expires: expiresAt,
			Path:    "/",
		})
		encoder.Encode(TokenResponse{
			Ok: true,
		})
	})
}
// type MyHandler struct {
// }

// type InvalidTokenResponse struct {
// 	Path string `json:"redirect_path"`
// }

// func (s *MyHandler) Handler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ok := middlewares.TokenValidater(r)
// 		if !ok {
// 			encoder := json.NewEncoder(w)
// 			encoder.SetEscapeHTML(false)
// 			encoder.Encode(&InvalidTokenResponse{
// 				Path: "/login",
// 			})
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

type TokenResponse struct {
	Ok bool `json"ok"`
}

func GetSessionToken(uid string, secret []byte) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	cliams := &models.AuthTokenClaims{
		UserId: uid,
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
}*/
