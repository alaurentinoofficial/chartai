package postgresdb_sql

import (
	"encoding/json"
	"time"

	charts_models "github.com/alaurentinoofficial/chartai/internal/charts/models"
	core_models "github.com/alaurentinoofficial/chartai/internal/core/models"
	databases_models "github.com/alaurentinoofficial/chartai/internal/databases/models"
)

func (d *Database) ToModel() *databases_models.DatabaseEntity {
	var lastSync *time.Time
	if d.LastSync.Valid {
		lastSync = &d.LastSync.Time
	}

	var schema *[]databases_models.Table
	if d.Schema.Valid {
		var tables []databases_models.Table
		if err := json.Unmarshal(d.Schema.RawMessage, &tables); err == nil { // Handle unmarshaling
			schema = &tables
		}
	}

	return &databases_models.DatabaseEntity{
		BaseEntity: core_models.BaseEntity{
			Id:         d.Id,
			CreatedAt:  d.CreatedAt,
			ModifiedAt: d.ModifiedAt,
			IsArchived: d.IsArchived,
			IsDeleted:  d.IsDeleted,
		},
		Name:             d.Name,
		Type:             databases_models.DatabaseType(d.Type),
		ConnectionString: d.ConnectionString,
		Schema:           schema,
		LastSync:         lastSync,
	}
}

func (d *Chart) ToModel() (charts_models.ChartEntity, error) {
	baseEntity := &core_models.BaseEntity{
		Id:         d.Id,
		CreatedAt:  d.CreatedAt,
		ModifiedAt: d.ModifiedAt,
		IsArchived: d.IsArchived,
		IsDeleted:  d.IsDeleted,
	}

	var categoricalColumnName *string
	if d.CategoricalColumnName.Valid {
		categoricalColumnName = &d.CategoricalColumnName.String
	}

	var valueColumnName *string
	if d.ValueColumnName.Valid {
		valueColumnName = &d.ValueColumnName.String
	}

	chart, err := charts_models.StringToChart(d.Type)
	if err != nil {
		return nil, err
	}

	switch chart.(type) {
	case charts_models.BarChartEntity:
		return charts_models.NewBarChartFull(
			baseEntity,
			d.Title,
			d.Query,
			d.DatabaseId,
			categoricalColumnName,
			valueColumnName,
		), nil
	default:
		return charts_models.NewLineChartFull(
			baseEntity,
			d.Title,
			d.Query,
			d.DatabaseId,
			categoricalColumnName,
			valueColumnName,
		), nil
	}
}
