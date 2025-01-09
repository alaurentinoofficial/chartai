package main

import (
	"net/http"

	endpoints "github.com/alaurentinoofficial/chartai/cmd/api/endpoints"
	"github.com/alaurentinoofficial/chartai/cmd/api/server"
	"github.com/alaurentinoofficial/chartai/infrastructure"
	"github.com/alaurentinoofficial/chartai/internal/charts"
	"github.com/alaurentinoofficial/chartai/internal/databases"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	godotenv.Load(".env.dev.local")

	fx.New(
		infrastructure.InfrastructureModule,
		databases.DatabasesModule,
		charts.ChartsModule,
		server.HttpModule,
		endpoints.EndpointsModule,

		fx.Invoke(func(*http.Server) {}),
	).Run()
}
