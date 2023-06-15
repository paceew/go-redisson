package test

import (
	"fmt"
	"github.com/paceew/go-redisson/limiter"
	"github.com/paceew/go-redisson/pkg/log"
	"github.com/paceew/go-redisson/redis"
	"testing"
	"time"
)

var glimiter limiter.GLimiter

func TestLimiter(t *testing.T) {
	redisentry := redis.NewRedisEntry("127.0.0.1:6379", "")
	log.SetDefaultLogConfig("./", "glimiter.log", "", 300, 1000, log.TRACE, false)
	logger := log.NewFieLogger("", nil, log.TRACE)
	glimiter = limiter.NewGLimiter(redisentry, limiter.WithLogger(logger))
	testGCounter(logger)
}

func testGCounter(logger log.FieldsLogger) {
	gcounter := glimiter.GetGCounter("counterlimite", 10, 20*time.Second)
	ok, err := gcounter.Limit(3)
	if err != nil {
		fmt.Println(err)
		logger.Error(err)
	}
	if ok {
		fmt.Println("allow 3")
		logger.Info("allow 3")
	} else {
		fmt.Println("not allow 3")
		logger.Info("not allow 3")
	}
	time.Sleep(5 * time.Second)

	ok, err = gcounter.Limit(4)
	if err != nil {
		fmt.Println(err)
		logger.Error(err)
	}
	if ok {
		fmt.Println("allow 4")
		logger.Info("allow 4")
	} else {
		fmt.Println("not allow 4")
		logger.Info("not allow 4")
	}

	time.Sleep(5 * time.Second)
	ok, err = gcounter.Limit(5)
	if err != nil {
		fmt.Println(err)
		logger.Error(err)
	}
	if ok {
		fmt.Println("allow 5")
		logger.Info("allow 5")
	} else {
		fmt.Println("not allow 5")
		logger.Info("not allow 5")
	}

	time.Sleep(10 * time.Second)
	ok, err = gcounter.Limit(6)
	if err != nil {
		fmt.Println(err)
		logger.Error(err)
	}
	if ok {
		fmt.Println("allow 6")
		logger.Info("allow 6")
	} else {
		fmt.Println("not allow 6")
		logger.Info("not allow 6")
	}
}
