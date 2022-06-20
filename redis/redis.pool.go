package redis

import "github.com/garyburd/redigo/redis"

// Pool redis连接池
var Pool *redis.Pool

// InitRedis 初始化redis连接池
func InitRedis() {
	Pool = &redis.Pool{ // 实例化一个连接池
		MaxIdle:     16,      // 最初的连接数量
		MaxActive:   1000000, // 最大连接数量
		IdleTimeout: 300,     // 连接关闭时间 300秒 (300秒不使用就自动关闭)
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "192.168.20.132:6379")
		},
	}
}
