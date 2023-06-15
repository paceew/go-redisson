package limiter

import (
	"github.com/paceew/go-redisson/pkg/log"
)

var defalultOpts *GLimiterOptions

func init() {
	defalultOpts = &GLimiterOptions{
		logger: log.NewEmptyLogger(),
		prefix: "GLIMITER",
	}
}

type GLimiterOptions struct {
	logger   log.FieldsLogger
	prefix   string
	unPrefix bool
}

type Option func(opts *GLimiterOptions)

func getOpts(opts ...Option) *GLimiterOptions {
	options := &GLimiterOptions{
		logger: defalultOpts.logger,
		prefix: defalultOpts.prefix,
	}
	for _, v := range opts {
		v(options)
	}
	return options
}

func WithLogger(logger log.FieldsLogger) Option {
	return func(opts *GLimiterOptions) {
		opts.logger = logger
	}
}

func WithPrefix(prefix string) Option {
	return func(opts *GLimiterOptions) {
		opts.prefix = prefix
		opts.unPrefix = false
	}
}

func WithUnPrefix() Option {
	return func(opts *GLimiterOptions) {
		opts.unPrefix = true
	}
}
