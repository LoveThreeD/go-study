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

// AddApplicationList 添加朋友，还需经过好友同意才会显示在好友列表
func AddApplicationList(c *gin.Context) {
	// params bind 参数绑定
	var params friend_dto.ReqFriendAdd
	if err := c.Bind(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.ErrParam)
		return
	}
	// search user
	if err := service.AddApplicationList(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.MsgFailed)
		return
	}
	response.ResSuccessWithData(c, response.OK)
}

// AddFriend 从申请列表中添加到好友列表
func AddFriend(c *gin.Context) {
	// params bind 参数绑定
	var params friend_dto.ReqFriendAdd
	if err := c.Bind(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.ErrParam)
		return
	}
	// search user
	if err := service.AddFriendList(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.MsgFailed)
		return
	}
	response.ResSuccessWithData(c, response.OK)
}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {
	// params bind 参数绑定
	var params friend_dto.ReqFriendAdd
	if err := c.Bind(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.ErrParam)
		return
	}
	// search user
	if err := service.DeleteFriendList(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.MsgFailed)
		return
	}
	response.ResSuccessWithData(c, response.OK)
}

// RecommendFriend 好友推荐
func RecommendFriend(c *gin.Context) {
	// params bind 参数绑定
	var params friend_dto.ReqRecommend
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.ErrParam)
		return
	}
	// search user
	friends, err := service.GetRecommendFriends(&params)
	if err != nil {
		logger.Errorf("%+v", err)
		response.ResFailedWithData(c, response.MsgFailed)
		return
	}

	response.ResSuccess(c, friends)
}
