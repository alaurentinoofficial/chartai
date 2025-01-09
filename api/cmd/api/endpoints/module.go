package endpoints

import "go.uber.org/fx"

var EndpointsModule = fx.Options(
	fx.Invoke(
		RegisterDatabaseEndpoints,
		RegisterChartEndpoints,
	),
)
