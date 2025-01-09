package core_handlers

import (
	"context"

	"github.com/google/uuid"
)

type HandlerFunc[T any, U any] func(ctx context.Context, request T) (*U, error)

type Authentication struct {
	AccountId uuid.UUID
}

type HandlerAuthenticatedFunc[T any, U any] func(ctx context.Context, auth Authentication, request T) (*U, error)
