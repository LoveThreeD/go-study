package service

import (
	"github.com/pkg/errors"
	"sTest/entity"
	"sTest/entity/dto"
	"sTest/entity/login_logout"
	"sTest/pkg/auth"
	"sTest/pkg/response"
	"sTest/repository/cache"
	"sTest/repository/data"
	"sTest/repository/document"
	"sTest/repository/document/mongo_key"
	"sTest/util"
	"strconv"
	"time"
)

func Login(account *entity.AccountData) (string, error) {
	// 检查用户账户密码
	userId, err := data.CheckUserByAccountAndPass(account.Account, account.Passwd)
	if err != nil {
		return "", err
	}

	// 改变mongo中用户的在线状态
	if err = document.UpdateElementByUserId(mongo_key.BaseIsOnline, userId, true); err != nil {
		return "", errors.Wrap(err, response.MsgFailed)
	}

	// 检查缓存
	if _, err := cache.GetUserCache(int(userId)); err != nil {
		return "", err
	}

	// 生成token
	token, err := auth.GenerateToken(int(userId))
	if err != nil {
		return "", err
	}
	return token, nil
}

func LoginOut(userId int64) error {
	if userId < 1 {
		return errors.New(response.MsgParamsError)
	}
	offlineTime := time.Now().Unix()
	// 改变mongo中用户的在线状态
	if err := document.UpdateTwoElementByUserId(userId, mongo_key.BaseIsOnline, mongo_key.BaseOfflineTime, false, offlineTime); err != nil {
		return errors.Wrap(err, response.MsgFailed)
	}
	return nil
}

func Register(param *login_logout.LoginReq) (*entity.AccountData, error) {
	userId, err := CreateUserGameData()
	if err != nil {
		return nil, errors.Wrap(err, response.MsgInitDataError)
	}
	// store user data in mongo
	item := dto.UserBaseData{
		UserId: userId,

		NickName: param.NickName,
		// 注意这里有一个头像属性,默认为空减少存储量,前端会判断没有头像就使用默认的头像
		IsOnline:    false,
		OfflineTime: -1,
		Age:         param.Age,
		Country:     param.Country,
		Points:      0,
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

	if err = data.InsertUser(userId, account, passwd, param.EquipmentID); err != nil {
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
