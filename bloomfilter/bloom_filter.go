package bloomfilter

import (
	"github.com/paceew/go-redisson/pkg/log"
)

const RedisMaxLength = 8 * 512 * 1024 * 1024

type GBloomFilter struct {
	redisOper RedisBloomFilterOperate
	opts      *GBloomFilterOptions
}

func (g *GBloomFilter) getKey(key string) string {
	if !g.opts.unPrefix {
		return g.opts.prefix + "_" + key
	}
	return key
}

func (g *GBloomFilter) getConfigKey(key string) string {
	key = key + "Config"
	if !g.opts.unPrefix {
		return g.opts.prefix + "_" + key
	}
	return key
}

func NewGBloomFilter(redisOper RedisBloomFilterOperate, opts ...Option) GBloomFilter {
	return GBloomFilter{
		redisOper: redisOper,
		opts:      getOpts(opts...),
	}
}

// GetGFilter GFilter, 获取key的布隆过滤器实例GFilter
func (g *GBloomFilter) GetGFilter(key string) *GFilter {
	return &GFilter{
		gBloomFilter: g,
		key:          g.getKey(key),
		configKey:    g.getConfigKey(key),
		hashfunc:     g.opts.hashfunc,
		checkConfig:  g.opts.checkConfig,
		logger:       g.opts.logger.WithPrefix("GBloomFilter").WithFields(log.Fields{"key": key}),
	}
}
