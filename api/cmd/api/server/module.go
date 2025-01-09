package server

import "go.uber.org/fx"

var HttpModule = fx.Options(
	fx.Provide(
		NewHTTPServer,
		NewServeMux,
	),
)
