package test

import (
	"github.com/paceew/go-redisson/bloomfilter"
	"github.com/paceew/go-redisson/pkg/log"
	"github.com/paceew/go-redisson/redis"
	"testing"
)

var gbl bloomfilter.GBloomFilter

func TestBloomfilter(t *testing.T) {
	redisentry := redis.NewRedisEntry("127.0.0.1:6379", "")
	log.SetDefaultLogConfig("./", "bloomfilter.log", "", 300, 1000, log.TRACE, false)
	logger := log.NewFieLogger("", nil, log.TRACE)
	gbl = bloomfilter.NewGBloomFilter(redisentry, bloomfilter.WithLogger(logger))
	testFilter(logger)
}

func testFilter(logger log.FieldsLogger) {
	filter := gbl.GetGFilter("blfilter")
	if err := filter.TryInit(1000, 0.0001); err != nil {
		logger.Errorf("filter try init err:%s", err)
	}

	for i := 1; i <= 100; i++ {
		if err := filter.Add(i); err != nil {
			logger.Errorf("filter %d add err:%s", i, err)
		}
	}

	n := 66
	if ok, err := filter.Contains(n); err != nil {
		logger.Errorf("filter contains err:%s", err)
	} else {
		logger.Infof("filter %d contains %s", n, ok)
	}

	n = 107
	if ok, err := filter.Contains(n); err != nil {
		logger.Errorf("filter contains err:%s", err)
	} else {
		logger.Infof("filter %d contains %s", n, ok)
	}

	filter2 := gbl.GetGFilter("blfilter")
	if err := filter2.TryInit(1000, 0.0001); err != nil {
		logger.Errorf("filter2 try init err:%s", err)
	}
	n = 66
	if ok, err := filter2.Contains(n); err != nil {
		logger.Errorf("filter2 contains err:%s", err)
	} else {
		logger.Infof("filter2 %d contains %s", n, ok)
	}

	n = 107
	if ok, err := filter2.Contains(n); err != nil {
		logger.Errorf("filter2 contains err:%s", err)
	} else {
		logger.Infof("filter2 %d contains %s", n, ok)
	}

	if err := filter.Add("one"); err != nil {
		logger.Errorf("filter one add err:%s", err)
	}

	if err := filter.Add("two"); err != nil {
		logger.Errorf("filter two add err:%s", err)
	}

	if err := filter.Add("three"); err != nil {
		logger.Errorf("filter three add err:%s", err)
	}

	if ok, err := filter2.Contains("two"); err != nil {
		logger.Errorf("filter2 contains two err:%s", err)
	} else {
		logger.Infof("filter2 two contains %s", ok)
	}

	if ok, err := filter2.Contains("tw"); err != nil {
		logger.Errorf("filter2 contains tw err:%s", err)
	} else {
		logger.Infof("filter2 tw contains %s", ok)
	}
}
