package databases_handlers

import (
	"context"

	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	databases_repositories "github.com/alaurentinoofficial/chartai/internal/databases/repositories"
	"github.com/google/uuid"
)

type GetDatabaseByIdRequest struct {
	Id string `json:"id" validate:"uuid"`
}

func GetDatabaseByIdHandler(
	databaseRepository databases_repositories.DatabaseRepository,
) core_handlers.HandlerFunc[GetDatabaseByIdRequest, DatabaseResponse] {
	return func(ctx context.Context, request GetDatabaseByIdRequest) (*DatabaseResponse, error) {
		databaseId := uuid.MustParse(request.Id)

		database, err := databaseRepository.GetById(ctx, databaseId)
		if err != nil {
			return nil, err
		}
		if database == nil {
			return nil, core_validations.ErrNotFound
		}

		return &DatabaseResponse{
			Id:     database.Id,
			Type:   database.Type,
			Name:   database.Name,
			Schema: database.Schema,
		}, nil
	}
}
