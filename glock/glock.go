package glock

import (
	"errors"

	"github.com/paceew/go-redisson/pkg/log"
	"github.com/paceew/go-redisson/pkg/util"
)

var (
	ErrNotObtain       = errors.New("not obtain the lock")
	ErrNotObtainedLock = errors.New("had not obtain the lock")
	ErrWrongTag        = errors.New("the lock does not belong to you")
	ErrNotExists       = errors.New("the lock does not exists")
)

type Glock struct {
	redisOper RedisGlockOperate
	opts      *GlockOptions
}

func (g *Glock) getSubscribeKey(locktype, key string) string {
	if !g.opts.unPrefix {
		return g.opts.prefix + "_PUBSUB_" + locktype + "_" + key
	}
	return "PUBSUB_" + key
}

func (g *Glock) getKey(locktype, key string) string {
	if !g.opts.unPrefix {
		return g.opts.prefix + "_LOCK_" + locktype + "_" + key
	}
	return key
}

// NewGlock Glock实例
func NewGlock(redisOper RedisGlockOperate, opts ...Option) Glock {
	return Glock{
		redisOper: redisOper,
		opts:      getOpts(opts...),
	}
}

// GetMutex GMutex, 获取key的不可重入锁GMutex
func (g *Glock) GetMutex(key string) *GMutex {
	dw := &distLock{
		glock:      g,
		key:        g.getKey("GMutex", key),
		publishKey: g.getSubscribeKey("GMutex", key),
		logger:     g.opts.logger.WithPrefix("GMutex").WithFields(log.Fields{"key": key}),
	}
	return &GMutex{dw}
}

// GetRWMutex GRWMutex, 获取key的读写锁GRWMutex
func (g *Glock) GetRWMutex(key string) *GRWMutex {
	dw := &distLock{
		glock:      g,
		key:        g.getKey("GRWMutex", key),
		publishKey: g.getSubscribeKey("GRWMutex", key),
		logger:     g.opts.logger.WithPrefix("GRWMutex").WithFields(log.Fields{"key": key}),
	}
	return &GRWMutex{lock: dw, readTag: READ_TAG}
}

// GetReLock GReLock, 获取key tag的可重入锁GReLock
//
// @tag 相同的key tag的锁互为可重入，如tag为空，则使用glock的tag（默认为 LocalIP() + ":" + PId()）
//
// 如果要基于goroutine使用的可重入锁，那么可以使用 GetReLock(key,LocalIP() + ":" + PId() + ":" + GoId())
func (g *Glock) GetReLock(key string, tag ...string) *GReLock {
	// ReLock 通过获取相同的tag实现可重入
	tag_val := g.opts.tag
	if len(tag) > 0 {
		tag_val = tag[0]
	}
	dw := &distLock{
		glock:      g,
		key:        g.getKey("GReLock", key),
		publishKey: g.getSubscribeKey("GReLock", key),
		logger:     g.opts.logger.WithPrefix("GReLock").WithFields(log.Fields{"key": key, "retag": tag_val}),
	}
	return &GReLock{lock: dw, tag: tag_val}
}

// GetReLockWithGoroutine, 获取key groutineid为tag的可重入锁GReLock
//
// 同一个groutine可重入,等价于GetReLock(key,LocalIP() + ":" + PId() + ":" + GoId())
func (g *Glock) GetReLockWithGoroutine(key string) *GReLock {
	// 使用goroutine id 进行可重入
	return g.GetReLock(key, util.LocalIP()+":"+util.PId()+":"+util.GoId())
}
