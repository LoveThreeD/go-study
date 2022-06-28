package timing

import (
	"github.com/asim/go-micro/v3/logger"
	"github.com/robfig/cron/v3"
	"sTest/repository/cache"
	"sTest/service"
)

func init() {
	c := cron.New()
	c.AddFunc("0 0 1 */1 *", func() {
		if err := cache.ExpireDelCache(service.GetLastPointsKey()); err != nil {
			logger.Error(err)
		}
	})
	c.Start()
}
