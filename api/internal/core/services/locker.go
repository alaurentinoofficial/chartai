package core_services

import (
	"context"
	"time"
)

type Locker interface {
	Acquire(
		ctx context.Context,
		name string,
		expiration time.Duration,
	) (Lock, error)
}

type Lock interface {
	Release()
}
