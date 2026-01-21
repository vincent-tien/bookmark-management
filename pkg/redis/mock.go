package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

// InitMockRedis creates and returns a mock Redis client for testing purposes.
// It uses miniredis to provide an in-memory Redis server implementation.
// Returns a Redis client connected to the mock server.
// Note: This requires network capabilities to bind to a TCP port.
// In Docker builds, ensure the build has network access (--network=host or --privileged).
func InitMockRedis(t *testing.T) *redis.Client {
	mock := miniredis.RunT(t)

	return redis.NewClient(&redis.Options{
		Addr: mock.Addr(),
	})
}
