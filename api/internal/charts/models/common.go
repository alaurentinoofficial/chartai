package charts_models

import (
	"time"

	cm "github.com/alaurentinoofficial/chartai/internal/core/models"
	core_models "github.com/alaurentinoofficial/chartai/internal/core/models"
	"github.com/google/uuid"
)

var registradedCharts map[string]ChartEntity = map[string]ChartEntity{}

func registerChart(chart ChartEntity) {
	registradedCharts[chart.GetType()] = chart
}

func GetRegistredCharts() map[string]ChartEntity {
	return registradedCharts
}

func GetRegistredChartsLen() int {
	return len(registradedCharts)
}

type ChartEntity interface {
	GetId() uuid.UUID
	GetCreatedAt() time.Time
	GetModifiedAt() time.Time
	GetIsDeleted() bool
	GetIsArchived() bool
	GetTitle() string
	GetType() string
	GetSchema() []ChartColumn
	GetQuery() string
	GetDatabaseId() uuid.UUID

	SetId(uuid.UUID)
	SetCreatedAt(time.Time)
	SetModifiedAt(time.Time)
	SetTitle(string)
	SetQuery(string)
	SetDatabaseId(uuid.UUID)
}

type baseChartEntity struct {
	cm.BaseEntity
	Title      string    `json:"title"`
	Query      string    `json:"query"`
	DatabaseId uuid.UUID `json:"databaseId"`
}

func newBaseChartEntity(
	title, query string,
	databaseId uuid.UUID,
) *baseChartEntity {
	now := time.Now()
	return &baseChartEntity{
		BaseEntity: core_models.BaseEntity{
			Id:         uuid.New(),
			CreatedAt:  now,
			ModifiedAt: now,
			IsArchived: false,
			IsDeleted:  false,
		},
		DatabaseId: databaseId,
		Title:      title,
		Query:      query,
	}
}

func newBaseChartEntityWBaseEntity(
	baseEntity *core_models.BaseEntity,
	title, query string,
	databaseId uuid.UUID,
) *baseChartEntity {
	return &baseChartEntity{
		BaseEntity: *baseEntity,
		DatabaseId: databaseId,
		Title:      title,
		Query:      query,
	}
}

func (c baseChartEntity) GetId() uuid.UUID {
	return c.BaseEntity.Id
}

func (c baseChartEntity) GetCreatedAt() time.Time {
	return c.BaseEntity.CreatedAt
}

func (c baseChartEntity) GetModifiedAt() time.Time {
	return c.BaseEntity.ModifiedAt
}

func (c baseChartEntity) GetIsArchived() bool {
	return c.BaseEntity.IsArchived
}

func (c baseChartEntity) GetIsDeleted() bool {
	return c.BaseEntity.IsDeleted
}

func (c baseChartEntity) GetTitle() string {
	return c.Title
}

func (c baseChartEntity) GetQuery() string {
	return c.Query
}

func (c baseChartEntity) GetDatabaseId() uuid.UUID {
	return c.DatabaseId
}

func (b baseChartEntity) SetId(id uuid.UUID) {
	b.Id = id
}

func (b baseChartEntity) SetCreatedAt(createdAt time.Time) {
	b.CreatedAt = createdAt
}

func (b baseChartEntity) SetModifiedAt(modifiedAt time.Time) {
	b.ModifiedAt = modifiedAt
}

func (c baseChartEntity) SetIsArchived(value bool) {
	c.BaseEntity.IsArchived = value
}

func (c baseChartEntity) SetIsDeleted(value bool) {
	c.BaseEntity.IsDeleted = value
}

func (b baseChartEntity) SetTitle(title string) {
	b.Title = title
}

func (b baseChartEntity) SetQuery(query string) {
	b.Query = query
}

func (b baseChartEntity) SetDatabaseId(databaseId uuid.UUID) {
	b.DatabaseId = databaseId
}
