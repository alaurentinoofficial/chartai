package databases

import (
	databases_handlers "github.com/alaurentinoofficial/chartai/internal/databases/handlers"
	databases_validations "github.com/alaurentinoofficial/chartai/internal/databases/validations"
	"go.uber.org/fx"
)

var DatabasesModule = fx.Options(
	fx.Provide(
		databases_handlers.CreateDatabaseHandler,
		databases_handlers.UpdateDatabaseHandler,
		databases_handlers.GetAllDatabasesHandler,
		databases_handlers.GetDatabaseByIdHandler,
		databases_handlers.DeleteDatabaseHandler,
		databases_handlers.RunQueryOnDatabaseByIdHandler,
	),
	fx.Invoke(
		databases_validations.RegisterValidations,
	),
)
