package charts_handlers

import (
	charts_models "github.com/alaurentinoofficial/chartai/internal/charts/models"
	"github.com/google/uuid"
)

type ChartResponse struct {
	Id         uuid.UUID `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Query      string    `json:"query"`
	DatabaseId uuid.UUID `json:"databaseId"`
}

func NewChartResponse(chart charts_models.ChartEntity) *ChartResponse {
	return &ChartResponse{
		Id:         chart.GetId(),
		Title:      chart.GetTitle(),
		Type:       chart.GetType(),
		Query:      chart.GetQuery(),
		DatabaseId: chart.GetDatabaseId(),
	}
}

type ChartDataResponse struct {
	Id         uuid.UUID          `json:"id"`
	Type       string             `json:"type"`
	Title      string             `json:"title"`
	Query      string             `json:"query"`
	Data       [](map[string]any) `json:"data"`
	DatabaseId uuid.UUID          `json:"databaseId"`
}

func NewChartDataResponse(chart charts_models.ChartEntity, data [](map[string]any)) *ChartDataResponse {
	return &ChartDataResponse{
		Id:         chart.GetId(),
		Title:      chart.GetTitle(),
		Type:       chart.GetType(),
		Query:      chart.GetQuery(),
		Data:       data,
		DatabaseId: chart.GetDatabaseId(),
	}
}
