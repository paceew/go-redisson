package main

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/paceew/go-redisson/bloomfilter"
	"github.com/paceew/go-redisson/glock"
	"github.com/paceew/go-redisson/limiter"
	"github.com/paceew/go-redisson/pkg/log"
	predis "github.com/paceew/go-redisson/redis"
)

func main() {
	pool := &redis.Pool{
		MaxIdle:     50,
		MaxActive:   0,
		IdleTimeout: time.Minute,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				"127.0.0.1:6379",
				redis.DialConnectTimeout(10*time.Second),
				redis.DialReadTimeout(5*time.Second),
				redis.DialWriteTimeout(5*time.Second),
				redis.DialDatabase(0),
				redis.DialPassword(""),
			)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	// 构建redisentry
	// redisentry是基于redigo，对于RedisGlockOperate接口的实现，可自定义RedisGlockOperate的实现
	redisentry := predis.NewRedisEntryWithPool(pool)
	// logger实现 log.FieldsLogger 日志接口
	log.SetDefaultLogConfig("./", "glock.log", "", 300, 1000, log.TRACE, false)
	logger := log.NewFieLogger("", log.Fields{}, log.TRACE)

	/********************************* glock ************************************/

	// glock实例，opts可选项：
	// glock.WithLogger(logger) 自定义日志，log.FieldsLogger日志接口
	// glock.WithPrefix("testll") redis key前缀
	// glock.WithUnPrefix() 不使用key前缀
	// glock.WithWatchDogTimeout(10*time.Second) 看门狗ttl
	// glock.WithTag("tag") 自定义tag，可重入锁会用到
	gl := glock.NewGlock(redisentry, glock.WithLogger(logger), glock.WithPrefix("testll"), glock.WithWatchDogTimeout(10*time.Second))

	/******************* Mutex ********************/

	// 获取Mutex
	mu := gl.GetMutex("testMu")
	// 最长等待30ms去尝试获取锁
	// 当不确定锁具体的持有时间时，使用这种方式获取锁会起一个协程开启看门狗策略
	// 当确定锁具体的持有时间时，可以使用mu.TryLock(30*time.Millisecond, 500*time.Second)代替
	if err := mu.TryLock(30 * time.Millisecond); err != nil {
		if errors.Is(err, glock.ErrNotObtain) {
			logger.Info("not botain")
		} else {
			logger.Errorf("lock err:%s", err.Error())
		}
		return
	}

	// do someting...

	if err := mu.UnLock(); err != nil {
		logger.Errorf("un lock err:%s", err.Error())
	}

	/******************* RWMutex ********************/

	// 获取RWMutex
	rwm := gl.GetRWMutex("testRWMu")
	// 读锁
	if err := rwm.TryRLock(30 * time.Millisecond); err != nil {
		if errors.Is(err, glock.ErrNotObtain) {
			logger.Info("not botain")
		} else {
			logger.Errorf("lock err:%s", err.Error())
		}
		return
	}

	// do someting...

	// 解读锁
	if err := rwm.RUnLock(); err != nil {
		logger.Errorf("un lock err:%s", err.Error())
	}

	// 写锁
	if err := rwm.TryLock(30 * time.Millisecond); err != nil {
		if errors.Is(err, glock.ErrNotObtain) {
			logger.Info("not botain")
		} else {
			logger.Errorf("lock err:%s", err.Error())
		}
		return
	}

	// do someting...

	// 解写锁
	if err := rwm.UnLock(); err != nil {
		logger.Errorf("un lock err:%s", err.Error())
	}

	/******************* ReLock ********************/

	// 获取ReLock
	// ReLock如果没指定tag，默认使用Glock的tag，
	// Glock如果没指定tag，Glock的默认tag为util.LocalIP() + ":" + util.PId()，意味着该程序内的可重入锁可重入
	// 如果要基于goroutine使用的可重入锁，那么可以使用rl := gl.GetReLockWithGoroutine("testRe")
	// 如果要自定义业务tag实现可重入锁，那么可以使用rl := gl.GetReLock("testRe","buskey")
	rl := gl.GetReLock("testRe")
	if err := rl.TryLock(30 * time.Millisecond); err != nil {
		if errors.Is(err, glock.ErrNotObtain) {
			logger.Info("not botain")
		} else {
			logger.Errorf("lock err:%s", err.Error())
		}
		return
	}

	// do someting...

	//
	if err := rl.UnLock(); err != nil {
		logger.Errorf("un lock err:%s", err.Error())
	}

	/********************************* bloomfilter ************************************/

	// bloomfilter实例，opts可选项：
	// bloomfilter.WithLogger(logger) 自定义日志，log.FieldsLogger日志接口
	// bloomfilter.WithPrefix("testll") redis key前缀
	// bloomfilter.WithUnPrefix() 不使用key前缀
	// bloomfilter.WithCheckConfig() 布隆过滤器每次操作前确认redis配置和实例配置是否一致，不一致返回ErrConfigUnequal错误
	// bloomfilter.WithHashfunc(hashfunc) 自定义hash函数，实现HashFunc接口
	gbl := bloomfilter.NewGBloomFilter(redisentry, bloomfilter.WithLogger(logger))

	filter := gbl.GetGFilter("blfilter")
	if err := filter.TryInit(1000, 0.0001); err != nil {
		logger.Errorf("filter try init err:%s", err)
	}

	if err := filter.Add(1); err != nil {
		logger.Errorf("filter 1 add err:%s", err)
	}
	if err := filter.Add(2); err != nil {
		logger.Errorf("filter 2 add err:%s", err)
	}
	if err := filter.Add("three"); err != nil {
		logger.Errorf("filter three add err:%s", err)
	}

	// 可设置过滤器过期时间
	if _, err := filter.Expired(24 * time.Hour); err != nil {
		logger.Errorf("filter PExpired err:%s", err)
	}

	exists, err := filter.Contains(2)
	if err != nil {
		logger.Errorf("filter contains err:%s", err)
	}

	if exists {
		logger.Info("filter 2 contains ")
	}

	// 可删除
	filter.Del()

	/********************************* limiter ************************************/
	glimiter := limiter.NewGLimiter(redisentry, limiter.WithLogger(logger))

	gcounter := glimiter.GetGCounter("counterlimite", 10, 20*time.Second)
	ok, err := gcounter.Limit(3)
	if err != nil {
		logger.Error(err)
	}
	if ok {
		logger.Info("allow 3")
	} else {
		logger.Info("not allow 3")
	}
}
