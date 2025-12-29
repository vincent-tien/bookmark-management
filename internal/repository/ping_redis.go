package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// PingRedis defines an interface for checking the connectivity to a Redis database.
// It provides a method to ping the Redis server for health and availability verification.
//
//go:generate mockery --name=PingRedis --filename=ping_redis.go
type PingRedis interface {
	Ping(ctx context.Context) error
}

type pingRedis struct {
	redisClient *redis.Client
}

func NewPingRedis(r *redis.Client) PingRedis {
	return &pingRedis{
		redisClient: r,
	}
}

func (p *pingRedis) Ping(ctx context.Context) error {
	if err := p.redisClient.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}
