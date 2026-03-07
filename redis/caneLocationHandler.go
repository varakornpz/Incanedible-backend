package myredis

import (
	"encoding/json"
	"time"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)


func PutCaneAddress(caneID string , location LatestCaneLocation) error{
	jsonData, err := json.Marshal(location)

	if err != nil {
		log.Error().Msgf("Failed to marshal location JSON: %v", err)
		return err
	}

	key := "cane_location:" + caneID

	err = RedisClient.Set(Ctx, key, jsonData, 24*time.Hour).Err()
	if err != nil {
		log.Error().Msgf("Failed to save location to Redis for cane %s: %v", caneID, err)
		return err
	}
	return nil
}


func GetLatestLocation(caneID string) (LatestCaneLocation, error) {
	var location LatestCaneLocation
	
	key := "cane_location:" + caneID

	val, err := RedisClient.Get(Ctx, key).Result()

	if err == redis.Nil {
		log.Warn().Msgf("No location data found for cane: %s", caneID)
		return location, err 
	} else if err != nil {
		log.Error().Msgf("Failed to get location from Redis for cane %s: %v", caneID, err)
		return location, err
	}

	
	err = json.Unmarshal([]byte(val), &location)
	if err != nil {
		log.Error().Msgf("Failed to unmarshal location JSON for cane %s: %v", caneID, err)
		return location, err
	}

	return location, nil
}