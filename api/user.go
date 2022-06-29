package api

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"study/entity/user"
	"study/pkg/response"
	"study/service"
)

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
