package infrastructure_services

import (
	core_services "github.com/alaurentinoofficial/chartai/internal/core/services"
	core_validations "github.com/alaurentinoofficial/chartai/internal/core/validations"
	"go.uber.org/fx"
)

var ServicesModule = fx.Options(
	fx.Provide(
		fx.Annotate(NewBcryptHashService, fx.As(new(core_services.HashService))),
		fx.Annotate(NewJwtTokenService, fx.As(new(core_services.TokenService))),
		fx.Annotate(NewPostgreUnitOfWork, fx.As(new(core_services.UnitOfWork))),
		fx.Annotate(NewRedisLocker, fx.As(new(core_services.Locker))),
		fx.Annotate(NewOpenAILlmService, fx.As(new(core_services.LlmService))),
		core_validations.NewStructValidator,
	),
)
