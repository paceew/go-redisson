package limiter

import (
	"github.com/paceew/go-redisson/pkg/log"
	"time"
)

type GLimiter struct {
	redisOper RedisLimiterOperate
	opts      *GLimiterOptions
}

func (g *GLimiter) getKey(limitertype, key string) string {
	if !g.opts.unPrefix {
		return g.opts.prefix + "_" + limitertype + "_" + key
	}
	return key
}

func NewGLimiter(redisOper RedisLimiterOperate, opts ...Option) GLimiter {
	return GLimiter{
		redisOper: redisOper,
		opts:      getOpts(opts...),
	}
}

// GetCounter GCounter, 获取key的计数限流器GCounter
//
// limitTime时间内计数达到limitThreshold阈值则触发限流，limitTime时间后会重置
func (g *GLimiter) GetGCounter(key string, limitThreshold int64, limitTime time.Duration) *GCounter {
	return &GCounter{
		glimiter:       g,
		limitTime:      limitTime,
		limitThreshold: limitThreshold,
		key:            g.getKey("GCounter", key),
		logger:         g.opts.logger.WithPrefix("GCounter").WithFields(log.Fields{"key": key}),
	}
}
