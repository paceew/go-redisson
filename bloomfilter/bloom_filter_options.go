package bloomfilter

import (
	"github.com/paceew/go-redisson/pkg/log"
)

var defalultOpts *GBloomFilterOptions

func init() {
	defalultOpts = &GBloomFilterOptions{
		logger:   log.NewEmptyLogger(),
		prefix:   "GBLOOMFILTER",
		hashfunc: HashMur3,
	}
}

type GBloomFilterOptions struct {
	logger      log.FieldsLogger
	prefix      string
	unPrefix    bool
	checkConfig bool
	hashfunc    HashFunc
}

type Option func(opts *GBloomFilterOptions)

func getOpts(opts ...Option) *GBloomFilterOptions {
	options := &GBloomFilterOptions{
		logger:   defalultOpts.logger,
		prefix:   defalultOpts.prefix,
		hashfunc: defalultOpts.hashfunc,
	}
	for _, v := range opts {
		v(options)
	}
	return options
}

func WithLogger(logger log.FieldsLogger) Option {
	return func(opts *GBloomFilterOptions) {
		opts.logger = logger
	}
}

func WithPrefix(prefix string) Option {
	return func(opts *GBloomFilterOptions) {
		opts.prefix = prefix
		opts.unPrefix = false
	}
}

func WithUnPrefix() Option {
	return func(opts *GBloomFilterOptions) {
		opts.unPrefix = true
	}
}

func WithCheckConfig() Option {
	return func(opts *GBloomFilterOptions) {
		opts.checkConfig = true
	}
}

func WithHashfunc(fun HashFunc) Option {
	return func(opts *GBloomFilterOptions) {
		opts.hashfunc = fun
	}
}
