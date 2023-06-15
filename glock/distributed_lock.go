package glock

import (
	"context"
	"errors"
	"github.com/paceew/go-redisson/pkg/log"
	"sync/atomic"
	"time"
)

type distLock struct {
	key, publishKey string
	tag             string
	glock           *Glock
	watchDog        *WatchDog
	logger          log.FieldsLogger
	// 加锁计数器
	counts uint64
}

// TryLock 尝试获取lock
//
// @maxWaiTime 最长等待时间，如果为0则只会获取一次，小于0会返回ErrNotObtain
// @holdingTime 持锁时间，获取到锁后生效，为本次获取锁的持有时间，如果为0，则是无限时长，开启看门狗策略
func (dm *distLock) tryLock(maxWaiTime, holdingTime time.Duration, tag string) error {
	if maxWaiTime < 0 {
		return ErrNotObtain
	}
	now := time.Now()
	// 尝试获取锁
	result, err := dm.tryObtain(holdingTime, tag)
	if err != nil {
		return err
	} else if result == 0 {
		// 获取锁成功
		return nil
	} else if maxWaiTime == 0 {
		// maxWaiTime = 0 的情况下只获取一次
		return ErrNotObtain
	}
	// else if result > maxWaiTime.Milliseconds() {
	// 	//剩余锁时间大于最大等待时间，直接返回失败
	// 	//TODO watchdog 机制会导致过期时间不一定准确，或者锁有可能提前释放
	// 	return ErrNotObtain
	// }

	/********* 订阅解锁，超时返回，没超时则继续尝试获取锁 *********/
	left := time.Duration(result) * time.Millisecond
	timeout := maxWaiTime - time.Since(now)
	if left < timeout {
		timeout = left
	}
	dm.logger.Tracef("subscribe wait %s %s", dm.publishKey, timeout)
	_, err = dm.glock.redisOper.SubscribeAndReceiveMessage(dm.publishKey, timeout)
	if err != nil {
		if errors.Is(err, ErrSubTimeout) {
			return ErrNotObtain
		}
		return err
	}
	dm.logger.Trace("subscribe receive message and try lock again")
	// 收到解锁，继续尝试获取锁
	return dm.tryLock(maxWaiTime-time.Since(now), holdingTime, tag)
}

// tryObtain 尝试获取锁
func (dm *distLock) tryObtain(ttl time.Duration, tag string) (int64, error) {
	watchDog := false
	if ttl == 0 {
		ttl = dm.glock.opts.watchDogTimeout
		watchDog = true
	}
	mttl, err := dm.glock.redisOper.ObtainLock(dm.key, tag, ttl)
	dm.logger.Tracef("try obtain lock,tag %s ,mttl:%d err:%v", tag, mttl, err)
	if err != nil {
		return -1, err
	} else if mttl == 0 {
		// 获取到锁就把 tag 存进去
		dm.tag = tag
		atomic.AddUint64(&dm.counts, 1)
		dm.logger.Tracef("obtain lock,tag %s ,counts %d", tag, dm.counts)
		if watchDog {
			// watch dog
			dm.watchByWatchdog()
		}
	}

	return mttl, nil
}

func (dm *distLock) watchByWatchdog() error {
	if dm.watchDog == nil {
		dm.logger.Trace("create watch dog")
		dm.watchDog = NewWatchDog(dm.glock.opts.watchDogTimeout/3, dm.RenewExpiration, dm.logger.WithFields(log.Fields{"WatchDog": dm.glock.opts.watchDogTimeout}))
	}
	dm.logger.Trace("open watch dog")
	return dm.watchDog.Watch(context.TODO())
}

func (dm *distLock) releaseWatchdog() error {
	if dm.watchDog != nil {
		dm.logger.Trace("release watch dog")
		return dm.watchDog.Release()
	}
	return nil
}

func (dm *distLock) RenewExpiration(ctx context.Context, interval time.Duration, logger log.FieldsLogger) (isend bool) {
	result, err := dm.glock.redisOper.RenewLock(dm.key, dm.tag, interval*3)
	if err != nil {
		logger.Errorf("RenewExpiration error:%s", err.Error())
		return false
	} else if result == 0 {
		logger.Error("RenewExpiration fail, the lock not exists")
		return true
	} else if result == -1 {
		logger.Error("RenewExpiration fail, the lock tag wrong")
		return true
	}
	dm.logger.Tracef("renew expiration %s %s", dm.tag, interval*3)
	return false
}

// Unlock
func (dm *distLock) unlock(tag ...string) error {
	tag_val := dm.tag
	if len(tag) > 0 {
		tag_val = tag[0]
	}
	result, err := dm.glock.redisOper.ReleaseLock(dm.key, tag_val, dm.publishKey, MSG_UNLOCK)
	dm.logger.Tracef("release lock %s result:%d err:%v", tag_val, result, err)
	if err != nil {
		return err
	}
	if result == 0 {
		// 锁不存在
		return ErrNotExists
	} else if result == -1 {
		// 锁不属于该tag
		return ErrWrongTag
	}
	if atomic.AddUint64(&dm.counts, ^uint64(0)) == 0 {
		return dm.releaseWatchdog()
	}
	return nil
}
