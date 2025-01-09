package databases_models

import (
	"time"

	cm "github.com/alaurentinoofficial/chartai/internal/core/models"
	core_models "github.com/alaurentinoofficial/chartai/internal/core/models"
	"github.com/google/uuid"
)

type DatabaseType int

const (
	PostgresDatabaseType DatabaseType = iota
)

type DatabaseEntity struct {
	cm.BaseEntity
	Name             string       `json:"name"`
	Type             DatabaseType `json:"type"`
	ConnectionString string       `json:"connectionString"`
	Schema           *[]Table     `json:"schema"`
	LastSync         *time.Time   `json:"lastSync"`
}

func NewDatabaseEntity(name string, dtype DatabaseType, connectionString string) *DatabaseEntity {
	now := time.Now()
	return &DatabaseEntity{
		BaseEntity: core_models.BaseEntity{
			Id:         uuid.New(),
			CreatedAt:  now,
			ModifiedAt: now,
			IsArchived: false,
			IsDeleted:  false,
		},
		Name:             name,
		Type:             dtype,
		ConnectionString: connectionString,
	}
}

type Table struct {
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
}

type Column struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	IsPrimaryKey bool        `json:"is_primary_key"`
	ForeignKey   *ForeignKey `json:"foreign_key,omitempty"`
}

type ForeignKey struct {
	TargetTable  string `json:"target_table"`
	TargetColumn string `json:"target_column"`
}
