package timerworker

import (
	"context"
	"errors"
	"github.com/paceew/go-redisson/pkg/log"
	"github.com/paceew/go-redisson/pkg/util"
	"github.com/tevino/abool"
	"time"
)

var (
	ErrNilWorker = errors.New("worker is nil")
	ErrNilLogger = errors.New("logger is nil")
	ErrHadRunned = errors.New("timer had runned")
	ErrNotRun    = errors.New("timer not run")
)

type Timer struct {
	ctx      context.Context
	cancel   context.CancelFunc
	Duration time.Duration
	ticker   *time.Ticker
	Logger   log.FieldsLogger
	run      abool.AtomicBool
	worker   Worker
}

func NewTimer(duration time.Duration, worker Worker, logger log.FieldsLogger) (*Timer, error) {
	if worker == nil {
		return nil, ErrNilWorker
	} else if logger == nil {
		return nil, ErrNilLogger
	}

	return &Timer{Duration: duration, Logger: logger, worker: worker}, nil
}

func (t *Timer) Stop() error {
	if !t.run.SetToIf(true, false) {
		return ErrNotRun
	}
	if t.cancel != nil {
		t.cancel()
	}
	t.worker.Done(t.ctx, t.Logger)
	t.ticker.Stop()
	return nil
}

func (t *Timer) done() {
	if t.run.SetToIf(true, false) {
		t.worker.Done(t.ctx, t.Logger)
		t.ticker.Stop()
	}
}

func (t *Timer) do(skipfirst bool) {
	t.ticker = time.NewTicker(t.Duration)
	defer func() {
		t.Logger.Trace("timer done...")
	}()
	t.Logger.Trace("timer doing...")
	if skipfirst {
		if t.worker.Do(t.ctx, t.Logger) {
			t.done()
			return
		}
	}
	for {
		select {
		case <-t.ticker.C:
			if t.worker.Do(t.ctx, t.Logger) {
				t.Stop()
			}
		case <-t.ctx.Done():
			t.done()
			return
		}
	}
}

func (t *Timer) Run(ctx context.Context, isblock, skipfirst bool) error {
	if !t.run.SetToIf(false, true) {
		return ErrHadRunned
	}
	t.ctx, t.cancel = context.WithCancel(ctx)
	if isblock {
		t.do(skipfirst)
	} else {
		go func() {
			util.Recover(t.Logger)
			t.do(skipfirst)
		}()
	}
	return nil
}
