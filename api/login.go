package api

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"sTest/entity"
	"sTest/entity/login_logout"
	"sTest/pkg/response"
	"sTest/service"
	"strconv"
)

// Login access+passwd
// ⽣成账号（8 =字⺟2 + ⽤户id）和密码（4）
func Login(c *gin.Context) {
	// get query params
	account := entity.AccountData{}
	if err := c.Bind(&account); err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}
	// get flag
	token, err := service.Login(&account)
	if err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusForbidden, "")
		return
	}
	response.ResSuccess(c, token)
}

func LogOut(c *gin.Context) {
	// get query params
	var userID int
	userIDStr := c.PostForm("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	// get flag
	err = service.LoginOut(int64(userID))
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusNotAcceptable, "")
		return
	}

	response.ResSuccess(c, "logout success")
}

// Register ⽤户注册（唯⼀id，机器码）⽣成账号（8 =字⺟2 + ⽤户id）和密码（4）
// 唯一id是自动生成的6位
// 机器码是前端传入的
// 账号是自动生成的8位
// 密码是4位.可以是固定的
func Register(c *gin.Context) {
	// get query params
	var param login_logout.LoginReq
	if err := c.Bind(&param); err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	equipmentID := param.EquipmentID
	nickName := param.NickName

	if equipmentID == "" || nickName == "" {
		err := errors.New("接收参数不正确,空参数")
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	accountData, err := service.Register(&param)
	if err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusNotAcceptable, "")
		return
	}

	response.ResSuccess(c, accountData)
}
