package service

/*积分服务*/

import (
	"fmt"
	"github.com/asim/go-micro/v3/logger"
	"github.com/garyburd/redigo/redis"
	"sTest/entity"
	r "sTest/pkg/redis"
	"sTest/repository/cache"
	"strconv"
	"time"
)

const (
	LeaderBoardName = "ranking"
)

/*增加积分
  1.add task integral
  2.add level integral
*/

// AddIntegral
// key userAccount:avatar url
func AddIntegral(key string, integral int) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	if _, err := conn.Do("zincrby", GetIntegralKey(), integral, key); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

/*
	get ranking
      1. 50                   zrevrangebyscore ranking +inf -inf limit 0 50
      2. self                 zscore ranking c
*/
func GetRanking(count int) (ranking []entity.LeaderBoardData, err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	// zrevrangebyscore ranking +inf -inf limit 0 50
	rels, err := redis.Strings(conn.Do("zrevrangebyscore", GetIntegralKey(), "+inf", "-inf", "WITHSCORES", "limit", 0, count))
	if err != nil {
		logger.Error(err)
		return
	}
	ranking = make([]entity.LeaderBoardData, 0)
	var item entity.LeaderBoardData
	for i, val := range rels {
		var userId int
		if i%2 == 0 {
			userId, err = strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			userCache, err := cache.GetUserCache(userId)
			if err != nil {
				return nil, err
			}

			item = entity.LeaderBoardData{}
			item.UserName = userCache.NickName
			item.AvatarURL = userCache.AvatarUrl
		} else {
			item.Number = (i + 1) / 2
			if item.Integral, err = strconv.Atoi(val); err != nil {
				logger.Error(err)
				return
			}
			ranking = append(ranking, item)
		}
	}
	return
}

func GetSelfIntegral(userID int) (ranking entity.LeaderBoardData, err error) {
	conn := r.Pool.Get()
	defer conn.Close()

	userCache, err := cache.GetUserCache(userID)
	if err != nil {
		logger.Error(err)
		return
	}

	number, err := redis.Int(conn.Do("ZREVRANK", GetIntegralKey(), userID))
	if err != nil {
		switch err.Error() {
		case "redigo: nil returned":
			number = 0
		default:
			logger.Error(err)
			return
		}
	}

	score, err := redis.Int(conn.Do("zscore", GetIntegralKey(), userID))
	if err != nil {
		switch err.Error() {
		case "redigo: nil returned":
			score = 0
		default:
			logger.Error(err)
			return
		}
	}

	ranking = entity.LeaderBoardData{
		Number:    number + 1,
		Integral:  score,
		UserName:  userCache.NickName,
		AvatarURL: userCache.AvatarUrl,
	}
	return ranking, nil
}

func GetIntegralKey() string {
	month := time.Now().Month().String()
	return fmt.Sprintf("%s:%d", LeaderBoardName, month)
}

func GetLastIntegralKey() string {
	month := time.Now().Month()
	return fmt.Sprintf("%s:%d", LeaderBoardName, month-1)
}
