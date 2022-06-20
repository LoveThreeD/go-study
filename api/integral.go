package api

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/gin-gonic/gin"
	"sTest/service"
	"strconv"
)

// TODO(待考虑) 返回数据中的排名是否可以不必返回,因为数组有序    (同样积分的会出现什么情况)

func GetRankingLimit50(c *gin.Context) {
	// userID 任务ID
	userIDStr := c.Query("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err})
		return
	}

	selfData, err := service.GetSelfIntegral(userID)
	if err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	data, err := service.GetRankingLimit50()
	if err != nil {
		logger.Error(err)
		c.JSON(200, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	m := make(map[string]interface{})
	m["top50"] = data
	m["selfData"] = selfData

	c.JSON(200, gin.H{"code": 200, "msg": m})
}
