package infrastructure_repositories

import (
	"context"
	"database/sql"

	postgresdb_sql "github.com/alaurentinoofficial/chartai/infrastructure/postgresdb/sql"
	infrastructure_services "github.com/alaurentinoofficial/chartai/infrastructure/services"
	charts_models "github.com/alaurentinoofficial/chartai/internal/charts/models"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	"github.com/google/uuid"
)

type PostgresChartRepository struct {
	queries *postgresdb_sql.Queries
}

func NewPostgresChartRepository(queries *postgresdb_sql.Queries) *PostgresChartRepository {
	return &PostgresChartRepository{queries: queries}
}

func (r *PostgresChartRepository) GetAll(ctx context.Context) ([]charts_models.ChartEntity, error) {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	charts, err := queries.GetCharts(ctx)

	results := make([]charts_models.ChartEntity, len(charts))
	for i, w := range charts {
		results[i], err = w.ToModel()
		if err != nil {
			return nil, err
		}
	}

	return results, err
}

func (r *PostgresChartRepository) GetById(ctx context.Context, databaseId uuid.UUID) (charts_models.ChartEntity, error) {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	database, err := queries.GetChartById(ctx, databaseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core_validations.ErrNotFound
		}
		return nil, err
	}

	return database.ToModel()
}

func (r *PostgresChartRepository) Create(ctx context.Context, database charts_models.ChartEntity) error {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	request := postgresdb_sql.CreateChartParams{
		Id:         database.GetId(),
		CreatedAt:  database.GetCreatedAt(),
		ModifiedAt: database.GetModifiedAt(),
		IsArchived: database.GetIsArchived(),
		IsDeleted:  database.GetIsDeleted(),
		Title:      database.GetTitle(),
		Type:       database.GetType(),
		Query:      database.GetQuery(),
		DatabaseId: database.GetDatabaseId(),
	}

	if chart, ok := database.(charts_models.BarChartEntity); ok {
		var valueColumnName sql.NullString
		if chart.ValueColumnName != nil {
			valueColumnName = sql.NullString{
				Valid:  true,
				String: *chart.ValueColumnName,
			}
		}

		var categoricalColumnName sql.NullString
		if chart.CategoricalColumnName != nil {
			categoricalColumnName = sql.NullString{
				Valid:  true,
				String: *chart.CategoricalColumnName,
			}
		}

		request.ValueColumnName = valueColumnName
		request.CategoricalColumnName = categoricalColumnName
	}

	if chart, ok := database.(charts_models.LineChartEntity); ok {
		var valueColumnName sql.NullString
		if chart.ValueColumnName != nil {
			valueColumnName = sql.NullString{
				Valid:  true,
				String: *chart.ValueColumnName,
			}
		}

		var categoricalColumnName sql.NullString
		if chart.CategoricalColumnName != nil {
			categoricalColumnName = sql.NullString{
				Valid:  true,
				String: *chart.CategoricalColumnName,
			}
		}

		request.ValueColumnName = valueColumnName
		request.CategoricalColumnName = categoricalColumnName
	}
	err := queries.CreateChart(ctx, request)

	return err
}

func (r *PostgresChartRepository) Update(ctx context.Context, database charts_models.ChartEntity) error {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	request := postgresdb_sql.UpdateChartByIdParams{
		Id:         database.GetId(),
		ModifiedAt: database.GetModifiedAt(),
		IsArchived: database.GetIsArchived(),
		IsDeleted:  database.GetIsDeleted(),
		Title:      database.GetTitle(),
		Type:       database.GetType(),
		Query:      database.GetQuery(),
	}

	if chart, ok := database.(charts_models.BarChartEntity); ok {
		var valueColumnName sql.NullString
		if chart.ValueColumnName != nil {
			valueColumnName = sql.NullString{
				Valid:  true,
				String: *chart.ValueColumnName,
			}
		}

		var categoricalColumnName sql.NullString
		if chart.CategoricalColumnName != nil {
			categoricalColumnName = sql.NullString{
				Valid:  true,
				String: *chart.CategoricalColumnName,
			}
		}

		request.ValueColumnName = valueColumnName
		request.CategoricalColumnName = categoricalColumnName
	}

	if chart, ok := database.(charts_models.LineChartEntity); ok {
		var valueColumnName sql.NullString
		if chart.ValueColumnName != nil {
			valueColumnName = sql.NullString{
				Valid:  true,
				String: *chart.ValueColumnName,
			}
		}

		var categoricalColumnName sql.NullString
		if chart.CategoricalColumnName != nil {
			categoricalColumnName = sql.NullString{
				Valid:  true,
				String: *chart.CategoricalColumnName,
			}
		}

		request.ValueColumnName = valueColumnName
		request.CategoricalColumnName = categoricalColumnName
	}

	err := queries.UpdateChartById(ctx, request)

	return err
}

func (r *PostgresChartRepository) SoftDelete(ctx context.Context, databaseId uuid.UUID) error {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	err := queries.SoftDeleteChartById(ctx, databaseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return core_validations.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *PostgresChartRepository) Delete(ctx context.Context, databaseId uuid.UUID) error {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	err := queries.DeleteChartById(ctx, databaseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return core_validations.ErrNotFound
		}
		return err
	}
	return nil
}
