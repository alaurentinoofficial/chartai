package charts_models

import (
	core_models "github.com/alaurentinoofficial/chartai/internal/core/models"
	"github.com/google/uuid"
)

func init() {
	registerChart(BarChartEntity{})
}

type BarChartEntity struct {
	baseChartEntity
	CategoricalColumnName *string
	ValueColumnName       *string
}

func NewBarChart(
	title string,
	query string,
	databaseId uuid.UUID,
) *BarChartEntity {
	base := newBaseChartEntity(title, query, databaseId)
	return &BarChartEntity{
		baseChartEntity: *base,
	}
}

func NewBarChartFull(
	baseEntity *core_models.BaseEntity,
	title, query string,
	databaseId uuid.UUID,
	categoricalColumnName *string,
	valueColumnName *string,
) *BarChartEntity {
	return &BarChartEntity{
		baseChartEntity:       *newBaseChartEntityWBaseEntity(baseEntity, title, query, databaseId),
		CategoricalColumnName: categoricalColumnName,
		ValueColumnName:       valueColumnName,
	}
}

func (c BarChartEntity) GetType() string {
	return "Bar"
}

func (c BarChartEntity) GetSchema() []ChartColumn {
	return categorical_schema
}
