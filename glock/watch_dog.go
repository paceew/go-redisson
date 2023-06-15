package glock

import (
	"context"
	"github.com/paceew/go-redisson/pkg/log"
	"github.com/paceew/go-redisson/pkg/timerworker"
	"time"
)

type WatchDog struct {
	watchDogTimeout time.Duration
	watchFor        WatchForFunc
	timer           *timerworker.Timer
}

type WatchForFunc func(ctx context.Context, interval time.Duration, logger log.FieldsLogger) (isend bool)

func NewWatchDog(watchDogTimeout time.Duration, watchFor WatchForFunc, logger log.FieldsLogger) *WatchDog {
	wd := &WatchDog{
		watchDogTimeout: watchDogTimeout,
		watchFor:        watchFor,
	}
	timer, _ := timerworker.NewTimer(watchDogTimeout, wd, logger)
	wd.timer = timer
	return wd
}

func (wd *WatchDog) Do(ctx context.Context, logger log.FieldsLogger) (isend bool) {
	if wd.watchFor == nil {
		return true
	}
	return wd.watchFor(ctx, wd.watchDogTimeout, logger)
}

func (wd *WatchDog) Done(ctx context.Context, logger log.FieldsLogger) {
}

func (wd *WatchDog) Watch(ctx context.Context) error {
	return wd.timer.Run(ctx, false, true)
}

func (wd *WatchDog) Release() error {
	return wd.timer.Stop()
}
