package infrastructure_repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	postgresdb_sql "github.com/alaurentinoofficial/chartai/infrastructure/postgresdb/sql"
	infrastructure_services "github.com/alaurentinoofficial/chartai/infrastructure/services"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	databases_models "github.com/alaurentinoofficial/chartai/internal/databases/models"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type PostgresDatabaseRepository struct {
	queries *postgresdb_sql.Queries
}

func NewPostgresDatabaseRepository(queries *postgresdb_sql.Queries) *PostgresDatabaseRepository {
	return &PostgresDatabaseRepository{queries: queries}
}

func (r *PostgresDatabaseRepository) GetAll(ctx context.Context) ([]databases_models.DatabaseEntity, error) {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	databases, err := queries.GetDatabases(ctx)

	results := make([]databases_models.DatabaseEntity, len(databases))
	for i, w := range databases {
		results[i] = *w.ToModel()
	}

	return results, err
}

func (r *PostgresDatabaseRepository) GetById(ctx context.Context, databaseId uuid.UUID) (*databases_models.DatabaseEntity, error) {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	database, err := queries.GetDatabaseById(ctx, databaseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core_validations.ErrNotFound
		}
		return nil, err
	}

	return database.ToModel(), err
}

func (r *PostgresDatabaseRepository) Create(ctx context.Context, database *databases_models.DatabaseEntity) error {
	var schema pqtype.NullRawMessage
	if database.Schema != nil {
		schemaJson, _ := json.Marshal(database.Schema)
		schema = pqtype.NullRawMessage{
			Valid:      true,
			RawMessage: schemaJson,
		}
	}

	var lastSync sql.NullTime
	if database.LastSync != nil {
		lastSync = sql.NullTime{
			Valid: true,
			Time:  *database.LastSync,
		}
	}

	queries := infrastructure_services.GetQueries(ctx, r.queries)
	err := queries.CreateDatabase(ctx, postgresdb_sql.CreateDatabaseParams{
		Id:               database.Id,
		CreatedAt:        database.CreatedAt,
		ModifiedAt:       database.ModifiedAt,
		IsArchived:       database.IsArchived,
		IsDeleted:        database.IsDeleted,
		Name:             database.Name,
		Type:             int32(database.Type),
		ConnectionString: database.ConnectionString,
		Schema:           schema,
		LastSync:         lastSync,
	})

	return err
}

func (r *PostgresDatabaseRepository) Update(ctx context.Context, database *databases_models.DatabaseEntity) error {
	var schema pqtype.NullRawMessage
	if database.Schema != nil {
		schemaJson, _ := json.Marshal(database.Schema)
		schema = pqtype.NullRawMessage{
			Valid:      true,
			RawMessage: schemaJson,
		}
	}

	var lastSync sql.NullTime
	if database.LastSync != nil {
		lastSync = sql.NullTime{
			Valid: true,
			Time:  *database.LastSync,
		}
	}

	queries := infrastructure_services.GetQueries(ctx, r.queries)
	err := queries.UpdateDatabaseById(ctx, postgresdb_sql.UpdateDatabaseByIdParams{
		Id:               database.Id,
		ModifiedAt:       time.Now(),
		IsArchived:       database.IsArchived,
		IsDeleted:        database.IsDeleted,
		Name:             database.Name,
		Type:             int32(database.Type),
		Schema:           schema,
		ConnectionString: database.ConnectionString,
		LastSync:         lastSync,
	})

	return err
}

func (r *PostgresDatabaseRepository) SoftDelete(ctx context.Context, databaseId uuid.UUID) error {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	err := queries.SoftDeleteDatabaseById(ctx, databaseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return core_validations.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *PostgresDatabaseRepository) Delete(ctx context.Context, databaseId uuid.UUID) error {
	queries := infrastructure_services.GetQueries(ctx, r.queries)
	err := queries.DeleteDatabaseById(ctx, databaseId)
	if err != nil {
		if err == sql.ErrNoRows {
			return core_validations.ErrNotFound
		}
		return err
	}
	return nil
}
