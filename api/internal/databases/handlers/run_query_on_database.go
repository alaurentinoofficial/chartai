package databases_handlers

import (
	"context"

	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	databases_repositories "github.com/alaurentinoofficial/chartai/internal/databases/repositories"
	databases_services "github.com/alaurentinoofficial/chartai/internal/databases/services"
	databases_validations "github.com/alaurentinoofficial/chartai/internal/databases/validations"
	"github.com/google/uuid"
)

type RunQueryOnDatabaseByIdRequest struct {
	Id    string `json:"id" validate:"uuid"`
	Query string `json:"query"`
}

func RunQueryOnDatabaseByIdHandler(
	databaseRepository databases_repositories.DatabaseRepository,
) core_handlers.HandlerFunc[RunQueryOnDatabaseByIdRequest, [](map[string]any)] {
	return func(ctx context.Context, request RunQueryOnDatabaseByIdRequest) (*[](map[string]any), error) {
		databaseId := uuid.MustParse(request.Id)

		database, err := databaseRepository.GetById(ctx, databaseId)
		if err != nil {
			return nil, err
		}
		if database == nil {
			return nil, core_validations.ErrNotFound
		}

		databaseConnection, err := databases_services.NewDatabaseConnectionFactory(
			database.Type,
			database.ConnectionString,
		)
		defer databaseConnection.Close(ctx)
		if err != nil {
			return nil, databases_validations.ErrInvalidDatabaseConnectionString
		}

		return databaseConnection.RunQuery(ctx, request.Query)
	}
}
