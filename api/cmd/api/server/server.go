package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	ic "github.com/alaurentinoofficial/chartai/infrastructure/configs"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle, config *ic.Config, mux *mux.Router) *http.Server {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                     // All origins
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, // Allowing only get, just an example
	})
	srv := &http.Server{Addr: fmt.Sprintf("localhost:%s", config.Port), Handler: c.Handler(mux)}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func NewServeMux() *mux.Router {
	return mux.NewRouter()
}
