package timerworker

import (
	"context"
	"github.com/paceew/go-redisson/pkg/log"
)

type Worker interface {
	Do(ctx context.Context, logger log.FieldsLogger) (end bool)
	Done(ctx context.Context, logger log.FieldsLogger)
}
