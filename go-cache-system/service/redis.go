/**
 * @Author: djh
 * @Description:
 * @File:  redis
 * @Version: 1.0.0
 * @Date: 2021/10/12 20:42
 */

package service

import (
	"cache-system/config"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	REDIS_SET    = "Set"
	REDIS_GET    = "Get"
	REDIS_EXPIRE = "expire"
)

var Redis redis.Conn

func NewRedis() {
	var err error
	Redis, err = redis.Dial("tcp", fmt.Sprintf("%s:%d",
		config.Configs.Redis.Host,
		config.Configs.Redis.Port),
		redis.DialDatabase(config.Configs.Redis.Database))
	if err != nil {
		panic(fmt.Sprintf("Redis connect err:%s", err))
	}
}

func NewRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     config.Configs.Redis.MaxIdle, //最大空闲连接数
		IdleTimeout: time.Duration(config.Configs.Redis.IdleTimeout) * time.Second,
		Wait:        true, //超过连接数后是否等待
		Dial: func() (redis.Conn, error) {
			redisUri := fmt.Sprintf("%s:%d", config.Configs.Redis.Host, config.Configs.Redis.Port)
			redisConn, err := redis.Dial("tcp", redisUri, redis.DialDatabase(config.Configs.Redis.Database))
			if err != nil {
				return nil, err
			}
			return redisConn, nil
		},
	}

}

func CloseRedis() {
	Redis.Close()
}
