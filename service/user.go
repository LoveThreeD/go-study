package service

import (
	"sTest/entity/user"
	"sTest/repository/document"
	"sTest/repository/document/mongo_key"
)

func UserInit(in *user.ReqUserBase) error {
	// 修改昵称和国家
	if err := document.UpdateTwoElementByUserId(in.UserId, mongo_key.BaseCountry, mongo_key.BaseNickName, in.Country, in.NickName); err != nil {
		return err
	}
	return nil
}
