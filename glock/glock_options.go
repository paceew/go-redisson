package glock

import (
	"github.com/paceew/go-redisson/pkg/log"
	"github.com/paceew/go-redisson/pkg/util"
	"time"
)

var defalultOpts *GlockOptions

func init() {
	defalultOpts = &GlockOptions{
		watchDogTimeout: 30 * time.Second,
		tag:             util.LocalIP() + ":" + util.PId(),
		logger:          log.NewEmptyLogger(),
		prefix:          "GLOCK",
	}
}

type GlockOptions struct {
	watchDogTimeout time.Duration
	//dist Lock pool
	//distLockPool
	tag string
	// rwLock    sync.RWLock
	// distLocks map[string]*distLock
	logger   log.FieldsLogger
	prefix   string
	unPrefix bool
}

type Option func(opts *GlockOptions)

func getOpts(opts ...Option) *GlockOptions {
	options := &GlockOptions{
		watchDogTimeout: defalultOpts.watchDogTimeout,
		tag:             defalultOpts.tag,
		logger:          defalultOpts.logger,
		prefix:          defalultOpts.prefix,
	}
	for _, v := range opts {
		v(options)
	}
	return options
}

func WithWatchDogTimeout(timeout time.Duration) Option {
	return func(opts *GlockOptions) {
		opts.watchDogTimeout = timeout
	}
}

func WithTag(tag string) Option {
	return func(opts *GlockOptions) {
		opts.tag = tag
	}
}

func WithLogger(logger log.FieldsLogger) Option {
	return func(opts *GlockOptions) {
		opts.logger = logger
	}
}

func WithPrefix(prefix string) Option {
	return func(opts *GlockOptions) {
		opts.prefix = prefix
		opts.unPrefix = false
	}
}

func WithUnPrefix() Option {
	return func(opts *GlockOptions) {
		opts.unPrefix = true
	}
}
