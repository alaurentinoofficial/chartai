package databases_handlers

import (
	"context"

	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	core_services "github.com/alaurentinoofficial/chartai/internal/core/services"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	databases_models "github.com/alaurentinoofficial/chartai/internal/databases/models"
	databases_repositories "github.com/alaurentinoofficial/chartai/internal/databases/repositories"
	databases_services "github.com/alaurentinoofficial/chartai/internal/databases/services"
	databases_validations "github.com/alaurentinoofficial/chartai/internal/databases/validations"
)

type CreateDatabaseRequest struct {
	Name             string `json:"name" validate:"required,DatabaseName"`
	Type             int32  `json:"type"`
	ConnectionString string `json:"connectionString"`
}

func CreateDatabaseHandler(
	unitOfWork core_services.UnitOfWork,
	databaseRepository databases_repositories.DatabaseRepository,
	structValidator *core_validations.StructValidator,
) core_handlers.HandlerFunc[CreateDatabaseRequest, DatabaseResponse] {
	return func(ctx context.Context, request CreateDatabaseRequest) (*DatabaseResponse, error) {
		// Validates
		if err := structValidator.Validate(request); err != nil {
			return nil, err
		}

		// Begin transactionCtx
		transactionCtx, err := unitOfWork.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer transactionCtx.Rollback()

		databaseConnection, err := databases_services.NewDatabaseConnectionFactory(
			databases_models.DatabaseType(request.Type),
			request.ConnectionString,
		)
		if err != nil {
			return nil, databases_validations.ErrInvalidDatabaseConnectionString
		}

		tables, err := databaseConnection.GetSchema(ctx)
		if err != nil {
			return nil, databases_validations.ErrInvalidDatabaseConnectionString
		}

		database := databases_models.NewDatabaseEntity(
			request.Name,
			databases_models.DatabaseType(request.Type),
			request.ConnectionString,
		)
		database.Schema = tables
		database.LastSync = &database.CreatedAt

		// Create Database
		err = databaseRepository.Create(transactionCtx, database)
		if err != nil {
			return nil, err
		}

		// Commit changes
		err = transactionCtx.Commit()
		if err != nil {
			return nil, err
		}

		return &DatabaseResponse{
			Id:     database.Id,
			Type:   database.Type,
			Name:   database.Name,
			Schema: database.Schema,
		}, nil
	}
}
