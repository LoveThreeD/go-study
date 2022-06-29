package timing

import (
	"study/repository/cache"

	"github.com/asim/go-micro/v3/logger"
	"github.com/robfig/cron/v3"
)

func init() {
	c := cron.New()
	if _, err := c.AddFunc("0 0 1 */1 *", func() {
		if err := cache.ExpireDelCache(cache.GetLastPointsKey()); err != nil {
			logger.Error(err)
		}
	}); err != nil {
		logger.Error(err)
	}
	c.Start()
}
