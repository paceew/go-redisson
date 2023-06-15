package glock

import (
	"time"
)

type GReLock struct {
	lock *distLock
	tag  string
}

// TryLock 加锁，尝试ReLock加锁一次
//
// @maxWaiTime 最长等待时间，如果为0则只会获取一次，小于0会返回ErrNotObtain
// @holdingTime 持锁时间，获取到锁后生效，为本次获取锁的持有时间，如果没选，则是无限时长，开启看门狗策略
func (rl *GReLock) TryLock(maxWaiTime time.Duration, holdingTimes ...time.Duration) error {
	var holdingTime time.Duration = 0
	if len(holdingTimes) > 0 {
		holdingTime = holdingTimes[0]
	}
	return rl.lock.tryLock(maxWaiTime, holdingTime, rl.tag)
}

// UnLock 解锁,解除ReLock一次
//
// 如果该实例获取ReLock不成功则不应该调用该方法，否则会有解锁其他相同key的ReLock的风险
// ，因为ReLock的tag都是一样的
func (rl *GReLock) UnLock() error {
	return rl.lock.unlock(rl.tag)
}
