package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"study/pkg/viper"
)

// Pool redis连接池
var Pool *redis.Pool

// InitRedis 初始化redis连接池
func init() {
	r := viper.Conf.Redis
	Pool = &redis.Pool{ // 实例化一个连接池
		MaxIdle:     16,      // 最初的连接数量
		MaxActive:   1000000, // 最大连接数量
		IdleTimeout: 300,     // 连接关闭时间 300秒 (300秒不使用就自动关闭)
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d", r.Address, r.Port))
		},
	}
}
