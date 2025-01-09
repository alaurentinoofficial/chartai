package core_services

import (
	"context"
)

type UnitOfWork interface {
	Begin(context context.Context) (Transaction, error)
}

type Transaction interface {
	context.Context
	Commit() error
	Rollback() error
}
