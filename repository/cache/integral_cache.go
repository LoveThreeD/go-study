package cache

import (
	redisBase "github.com/garyburd/redigo/redis"
	"sTest/pkg/redis"
)

// ExpireDelCache 过期删除缓存
func ExpireDelCache(key string) (err error) {
	conn := redis.Pool.Get()
	defer conn.Close()

	if _, err := redisBase.String(conn.Do("DEL", key)); err != nil {
		return err
	}
	return nil
}
