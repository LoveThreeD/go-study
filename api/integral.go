package api

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"study/pkg/response"
	"study/service"
)

// XXX(待考虑) 返回数据中的排名是否可以不必返回,因为数组有序    (同样积分的会出现什么情况)

func GetRanking(c *gin.Context) {
	// userID 任务ID
	userIDStr := c.Query("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusPreconditionFailed, "")
		return
	}

	selfData, err := service.GetPointsById(userID)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}

	data, err := service.GetRanking(50)
	if err != nil {
		logger.Error(err)
		response.ResFail(c, http.StatusFailedDependency, "")
		return
	}

	m := make(map[string]interface{})
	m["top50"] = data
	m["selfData"] = selfData
	response.ResSuccess(c, m)
}
