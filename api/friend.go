package api

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"sTest/entity/friend_dto"
	"sTest/pkg/response"
	"sTest/service"
)

func SearchFriend(c *gin.Context) {

	var params friend_dto.ReqFriendSearch
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	users, err := service.SearchUser(&params)
	if err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}

	response.ResSuccess(c, users)
}

// RequestFriend 添加朋友，还需经过好友同意才会显示在好友列表
func RequestFriend(c *gin.Context) {

	var params friend_dto.ReqFriendAdd
	if err := c.Bind(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	// 已经在好友列表中的话没必要走下去
	if err := service.ExistsInFriends(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}

	if err := service.RequestFriend(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}
	response.ResSuccess(c, "")
}

// ConfirmFriend 同意/拒绝
func ConfirmFriend(c *gin.Context) {

	var params friend_dto.ReqFriendAdd
	if err := c.Bind(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	// 已经在好友列表中的话没必要走下去
	if err := service.ExistsInFriends(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}

	// 拒绝
	if params.Agree == 0 {
		if err := service.NotPass(&params); err != nil {
			logger.Errorf("%+v", err)
			response.ResFail(c, http.StatusFailedDependency, "")
			return
		}
		response.ResSuccess(c, true)
		return
	}

	// search user
	if err := service.AddFriendList(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}
	response.ResSuccess(c, "")
}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {
	// params bind 参数绑定
	var params friend_dto.ReqFriendAdd
	if err := c.Bind(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}
	// search user
	if err := service.DeleteFriendList(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}
	response.ResSuccess(c, "")
}

// RecommendFriend 好友推荐
func RecommendFriend(c *gin.Context) {
	// params bind 参数绑定
	var params friend_dto.ReqRecommend
	if err := c.ShouldBindQuery(&params); err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}
	// search user
	friends, err := service.GetRecommendFriends(&params)
	if err != nil {
		logger.Errorf("%+v", err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}

	response.ResSuccess(c, friends)
}
