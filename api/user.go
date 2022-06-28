package api

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"sTest/entity/user"
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
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	entity, err := service.GetUserByID(int(userID))
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusFailedDependency, "")
	}

	response.ResSuccess(c, entity)
}

// UpdateUserBaseMessage 当登陆成功，进入“我的家园”，提示玩家给自己设置昵称，并设置国家
func UpdateUserBaseMessage(c *gin.Context) {
	var val user.ReqUserBase
	if err := c.Bind(&val); err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
	}
	if err := service.UserInit(&val); err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusFailedDependency, "")
	}
	response.ResSuccess(c, true)
}
