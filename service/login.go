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

	/*updateUserStatusSQL := "update t_base_data set is_online = 1 where user_id = ?"
	if _, err := m.DB.Exec(updateUserStatusSQL, userID); err != nil {
		logger.Error(err)
		return "", err
	}*/

	//改变mongo中用户的在线状态
	if err = document.UpdateElementByUserId(_isOnline, userID, true); err != nil {
		return "", errors.Wrap(err, response.MsgFailed)
	}

	//检查缓存
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
		return errors.New("非法参数")
	}
	offlineTime := time.Now().Unix()
	//改变mongo中用户的在线状态
	if err = document.UpdateTwoElementByUserId(userId, _isOnline, _offlineTime, false, offlineTime); err != nil {
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
		BaseData: entity.BaseData{
			NickName:    param.NickName,
			AvatarURL:   "default.avatar.URL",
			Score:       0,
			IsOnline:    false,
			OfflineTime: -1,
		},
		Age:      param.Age,
		Country:  param.Country,
		Integral: 0,
		/*ApplicationList: []int64{},
		AlreadyAppliedList: []int64{},
		NoPassList: []int64{},
		NoPassRecord: []int64{},*/
		Friends: []int64{},
		Applied: []dto.Applied{},
	}

	if err = document.CreateUser(&item); err != nil {
		return nil, err
	}

	// get random char
	twoChar := util.RandNCharAccount(2)
	passwd := util.RandNCharPasswd(4)

	// get access
	account := twoChar + strconv.Itoa(int(userId))

	// baseDataSql := "insert into t_base_data(user_id,nickname,avatarURL,score,isOnline,offlineTime) values(?,?,?,?,?,?)"
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
