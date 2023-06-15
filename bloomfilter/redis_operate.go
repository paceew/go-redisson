package bloomfilter

import (
	"time"
)

type RedisBloomFilterOperate interface {
	HMGet(keys ...interface{}) (map[string]string, error)
	HMSet(kvs ...interface{}) error
	SetBit(key string, offset uint64, value int) error
	GetBit(key string, offset uint64) (int64, error)
	Expire(key string, second time.Duration) (int64, error)
	Del(key string) error
}
