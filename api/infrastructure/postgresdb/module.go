package postgresdb

import (
	"context"
	"database/sql"
	"log"

	infrastructure_config "github.com/alaurentinoofficial/chartai/infrastructure/configs"
	postgresdb_sql "github.com/alaurentinoofficial/chartai/infrastructure/postgresdb/sql"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func NewPostgresdbConnection(lc fx.Lifecycle, config *infrastructure_config.Config) *sql.DB {
	conn, err := sql.Open("postgres", config.DatabaseConnection)
	if err != nil {
		log.Fatal(err.Error())
	}

	conn.SetMaxOpenConns(30)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			conn.Close()
			return nil
		},
	})
	return conn
}

func NewPostgresdbQueries(conn *sql.DB) *postgresdb_sql.Queries {
	return postgresdb_sql.New(conn)
}
