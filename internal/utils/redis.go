package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)


func NewRedis(addr, pw string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,  
		DB:       db,  
	})
	
	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}




