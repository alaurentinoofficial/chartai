package charts_handlers

import (
	"context"

	charts_repositories "github.com/alaurentinoofficial/chartai/internal/charts/repositories"
	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	databases_handlers "github.com/alaurentinoofficial/chartai/internal/databases/handlers"
	"github.com/google/uuid"
)

type GetChartDataByIdRequest struct {
	Id string `json:"id" validate:"uuid"`
}

func GetChartDataByIdHandler(
	databaseRepository charts_repositories.ChartRepository,
	structValidator *core_validations.StructValidator,
	runQueryOnDatabase core_handlers.HandlerFunc[databases_handlers.RunQueryOnDatabaseByIdRequest, [](map[string]any)],
) core_handlers.HandlerFunc[GetChartDataByIdRequest, ChartDataResponse] {
	return func(ctx context.Context, request GetChartDataByIdRequest) (*ChartDataResponse, error) {
		// Validates
		if err := structValidator.Validate(request); err != nil {
			return nil, err
		}

		// Create Chart
		chartId := uuid.MustParse(request.Id)
		chart, err := databaseRepository.GetById(ctx, chartId)
		if err != nil {
			return nil, err
		}

		data, err := runQueryOnDatabase(
			ctx,
			databases_handlers.RunQueryOnDatabaseByIdRequest{
				Id:    chart.GetDatabaseId().String(),
				Query: chart.GetQuery(),
			},
		)
		if err != nil {
			return nil, err
		}

		return NewChartDataResponse(chart, *data), nil
	}
}
