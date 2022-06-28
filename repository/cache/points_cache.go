package cache

import (
	"github.com/garyburd/redigo/redis"
	r "sTest/pkg/redis"
	"sTest/service"
)

// ExpireDelCache 过期删除缓存
func ExpireDelCache(key string) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()

	if _, err := redis.String(conn.Do("DEL", key)); err != nil {
		return err
	}
	return nil
}

/*
	增加积分
*/
func AddPoints(key string, points int) error {
	conn := r.Pool.Get()
	defer conn.Close()
	if _, err := conn.Do("zincrby", service.GetPointsKey(), points, key); err != nil {
		return err
	}
	return nil
}

/*
	获取排行榜的排序,从高到低
*/
func GetRanks(count int) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	rels, err := redis.Strings(conn.Do("zrevrangebyscore", service.GetPointsKey(), "+inf", "-inf", "WITHSCORES", "limit", 0, count))
	if err != nil {
		return nil, err
	}
	return rels, nil
}

/*
	获取某人的排名  1,2,3,4...
    该处是当某人无积分时返回0
*/
func GetOneRanks(userId int) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	number, err := redis.Int(conn.Do("ZREVRANK", service.GetPointsKey(), userId))
	if err != nil {
		switch err.Error() {
		case "redigo: nil returned":
			number = 0
		default:
			return 0, err
		}
	}
	return number, nil
}

/*
	获取某人的积分
*/
func GetOnePoints(userId int) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	score, err := redis.Int(conn.Do("ZSCORE", service.GetPointsKey(), userId))
	if err != nil {
		switch err.Error() {
		case "redigo: nil returned":
			score = 0
		default:
			return 0, err
		}
	}
	return score, nil
}
