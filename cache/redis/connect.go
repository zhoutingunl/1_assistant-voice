package myredis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
	//redisHost = "192.168.190.128:6379"
	redisHost = "192.168.171.129:6379"
	redisPass = "testUpLoad"
)

func InitRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 1.打开连接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println("redis dial err:", err)
				return nil, err
			}
			// 2.访问验证
			if _, err := c.Do("AUTH", redisPass); err != nil {
				c.Close()
				fmt.Println("redis auth err:", err)
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				fmt.Println("redis ping err:", err)
				return err
			}
			return err
		},
	}
}

func init() {
	pool = InitRedisPool()
	if pool.TestOnBorrow(pool.Get(), time.Now()) == nil {
		fmt.Println("redis connect success")
	}
}

func InitRedis() *redis.Pool {
	return pool
}
