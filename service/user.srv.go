package service

import (
	"github.com/asim/go-micro/v3/logger"
	"sTest/entity"
	m "sTest/mysql"
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

func UserInit(in *entity.BaseData) (out *entity.BaseData, err error) {
	selectSQL := "select exists(select user_id from t_account_data where user_id = ? limit 1)"
	var ok bool
	if err = m.DB.Get(&ok, selectSQL, in.UserID); err != nil {
		logger.Error(err)
		return nil, err
	}

	in.AvatarURL = "default.avatar.URL"
	in.Score = 0
	in.IsOnline = true
	in.OfflineTime = -1

	insertSQL := "insert into t_base_data(user_id,nickname,avatar_url,score,is_online,offline_time) values(?,?,?,?,?,?)"
	if _, err := m.DB.Exec(insertSQL, in.UserID, in.NickName, in.AvatarURL, in.Score, in.IsOnline, in.OfflineTime); err != nil {
		logger.Error(err)
		return nil, err
	}
	return in, nil
}
