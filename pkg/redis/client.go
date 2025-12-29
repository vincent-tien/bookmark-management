package redis

import "github.com/redis/go-redis/v9"

// NewClient creates and returns a new Redis client instance.
// It loads configuration from environment variables using the provided prefix.
// Returns a Redis client and an error if configuration loading or client creation fails.
func NewClient(envPrefix string) (*redis.Client, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return client, nil
}
