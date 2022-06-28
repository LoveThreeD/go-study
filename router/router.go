package router

import (
	"github.com/gin-gonic/gin"
	"sTest/api"
)

func InitRouter(router *gin.Engine) {
	router.POST("/login", api.Login)
	router.POST("/register", api.Register)
	router.POST("/logout", api.LogOut)

	/*router.Use(middleware.UserAuthMiddleware(
		middleware.AllowPathPrefixSkipper("/login", "/register"),
	))*/
	user := router.Group("/user")
	{
		user.PUT("", api.UpdateUserBaseMessage)
	}

	level := router.Group("/level")
	{
		// 进入关卡
		level.GET("enterLevel", api.EnterLevel)
		// 完成任务
		level.PUT("missionAccomplished", api.FinishTask)
		// 完成关卡
		level.PUT("leave", api.FinishLevel)
	}

	friend := router.Group("/friend")
	{
		// 搜索好友
		friend.GET("search", api.SearchFriend)
		// 申请好友
		friend.POST("application", api.AddApplicationList)
		// 确认好友
		friend.POST("ack", api.Ack)
		// 删除好友
		friend.DELETE("delete", api.DeleteFriend)
		// 推荐好友
		friend.GET("recommendList", api.RecommendFriend)
	}

	// 获取榜单
	router.GET("/integral", api.GetRanking)

}
