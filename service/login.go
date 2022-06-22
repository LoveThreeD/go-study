package service

import (
	"errors"
	"github.com/asim/go-micro/v3/logger"
	"sTest/entity"
	"sTest/pkg/auth"
	m "sTest/pkg/mysql"
	"sTest/repository/cache"
	"sTest/util"
	"strconv"
	"time"
)

func Login(account *entity.AccountData) (token string, err error) {
	// select mysql by account
	sqlStr := "select user_id from t_account_data where account = ? and passwd = ?"

	var userID int64
	if err = m.DB.Get(&userID, sqlStr, account.Account, account.Passwd); err != nil {
		logger.Error(err)
		err = errors.New("账号或密码错误，请输入正确的账号和密码")
		return "", err
	}

	updateUserStatusSQL := "update t_base_data set is_online = 1 where user_id = ?"
	if _, err := m.DB.Exec(updateUserStatusSQL, userID); err != nil {
		logger.Error(err)
		return "", err
	}

	//检查缓存
	if _, err := cache.GetUserCache(int(userID)); err != nil {
		logger.Error(err)
		return "", err
	}

	// 生成token
	token, err = auth.GenerateToken(int(userID))
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return
}

func LoginOut(userId int64) (err error) {
	if userId < 1 {
		return errors.New("非法参数")
	}
	offlineTime := time.Now().Unix()
	updateUserStatusSQL := "update t_base_data set is_online = false,offline_time = ? where user_id = ?"
	result, err := m.DB.Exec(updateUserStatusSQL, offlineTime, userId)
	if err != nil {
		logger.Error(err)
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		logger.Error(err)
		return err
	}
	if count <= 0 {
		return errors.New("更新行数为0")
	}
	return nil
}

func Register(equipmentID, nickName string) (v *entity.AccountData, err error) {

	var in entity.BaseData
	in.AvatarURL = "default.avatar.URL"
	in.Score = 0
	in.IsOnline = true
	in.OfflineTime = -1
	in.NickName = nickName

	insertSQL := "insert into t_base_data(nickname,avatar_url,score,is_online,offline_time) values(?,?,?,?,?)"
	result, err := m.DB.Exec(insertSQL, in.NickName, in.AvatarURL, in.Score, in.IsOnline, in.OfflineTime)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// get random char
	twoChar := util.RandNCharAccount(2)
	passwd := util.RandNCharPasswd(4)

	// get access
	account := twoChar + strconv.Itoa(int(userId))

	// baseDataSql := "insert into t_base_data(user_id,nickname,avatarURL,score,isOnline,offlineTime) values(?,?,?,?,?,?)"
	sqlStr := "insert into t_account_data(user_id,passwd,equipment_id,account) values(?,?,?,?)"
	if _, err = m.DB.Exec(sqlStr, userId, passwd, equipmentID, account); err != nil {
		logger.Error(err)
		return nil, err
	}

	accountData := &entity.AccountData{
		UserID:      int(userId),
		Passwd:      passwd,
		EquipmentID: equipmentID,
		Account:     account,
	}
	return accountData, nil
}
