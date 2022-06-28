package data

import (
	"github.com/asim/go-micro/v3/logger"
	"sTest/entity/dto"
	m "sTest/pkg/mysql"
)

func GetNickNameAndAvatar(userId int) (c *dto.UserCache, err error) {
	sqlStr := "select nickname,avatar_url,is_online from t_base_data where user_id = ?"
	c = &dto.UserCache{}
	if err = m.DB.Get(c, sqlStr, userId); err != nil {
		logger.Error(err)
		return nil, err
	}
	return c, nil
}
