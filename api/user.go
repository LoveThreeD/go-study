package api

import (
	"errors"
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"sTest/entity"
	"sTest/pkg/response"
	"sTest/service"
	"strconv"
)

// GetUserByUserID 根据玩家ID获取玩家基础数据
func GetUserByUserID(c *gin.Context) {
	userIDStr := c.Query("userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}

	entity, err := service.GetUserByID(int(userID))
	if err != nil {
		logger.Error(err)
		err = errors.New("用户ID不正确")
		response.ResFailed(c)
	}

	response.ResSuccess(c, entity)
}

// UserInit 当登陆成功，进入“我的家园”，提示玩家给自己设置昵称，并选择性别
// 玩家基础数据表初始化
func UserInit(c *gin.Context) {
	var val entity.BaseData
	if err := c.Bind(&val); err != nil {
		logger.Error(err)
		response.ResFailed(c)
	}

	data, err := service.UserInit(&val)
	if err != nil {
		logger.Error(err)
		response.ResFailed(c)
	}
	response.ResSuccess(c, data)
}
