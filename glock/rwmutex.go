package glock

import (
	"github.com/paceew/go-redisson/pkg/util"
	"time"
)

type GRWMutex struct {
	lock    *distLock
	readTag string
}

const (
	READ_TAG string = "GLOCK_GRWMUTEX_READ_TAG"
)

// TryRLock 读加锁，尝试RLock加锁一次
//
// @maxWaiTime 最长等待时间，如果为0则只会获取一次，小于0会返回ErrNotObtain
// @holdingTime 持锁时间，获取到锁后生效，为本次获取锁的持有时间，如果没选，则是无限时长，开启看门狗策略
func (rw *GRWMutex) TryRLock(maxWaiTime time.Duration, holdingTimes ...time.Duration) error {
	var holdingTime time.Duration = 0
	if len(holdingTimes) > 0 {
		holdingTime = holdingTimes[0]
	}
	// 读锁用相同的tag
	return rw.lock.tryLock(maxWaiTime, holdingTime, rw.readTag)
}

// RUnLock 解锁,解除RLock一次
//
// 如果该实例获取RLock不成功则不应该调用该方法，否则会有解锁其他相同key的RLock的风险
// ，因为RLock的tag都是一样的
func (rw *GRWMutex) RUnLock() error {
	return rw.lock.unlock(rw.readTag)
}

// TryLock 写加锁，尝试WLock加锁
//
// @maxWaiTime 最长等待时间，如果为0则只会获取一次，小于0会返回ErrNotObtain
// @holdingTime 持锁时间，获取到锁后生效，为本次获取锁的持有时间，如果没选，则是无限时长，开启看门狗策略
func (rw *GRWMutex) TryLock(maxWaiTime time.Duration, holdingTimes ...time.Duration) error {
	var holdingTime time.Duration = 0
	if len(holdingTimes) > 0 {
		holdingTime = holdingTimes[0]
	}
	// 写锁每次用不同的tag
	tag := util.GenUUID()
	return rw.lock.tryLock(maxWaiTime, holdingTime, tag)
}

// UnLock 解锁，解除单个WLock
//
// 如果该实例获取WLock不成功则不应该调用该方法，虽然调用该方法不会解锁其他相同key已获取的WLock
// ，因为WLock的tag都是随机生成不一样的
func (rw *GRWMutex) UnLock() error {
	return rw.lock.unlock()
}
