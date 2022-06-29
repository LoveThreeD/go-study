package service

/* 积分服务*/

import (
	"strconv"
	"study/entity"
	r "study/pkg/redis"
	"study/repository/cache"
)

/*
  增加积分
*/
func addPoints(key string, points int) error {
	return cache.AddPoints(key, points)
}

/*
	获取排行榜,榜单前50名以及自己排名
*/
func GetRanking(count int) ([]*entity.LeaderBoardData, error) {
	ranks, err := cache.GetRanks(count)
	if err != nil {
		return nil, err
	}
	ranking := make([]*entity.LeaderBoardData, len(ranks)/2)

	// 从redis中获取的数据格式为   [k1][v1][k2][v2]....    k为用户的userId  v为用户的排行榜积分
	// 解析缓存排行榜数据,用userId获取用户信息填充
	for i := 0; i < len(ranks); i = i + 2 {
		userId, err := strconv.Atoi(ranks[i])
		if err != nil {
			return nil, err
		}
		userCache, err := cache.GetUserCache(userId)
		if err != nil {
			return nil, err
		}
		item := entity.LeaderBoardData{
			UserName:  userCache.NickName,
			AvatarURL: userCache.AvatarUrl,
			// i = 0,2,4,6...   排名结果就为 = 1,2,3,4...
			Number: (i + 2) / 2,
		}
		if item.Points, err = strconv.Atoi(ranks[i+1]); err != nil {
			return nil, err
		}
		ranking = append(ranking, &item)
	}
	return ranking, nil
}

func GetPointsById(userId int) (*entity.LeaderBoardData, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	// 获取用户数据
	userCache, err := cache.GetUserCache(userId)
	if err != nil {
		return nil, err
	}
	// 获取排名
	number, err := cache.GetOneRanks(userId)
	if err != nil {
		return nil, err
	}
	// 获取积分
	points, err := cache.GetOnePoints(userId)
	if err != nil {
		return nil, err
	}

	ranking := &entity.LeaderBoardData{
		Number:    number + 1,
		Points:    points,
		UserName:  userCache.NickName,
		AvatarURL: userCache.AvatarUrl,
	}
	return ranking, nil
}
