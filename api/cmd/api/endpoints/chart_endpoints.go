package endpoints

import (
	"github.com/alaurentinoofficial/chartai/cmd/api/server"
	charts_handlers "github.com/alaurentinoofficial/chartai/internal/charts/handlers"
	core_handlers "github.com/alaurentinoofficial/chartai/internal/core/handlers"
	"github.com/gorilla/mux"
)

func RegisterChartEndpoints(
	mux *mux.Router,
	getChartDataById core_handlers.HandlerFunc[charts_handlers.GetChartDataByIdRequest, charts_handlers.ChartDataResponse],
	getAllCharts core_handlers.HandlerFunc[charts_handlers.GetAllChartsRequest, []charts_handlers.ChartDataResponse],
	createChart core_handlers.HandlerFunc[charts_handlers.CreateChartRequest, charts_handlers.ChartResponse],
) {
	mux.HandleFunc("/v1/charts", server.HttpHandler(getAllCharts, false)).Methods("GET")
	mux.HandleFunc("/v1/charts/{id}/data", server.HttpHandler(getChartDataById, false)).Methods("GET")
	mux.HandleFunc("/v1/charts", server.HttpHandler(createChart, true)).Methods("POST")
}
