package infrastructure

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"bracelet-ticket-system-be/pkg/xlogger"
)

var redisClient redis.UniversalClient

func redisSetup() {
	logger := xlogger.Logger
	if cfg.RedisClusterMode {
		redisClient = redis.NewClusterClient(
			&redis.ClusterOptions{
				Addrs:    cfg.RedisAddress,
				Password: cfg.RedisPassword,
			},
		)
	} else {
		redisOpt, err := redis.ParseURL(cfg.RedisUrl)
		if err != nil {
			log.Fatalf("Failed to parse Redis URL: %v", err)
		}
		redisClient = redis.NewClient(
			redisOpt,
		)
	}

	ctx := context.Background()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		logger.Error().Err(err).Msgf("Failed connect to  redis: %v", err)
	}
	logger.Info().Msg("Successfully connected to the redis!")
}
