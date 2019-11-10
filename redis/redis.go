package main

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	gRedisPool *redis.Pool
)

// 创建redis连接
func NewRedis(host, password string) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     50,   //最大空闲50
		MaxActive:   100,  //最大连接可用数
		IdleTimeout: 60 * time.Second,
		Wait: true,  //超出最大连接可用后等待
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host, redis.DialPassword(password))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}, nil
}
