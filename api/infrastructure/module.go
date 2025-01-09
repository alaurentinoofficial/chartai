package infrastructure

import (
	infrastructure_config "github.com/alaurentinoofficial/chartai/infrastructure/configs"
	"github.com/alaurentinoofficial/chartai/infrastructure/postgresdb"
	infrastructure_repositories "github.com/alaurentinoofficial/chartai/infrastructure/postgresdb/repositories"
	infrastructure_redis "github.com/alaurentinoofficial/chartai/infrastructure/redis"
	infrastructure_services "github.com/alaurentinoofficial/chartai/infrastructure/services"
	"go.uber.org/fx"
)

var InfrastructureModule = fx.Options(
	fx.Provide(
		infrastructure_config.NewConfig,
		infrastructure_redis.NewRedisConnection,
		postgresdb.NewPostgresdbConnection,
		postgresdb.NewPostgresdbQueries,
	),

	infrastructure_repositories.RepositoriesModule,
	infrastructure_services.ServicesModule,
)
