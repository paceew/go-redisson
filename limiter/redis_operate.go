package limiter

import (
	"time"
)

type RedisLimiterOperate interface {
	IncrBy(key string, increment int64) (int64, error)
	PExpire(key string, millisecond time.Duration) (int64, error)
	PTTL(key string) (int64, error)
}
