package api

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"study/pkg/response"
	"study/service"
)

func EnterLevel(c *gin.Context) {
	// 1.获取关卡ID 用户id
	curLevelStr := c.Query("curLevel")
	idStr := c.Query("userID")

	curLevel, err := strconv.Atoi(curLevelStr)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	// 2.判断关卡是否可进入
	gameData, err := service.EnterLevel(curLevel, userID)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}
	response.ResSuccess(c, gameData)
}

func FinishTask(c *gin.Context) {
	// userID 任务ID
	userIDStr := c.PostForm("userID")
	taskIDStr := c.PostForm("taskID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	if err := service.FinishTask(userID, taskID); err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}
	response.ResSuccess(c, true)
}

func FinishLevel(c *gin.Context) {
	// userID 任务ID
	userIDStr := c.PostForm("userID")
	levelIDStr := c.PostForm("levelID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}
	levelID, err := strconv.Atoi(levelIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	if err := service.FinishLevel(userID, levelID); err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}
	response.ResSuccess(c, true)
}
