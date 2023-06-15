package limiter

import (
	"github.com/paceew/go-redisson/pkg/log"
	"time"
)

type GCounter struct {
	glimiter       *GLimiter
	key            string
	limitThreshold int64
	limitTime      time.Duration
	logger         log.FieldsLogger
}

// Limit 限流，GCounter计数限流器限流
//
// 增加allow计数，在limitTime时间内计数达到limitThreshold阈值则触发限流，limitTime时间后会重置
func (c *GCounter) Limit(n int64) (allow bool, err error) {
	key := c.key
	var current int64 = 0
	current, err = c.glimiter.redisOper.IncrBy(key, n)
	c.logger.Tracef("Incr %s %d current:%d, error:%v ", key, n, current, err)
	if err != nil {
		return
	}
	if current == n {
		if _, err = c.glimiter.redisOper.PExpire(key, c.limitTime); err != nil {
			c.logger.Tracef("PExpire %s %s, error:%v ", key, c.limitTime, err)
			// 失败重试一次
			if _, err = c.glimiter.redisOper.PExpire(key, c.limitTime); err != nil {
				c.logger.Errorf("PExpire %s %s,both times error:%s", key, c.limitTime, err.Error())
			}
		}
	}

	if current > c.limitThreshold {
		// 触发限流检查ttl，如果不存在ttl则重新设置
		if ttl, err := c.glimiter.redisOper.PTTL(key); err == nil && ttl == -1 {
			c.logger.Warnf("%s limit trigger,threshold:%d current:%d,but not have ttl,reset expire:%s", key, current, c.limitTime)
			c.glimiter.redisOper.PExpire(key, c.limitTime)
		}
		return false, err
	}

	return true, err
}
