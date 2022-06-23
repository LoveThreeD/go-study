package api

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"sTest/entity/friend_dto"
	"sTest/pkg/response"
	"sTest/service"
)

func SearchFriend(c *gin.Context) {

	// params bind 参数绑定
	var params friend_dto.ReqFriendSearch
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.ErrParam)
		return
	}

	// search user
	users, err := service.SearchUser(&params)
	if err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.MsgFailed)
		return
	}

	response.ResSuccess(c, users)
}
