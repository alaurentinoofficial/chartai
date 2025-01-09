package endpoints

import (
	"github.com/alaurentinoofficial/chartai/cmd/api/server"
	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	databases_handlers "github.com/alaurentinoofficial/chartai/internal/databases/handlers"
	"github.com/gorilla/mux"
)

func RegisterDatabaseEndpoints(
	mux *mux.Router,
	createDatabase core_handlers.HandlerFunc[databases_handlers.CreateDatabaseRequest, databases_handlers.DatabaseResponse],
	updateDatabase core_handlers.HandlerFunc[databases_handlers.UpdateDatabaseRequest, any],
	getAllDatabases core_handlers.HandlerFunc[databases_handlers.GetAllDatabasesRequest, []databases_handlers.DatabaseResponse],
	getDatabaseById core_handlers.HandlerFunc[databases_handlers.GetDatabaseByIdRequest, databases_handlers.DatabaseResponse],
	deleteDatabase core_handlers.HandlerFunc[databases_handlers.DeleteDatabaseRequest, any],
) {
	mux.HandleFunc("/v1/databases", server.HttpHandler(getAllDatabases, false)).Methods("GET")
	mux.HandleFunc("/v1/databases/{id}", server.HttpHandler(getDatabaseById, false)).Methods("GET")

	mux.HandleFunc("/v1/databases", server.HttpHandler(createDatabase, true)).Methods("POST")
	mux.HandleFunc("/v1/databases/{id}", server.HttpHandler(updateDatabase, true)).Methods("PUT")
	mux.HandleFunc("/v1/databases/{id}", server.HttpHandler(deleteDatabase, true)).Methods("DELETE")
}
