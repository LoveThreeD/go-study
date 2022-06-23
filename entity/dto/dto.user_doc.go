package dto

import "sTest/entity"

type UserBaseData struct {
	BaseData entity.BaseData
	UserId   int64 `db:"user_id"`
	Age      int   `db:"age"`
}
