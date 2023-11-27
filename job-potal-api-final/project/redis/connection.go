package redis

import (
	"fmt"
	"project/config"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	cfg := config.GetConfig()
	//fmt.Sprintf("Addr:%s Password:%s DB:%s",cfg.RadiesConfig.Addr,cfg.RadiesConfig.Password,cfg.RadiesConfig.DB)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(":%s", cfg.RadiesConfig.Addr),     //"localhost:6379",
		Password: fmt.Sprintf(":%s", cfg.RadiesConfig.Password), //"", // no password set
		DB:       cfg.RadiesConfig.DB,
	})
	return rdb

}
