package charts_repositories

import (
	"context"

	charts_models "github.com/alaurentinoofficial/chartai/internal/charts/models"
	"github.com/google/uuid"
)

type ChartRepository interface {
	GetAll(ctx context.Context) ([]charts_models.ChartEntity, error)
	GetById(ctx context.Context, databaseId uuid.UUID) (charts_models.ChartEntity, error)
	Create(ctx context.Context, database charts_models.ChartEntity) error
	Update(ctx context.Context, database charts_models.ChartEntity) error
	SoftDelete(ctx context.Context, databaseId uuid.UUID) error
	Delete(ctx context.Context, databaseId uuid.UUID) error
}
