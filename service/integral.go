package service

/*积分服务*/

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/garyburd/redigo/redis"
	"sTest/entity"
	m "sTest/pkg/mysql"
	r "sTest/pkg/redis"
	"strconv"
	"strings"
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
	if _, err := conn.Do("zincrby", LeaderBoardName, integral, key); err != nil {
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
func GetRankingLimit50() (ranking []entity.LeaderBoardData, err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	// zrevrangebyscore ranking +inf -inf limit 0 50
	rels, err := redis.Strings(conn.Do("zrevrangebyscore", LeaderBoardName, "+inf", "-inf", "WITHSCORES", "limit", 0, 50))
	if err != nil {
		logger.Error(err)
		return
	}
	ranking = make([]entity.LeaderBoardData, 0)
	var item entity.LeaderBoardData
	for i, val := range rels {
		if i%2 == 0 {
			item = entity.LeaderBoardData{}
			splits := strings.Split(val, ":")
			item.UserName = splits[0]
			item.AvatarURL = splits[1]
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
	var account string
	findAccountSQL := "select account from t_account_data where user_id = ?"
	if err = m.DB.Get(&account, findAccountSQL, userID); err != nil {
		logger.Error(err)
		return
	}
	var avatarURL string
	findAvatarSQL := "select avatar_url from t_base_data where user_id = ?"
	if err = m.DB.Get(&avatarURL, findAvatarSQL, userID); err != nil {
		logger.Error(err)
		return
	}

	conn := r.Pool.Get()
	defer conn.Close()
	number, err := redis.Int(conn.Do("ZREVRANK", LeaderBoardName, account+":"+avatarURL))
	if err != nil {
		logger.Error(err)
		return
	}

	score, err := redis.Int(conn.Do("zscore", LeaderBoardName, account+":"+avatarURL))
	if err != nil {
		logger.Error(err)
		return
	}

	ranking = entity.LeaderBoardData{
		Number:    number + 1,
		Integral:  score,
		UserName:  account,
		AvatarURL: avatarURL,
	}

	return
}
