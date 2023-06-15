package glock

import (
	"github.com/paceew/go-redisson/pkg/util"
	"time"
)

type GMutex struct {
	lock *distLock
}

// TryLock 加锁，尝试Mutex加锁
//
// @maxWaiTime 最长等待时间，如果为0则只会获取一次，小于0会返回ErrNotObtain
// @holdingTime 持锁时间，获取到锁后生效，为本次获取锁的持有时间，如果没选，则是无限时长，开启看门狗策略
func (rw *GMutex) TryLock(maxWaiTime time.Duration, holdingTimes ...time.Duration) error {
	var holdingTime time.Duration = 0
	if len(holdingTimes) > 0 {
		holdingTime = holdingTimes[0]
	}
	// Mutex每次用不同的tag
	tag := util.GenUUID()
	return rw.lock.tryLock(maxWaiTime, holdingTime, tag)
}

// UnLock 解锁，解除单个Mutex
//
// 如果该实例获取Mutex不成功则不应该调用该方法，虽然调用该方法不会解锁其他相同key已获取的Mutex
// ，因为Mutex的tag都是随机生成不一样的
func (rw *GMutex) UnLock() error {
	return rw.lock.unlock()
}
