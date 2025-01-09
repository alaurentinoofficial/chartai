package infrastructure_repositories

import (
	charts_repositories "github.com/alaurentinoofficial/chartai/internal/charts/repositories"
	databases_repositories "github.com/alaurentinoofficial/chartai/internal/databases/repositories"
	"go.uber.org/fx"
)

var RepositoriesModule = fx.Options(
	fx.Provide(
		fx.Annotate(NewPostgresDatabaseRepository, fx.As(new(databases_repositories.DatabaseRepository))),
		fx.Annotate(NewPostgresChartRepository, fx.As(new(charts_repositories.ChartRepository))),
	),
)
