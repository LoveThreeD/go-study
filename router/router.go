package router

import (
	"github.com/gin-gonic/gin"
	"sTest/api"
)

func InitRouter(router *gin.Engine) {
	router.POST("/login", api.Login)
	router.POST("/register", api.Register)
	router.POST("/logout", api.LogOut)

	user := router.Group("/user")
	{
		// 获取信息
		user.GET("", api.GetUserByUserID)
		// 进入我的家园，设置昵称,设置其它
		user.POST("/init", api.UserInit)
	}

	level := router.Group("/level")
	{
		// 进入关卡
		level.GET("enterLevel", api.EnterLevel)
		// 完成任务
		level.PUT("missionAccomplished", api.MissionAccomplished)
		// 完成关卡
		level.PUT("leave", api.Leave)
	}
	// 获取榜单
	router.GET("/integral", api.GetRankingLimit50)

}
