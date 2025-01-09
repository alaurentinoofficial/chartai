package databases_repositories

import (
	"context"

	databases_models "github.com/alaurentinoofficial/chartai/internal/databases/models"
	"github.com/google/uuid"
)

type DatabaseRepository interface {
	GetAll(ctx context.Context) ([]databases_models.DatabaseEntity, error)
	GetById(ctx context.Context, databaseId uuid.UUID) (*databases_models.DatabaseEntity, error)
	Create(ctx context.Context, database *databases_models.DatabaseEntity) error
	Update(ctx context.Context, database *databases_models.DatabaseEntity) error
	SoftDelete(ctx context.Context, databaseId uuid.UUID) error
	Delete(ctx context.Context, databaseId uuid.UUID) error
}
