package cache

import (
	"encoding/json"
	"github.com/asim/go-micro/v3/logger"
	redisBase "github.com/garyburd/redigo/redis"
	"sTest/entity/dto"
	"sTest/pkg/redis"
	"sTest/repository/document"
)

const (
	_expire = 2592000
)

func AddUser(userId int, cache *dto.UserCache) (err error) {
	conn := redis.Pool.Get()
	defer conn.Close()

	bytes, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	if _, err := conn.Do("SET", userId, string(bytes), "ex", _expire); err != nil {
		return err
	}
	return nil
}

func ExistsUserCache(userId int) (ok bool, err error) {
	conn := redis.Pool.Get()
	defer conn.Close()

	ok, err = redisBase.Bool(conn.Do("EXISTS", userId))
	if err != nil {
		return false, err
	}
	return ok, nil
}

func GetUserCache(userId int) (cache *dto.UserCache, err error) {
	ok, err := ExistsUserCache(int(userId))
	if err != nil {
		return nil, err
	}
	// 缓存未命中
	if !ok {
		userCache, err := document.SelectUserByUserId(userId)
		if err != nil {
			return nil, err
		}
		if err = AddUser(userId, userCache); err != nil {
			logger.Error(err)
			return nil, err
		}
	}
	// 缓存命中
	cache, err = getUser(userId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return cache, nil
}

func getUser(userId int) (cache *dto.UserCache, err error) {
	conn := redis.Pool.Get()
	defer conn.Close()

	reply, err := redisBase.String(conn.Do("GET", userId))
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(reply), &cache); err != nil {
		return nil, err
	}
	return cache, nil
}
