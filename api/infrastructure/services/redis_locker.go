package infrastructure_services

import (
	"time"

	"github.com/alaurentinoofficial/chartai/internal/core/services"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type redisLocker struct {
	client *redis.Client
}

func NewRedisLocker(client *redis.Client) *redisLocker {
	return &redisLocker{client: client}
}

func (l *redisLocker) Acquire(
	ctx context.Context,
	name string,
	expiration time.Duration,
) (core_services.Lock, error) {
	resultChan := make(chan *RedisLock)
	errorChan := make(chan error)

	go func() {
		for {
			wasAcquired, err := l.client.SetNX(ctx, name, "1", expiration).Result()
			if err != nil {
				errorChan <- err
				return
			}

			// If lock was acquired, send it to result channel and exit
			if wasAcquired {
				resultChan <- &RedisLock{
					ctx:    ctx,
					client: l.client,
					name:   name,
				}
				return
			}

			// Wait for the retry interval or context timeout before trying again
			select {
			case <-time.After(time.Millisecond * 50):
				// Continue the loop for another attempt
			case <-ctx.Done():
				// If context timeout is reached, exit
				errorChan <- ctx.Err()
				return
			}
		}
	}()

	// Listen for either a successful lock acquisition or an error/timeout
	select {
	case lock := <-resultChan:
		return lock, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		// If the original context is canceled
		return nil, ctx.Err()
	}
}

type RedisLock struct {
	ctx    context.Context
	client *redis.Client
	name   string
}

func (l *RedisLock) Release() {
	l.client.Del(context.Background(), l.name)
}
