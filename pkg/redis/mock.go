package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// InitMockRedis creates and returns a mock Redis client for testing purposes.
// It uses miniredis to provide an in-memory Redis server implementation.
// Returns a Redis client connected to the mock server.
func InitMockRedis(t *testing.T) *redis.Client {
	mock := miniredis.RunT(t)

	return redis.NewClient(&redis.Options{
		Addr: mock.Addr(),
	})
}
