package service

import (
	"errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/garyburd/redigo/redis"
	"sTest/entity"
	m "sTest/pkg/mysql"
	r "sTest/pkg/redis"
	"sTest/util"
	"strconv"
	"time"
)

func Login(account *entity.AccountData) (ok bool, err error) {
	// select mysql by account
	sqlStr := "select user_id from t_account_data where account = ? and passwd = ?"

	var userID int64
	err = m.DB.Get(&userID, sqlStr, account.Account, account.Passwd)
	if err != nil {
		logger.Error(err)
		err = errors.New("账号或密码错误，请输入正确的账号和密码")
		return false, err
	}

	updateUserStatusSQL := "update t_base_data set is_online = true where user_id = ?"
	if _, err := m.DB.Exec(updateUserStatusSQL, userID); err != nil {
		logger.Error(err)
		return false, err
	}

	return userID > 0, nil
}

func LoginOut(userID int64) (ok bool, err error) {
	offlineTime := time.Now().Unix()
	updateUserStatusSQL := "update t_base_data set is_online = false,offline_time = ? where user_id = ?"
	if _, err := m.DB.Exec(updateUserStatusSQL, offlineTime, userID); err != nil {
		logger.Error(err)
		return false, err
	}

	return userID > 0, nil
}

// TODO(已完成): 注册问题,[[只有不同的设备码才可以注册,相同设备码不能注册]]

func Register(equipmentID string) (v *entity.AccountData, err error) {
	// get random char
	twoChar := util.RandNCharAccount(2)
	passwd := util.RandNCharPasswd(4)

	// get userId
	rConn := r.Pool.Get()
	userID, err := redis.Int(rConn.Do("incr", "userId"))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	// get access
	account := twoChar + strconv.Itoa(userID)

	// baseDataSql := "insert into t_base_data(user_id,nickname,avatarURL,score,isOnline,offlineTime) values(?,?,?,?,?,?)"
	sqlStr := "insert into t_account_data(user_id,passwd,equipment_id,account) values(?,?,?,?)"
	if _, err = m.DB.Exec(sqlStr, userID, passwd, equipmentID, account); err != nil {
		logger.Error(err)
		return nil, err
	}

	accountData := &entity.AccountData{
		UserID:      userID,
		Passwd:      passwd,
		EquipmentID: equipmentID,
		Account:     account,
	}
	return accountData, nil
}
