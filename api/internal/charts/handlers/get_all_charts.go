package charts_handlers

import (
	"context"

	charts_repositories "github.com/alaurentinoofficial/chartai/internal/charts/repositories"
	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	databases_handlers "github.com/alaurentinoofficial/chartai/internal/databases/handlers"
)

type GetAllChartsRequest struct{}

func GetAllChartsHandler(
	databaseRepository charts_repositories.ChartRepository,
	runQueryOnDatabase core_handlers.HandlerFunc[databases_handlers.RunQueryOnDatabaseByIdRequest, [](map[string]any)],
) core_handlers.HandlerFunc[GetAllChartsRequest, []ChartDataResponse] {
	return func(ctx context.Context, request GetAllChartsRequest) (*[]ChartDataResponse, error) {
		// Create Chart
		charts, err := databaseRepository.GetAll(ctx)
		if err != nil {
			return nil, err
		}

		responseCharts := make([]ChartDataResponse, len(charts))
		for i, chart := range charts {
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

			responseCharts[i] = *NewChartDataResponse(chart, *data)
		}

		return &responseCharts, nil
	}
}
