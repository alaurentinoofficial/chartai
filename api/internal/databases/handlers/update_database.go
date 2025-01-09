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
	"github.com/google/uuid"
)

type UpdateDatabaseRequest struct {
	Id               string  `json:"id" validate:"uuid"`
	Name             *string `json:"name"`
	Type             *int32  `json:"type"`
	ConnectionString *string `json:"connectionString"`
}

func UpdateDatabaseHandler(
	unitOfWork core_services.UnitOfWork,
	databaseRepository databases_repositories.DatabaseRepository,
	structValidator *core_validations.StructValidator,
) core_handlers.HandlerFunc[UpdateDatabaseRequest, any] {
	return func(ctx context.Context, request UpdateDatabaseRequest) (*any, error) {
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

		databaseId := uuid.MustParse(request.Id)

		database, err := databaseRepository.GetById(transactionCtx, databaseId)
		if err != nil {
			return nil, err
		}

		if request.Type != nil {
			database.Type = databases_models.DatabaseType(*request.Type)
		}

		if request.ConnectionString != nil {
			database.ConnectionString = *request.ConnectionString
		}

		if request.Name != nil {
			database.Name = *request.Name
		}

		if request.Type != nil || request.ConnectionString != nil {
			databaseConnection, err := databases_services.NewDatabaseConnectionFactory(
				database.Type,
				database.ConnectionString,
			)
			if err != nil {
				return nil, databases_validations.ErrInvalidDatabaseConnectionString
			}

			tables, err := databaseConnection.GetSchema(ctx)
			if err != nil {
				return nil, databases_validations.ErrInvalidDatabaseConnectionString
			}

			database.Schema = tables
		}

		// Update Database
		err = databaseRepository.Update(transactionCtx, database)
		if err != nil {
			return nil, err
		}

		// Commit changes
		err = transactionCtx.Commit()
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
