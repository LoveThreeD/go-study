package api

import (
	"errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"sTest/entity"
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
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// get flag
	ok, err := service.Login(&account)
	if err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	if !ok {
		logger.Warnf("账户或者密码错误: mac地址....")
		c.JSON(200, gin.H{"code": 401, "msg": "账号或密码错误，请输入正确的账号和密码"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "Login Success"})
}

func LogOut(c *gin.Context) {
	// get query params
	var userID int
	userIDStr := c.PostForm("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err})
		return
	}

	// get flag
	ok, err := service.LoginOut(int64(userID))
	if err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err})
		return
	}

	if !ok {
		logger.Warnf("登出失败!")
		c.JSON(200, gin.H{"code": 401, "msg": "登出失败!"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "LogOut Success"})
}

// Register ⽤户注册（唯⼀id，机器码）⽣成账号（8 =字⺟2 + ⽤户id）和密码（4）
// 唯一id是自动生成的6位
// 机器码是前端传入的
// 账号是自动生成的8位
// 密码是4位.可以是固定的
func Register(c *gin.Context) {
	// get query params
	equipmentID := c.PostForm("equipmentId")
	if equipmentID == "" {
		err := errors.New("接收参数不正确,空参数")
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// init t_account_data and t_base_data
	accountData, err := service.Register(equipmentID)
	if err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// init game data
	if _, err = service.InitUserGameData(accountData.UserID); err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": accountData})
}
