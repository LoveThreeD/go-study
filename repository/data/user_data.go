package data

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/pkg/errors"
	"sTest/entity/dto"
	m "sTest/pkg/mysql"
)

func GetNickNameAndAvatar(userId int) (c *dto.UserCache, err error) {
	sqlStr := "SELECT nickname,avatar_url,is_online FROM base_data WHERE user_id = ?"
	c = &dto.UserCache{}
	if err = m.DB.Get(c, sqlStr, userId); err != nil {
		logger.Error(err)
		return nil, err
	}
	return c, nil
}

/*
	根据用户账户和密码查询用户是否存在
*/
func CheckUserByAccountAndPass(account, pass string) (int64, error) {
	// select mysql by account
	sqlStr := "SELECT user_id FROM account_data WHERE account = ? AND passwd = ?"

	var userId int64
	if err := m.DB.Get(&userId, sqlStr, account, pass); err != nil {
		return 0, err
	}
	if userId < 1 {
		return 0, errors.New("select fail")
	}
	return userId, nil
}

/*
	插入用户数据
*/
func InsertUser(userId int64, account, passwd, equipmentID string) error {
	query := "INSERT INTO account_data(user_id,passwd,equipment_id,account) VALUES(?,?,?,?)"
	result, err := m.DB.Exec(query, userId, passwd, equipmentID, account)
	if err != nil {
		return err
	}
	if _, err := result.LastInsertId(); err != nil {
		return err
	}
	return nil
}
