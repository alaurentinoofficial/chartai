package databases_handlers

import (
	"context"

	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	core_services "github.com/alaurentinoofficial/chartai/internal/core/services"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	databases_repositories "github.com/alaurentinoofficial/chartai/internal/databases/repositories"
	"github.com/google/uuid"
)

type DeleteDatabaseRequest struct {
	Id string `json:"id" validate:"uuid"`
}

func DeleteDatabaseHandler(
	unitOfWork core_services.UnitOfWork,
	databaseRepository databases_repositories.DatabaseRepository,
	structValidator *core_validations.StructValidator,
) core_handlers.HandlerFunc[DeleteDatabaseRequest, any] {
	return func(ctx context.Context, request DeleteDatabaseRequest) (*any, error) {
		// Validates
		if err := structValidator.Validate(request); err != nil {
			return nil, err
		}

		databaseId := uuid.MustParse(request.Id)

		// Begin transactionCtx
		transactionCtx, err := unitOfWork.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer transactionCtx.Rollback()

		// Delete Database
		err = databaseRepository.SoftDelete(transactionCtx, databaseId)
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
