package databases_handlers

import (
	"context"

	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	databases_repositories "github.com/alaurentinoofficial/chartai/internal/databases/repositories"
)

type GetAllDatabasesRequest struct{}

func GetAllDatabasesHandler(
	databaseRepository databases_repositories.DatabaseRepository,
) core_handlers.HandlerFunc[GetAllDatabasesRequest, []DatabaseResponse] {
	return func(ctx context.Context, request GetAllDatabasesRequest) (*[]DatabaseResponse, error) {
		// Create Database
		databases, err := databaseRepository.GetAll(ctx)
		if err != nil {
			return nil, err
		}

		responseDatabases := make([]DatabaseResponse, len(databases))
		for i, database := range databases {
			responseDatabases[i] = DatabaseResponse{
				Id:     database.Id,
				Type:   database.Type,
				Name:   database.Name,
				Schema: database.Schema,
			}
		}

		return &responseDatabases, nil
	}
}
