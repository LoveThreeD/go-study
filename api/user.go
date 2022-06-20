package api

import (
	"errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"sTest/entity"
	"sTest/service"
	"strconv"
)

// GetUserByUserID 根据玩家ID获取玩家基础数据
func GetUserByUserID(c *gin.Context) {
	userIDStr := c.Query("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": "参数不正确"})
		return
	}

	entity, err := service.GetUserByID(int(userID))
	if err != nil {
		logger.Error(err)
		err = errors.New("用户ID不正确")
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": entity})
}

// UserInit 当登陆成功，进入“我的家园”，提示玩家给自己设置昵称，并选择性别
// 玩家基础数据表初始化
func UserInit(c *gin.Context) {
	var val entity.BaseData
	if err := c.Bind(&val); err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err})
		return
	}

	data, err := service.UserInit(&val)
	if err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": data})
}
