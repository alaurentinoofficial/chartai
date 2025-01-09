package core_models

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	Id         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
	IsArchived bool      `json:"archived"`
	IsDeleted  bool      `json:"deleted"`
}
