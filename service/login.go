package service

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/pkg/errors"
	"sTest/entity"
	"sTest/entity/dto"
	"sTest/entity/login_logout"
	"sTest/pkg/auth"
	m "sTest/pkg/mysql"
	"sTest/pkg/response"
	"sTest/repository/cache"
	"sTest/repository/document"
	"sTest/repository/document/mongo_key"
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
		return "", err
	}

	// 改变mongo中用户的在线状态
	if err = document.UpdateElementByUserId(mongo_key.BaseIsOnline, userID, true); err != nil {
		return "", errors.Wrap(err, response.MsgFailed)
	}

	// 检查缓存
	if _, err := cache.GetUserCache(int(userID)); err != nil {
		return "", err
	}

	// 生成token
	token, err = auth.GenerateToken(int(userID))
	if err != nil {
		return "", err
	}
	return
}

func LoginOut(userId int64) (err error) {
	if userId < 1 {
		return errors.New(response.MsgParamsError)
	}
	offlineTime := time.Now().Unix()
	// 改变mongo中用户的在线状态
	if err = document.UpdateTwoElementByUserId(userId, mongo_key.BaseIsOnline, mongo_key.BaseOfflineTime, false, offlineTime); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}
	return nil
}

func Register(param *login_logout.LoginReq) (v *entity.AccountData, err error) {
	userId, err := InitUserGameData()
	if err != nil {
		return nil, errors.Wrap(err, response.MsgInitDataError)
	}

	// store user data in mongo
	item := dto.UserBaseData{
		UserId: userId,

		NickName:    param.NickName,
		AvatarURL:   "default.avatar.URL",
		IsOnline:    false,
		OfflineTime: -1,
		Age:         param.Age,
		Country:     param.Country,
		Integral:    0,
		Friends:     []int64{},
		Applied:     []dto.Applied{},
	}

	if err = document.CreateUser(&item); err != nil {
		return nil, err
	}

	// get random char
	twoChar := util.RandNCharAccount(2)
	passwd := util.RandNCharPasswd(4)

	// get access
	account := twoChar + strconv.Itoa(int(userId))

	sqlStr := "insert into t_account_data(user_id,passwd,equipment_id,account) values(?,?,?,?)"
	if _, err = m.DB.Exec(sqlStr, userId, passwd, param.EquipmentID, account); err != nil {
		logger.Error(err)
		return nil, err
	}

	accountData := &entity.AccountData{
		UserID:      int(userId),
		Passwd:      passwd,
		EquipmentID: param.EquipmentID,
		Account:     account,
	}
	return accountData, nil
}
