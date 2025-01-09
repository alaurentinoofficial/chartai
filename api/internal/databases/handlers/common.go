package databases_handlers

import (
	databases_models "github.com/alaurentinoofficial/chartai/internal/databases/models"
	"github.com/google/uuid"
)

type DatabaseResponse struct {
	Id     uuid.UUID                     `json:"id"`
	Type   databases_models.DatabaseType `json:"type"`
	Name   string                        `json:"name"`
	Schema *[]databases_models.Table     `json:"schema"`
}
