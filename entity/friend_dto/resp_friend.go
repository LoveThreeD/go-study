package friend_dto

import "sTest/entity/dto"

type RespFriendRecommend struct {
	UserId      int64
	Age         int
	Country     string
	Points      int
	NickName    string
	AvatarURL   string
	IsOnline    bool
	OfflineTime int64
}

func ConventRespFriendRecommend(item *dto.UserBaseData) RespFriendRecommend {
	result := RespFriendRecommend{}
	result.UserId = item.UserId
	result.Age = item.Age
	result.Country = item.Country
	result.Points = item.Points
	result.NickName = item.NickName
	result.AvatarURL = item.AvatarURL
	result.IsOnline = item.IsOnline
	result.OfflineTime = item.OfflineTime
	return result
}
