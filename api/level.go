package api

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"sTest/pkg/response"
	"sTest/service"
	"strconv"
)

func EnterLevel(c *gin.Context) {
	// 1.获取关卡ID 用户id
	curLevelStr := c.Query("curLevel")
	idStr := c.Query("userID")

	curLevel, err := strconv.Atoi(curLevelStr)
	if err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}

	// 2.判断关卡是否可进入
	gameData, err := service.EnterLevel(curLevel, userID)
	if err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}
	response.ResSuccess(c, gameData)
}

func MissionAccomplished(c *gin.Context) {
	// userID 任务ID
	userIDStr := c.PostForm("userID")
	taskIDStr := c.PostForm("taskID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}

	if err := service.MissionAccomplished(userID, taskID); err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}
	response.ResSuccess(c, true)
}

func Leave(c *gin.Context) {
	// userID 任务ID
	userIDStr := c.PostForm("userID")
	levelIDStr := c.PostForm("levelID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}
	levelID, err := strconv.Atoi(levelIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}

	if err := service.Leave(userID, levelID); err != nil {
		logger.Error(err)
		response.ResFailed(c)
		return
	}
	response.ResSuccess(c, true)
}
