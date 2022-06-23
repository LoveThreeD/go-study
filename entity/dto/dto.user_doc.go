package dto

import "sTest/entity"

type UserBaseData struct {
	BaseData entity.BaseData
	UserId   int64
	Age      int
	Country  string
	Integral int
	// 申请列表 Application list
	ApplicationList []int64
	// 好友列表 Friend list
	Friends []int64
}
