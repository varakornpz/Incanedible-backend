package myredis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/varakornpz/providers"
)

type LatestCaneLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     providers.AppConf.REDISADDRESS,
		Password: "", // ถ้ามีรหัสผ่านค่อยเพิ่มจาก config
		DB:       0,
	})
    _, err := RedisClient.Ping(Ctx).Result()

	if err != nil {
		log.Fatal().Msgf("Failed to connect to Redis: %v", err)
	}

	log.Info().Msg("Redis Initialized successfully")
}
