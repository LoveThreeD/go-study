package service

import (
	"github.com/asim/go-micro/v3/logger"
	"sTest/entity"
	"sTest/entity/user"
	m "sTest/pkg/mysql"
	"sTest/repository/document"
	"sTest/repository/document/mongo_key"
)

func GetUserByID(userID int) (e *entity.BaseData, err error) {
	sqlStr := "select user_id,nickname,avatar_url,score,is_online,offline_time from t_base_data where user_id = ?"

	e = &entity.BaseData{}
	err = m.DB.Get(e, sqlStr, userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return e, nil
}

func UserInit(in *user.ReqUserBase) (err error) {

	// 修改昵称和国家
	if err = document.UpdateTwoElementByUserId(in.UserId, mongo_key.BaseCountry, mongo_key.BaseNickName, in.Country, in.NickName); err != nil {
		return err
	}

	return nil
}
