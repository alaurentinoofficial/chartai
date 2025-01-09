package charts

import (
	charts_handlers "github.com/alaurentinoofficial/chartai/internal/charts/handlers"
	"go.uber.org/fx"
)

var ChartsModule = fx.Options(
	fx.Provide(
		charts_handlers.CreateChartHandler,
		charts_handlers.GetAllChartsHandler,
		charts_handlers.GetChartDataByIdHandler,
	),
)
