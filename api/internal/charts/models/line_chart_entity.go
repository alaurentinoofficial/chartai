package charts_models

import (
	core_models "github.com/alaurentinoofficial/chartai/internal/core/models"
	"github.com/google/uuid"
)

func init() {
	registerChart(LineChartEntity{})
}

type LineChartEntity struct {
	baseChartEntity
	CategoricalColumnName *string
	ValueColumnName       *string
}

func NewLineChart(
	title string,
	query string,
	databaseId uuid.UUID,
) *LineChartEntity {
	base := newBaseChartEntity(title, query, databaseId)
	return &LineChartEntity{
		baseChartEntity: *base,
	}
}

func NewLineChartFull(
	baseEntity *core_models.BaseEntity,
	title, query string,
	databaseId uuid.UUID,
	categoricalColumnName *string,
	valueColumnName *string,
) *LineChartEntity {
	return &LineChartEntity{
		baseChartEntity:       *newBaseChartEntityWBaseEntity(baseEntity, title, query, databaseId),
		CategoricalColumnName: categoricalColumnName,
		ValueColumnName:       valueColumnName,
	}
}

func (c LineChartEntity) GetType() string {
	return "Line"
}

func (c LineChartEntity) GetSchema() []ChartColumn {
	return categorical_schema
}
