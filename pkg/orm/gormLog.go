package orm

import (
	"context"
	"errors"
	"github.com/aiwuTech/fileLogger"
	"github.com/paceew/go-redisson/pkg/log"
	glog "gorm.io/gorm/logger"
	"strconv"
	"time"
)

var (
	gormLogSingletons *fileLogger.FileLogger
)

//	type Interface interface {
//	    LogMode(LogLevel) Interface
//	    Info(context.Context, string, ...interface{})
//	    Warn(context.Context, string, ...interface{})
//	    Error(context.Context, string, ...interface{})
//	    Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
//	}
type GormLog struct {
	Logger     *fileLogger.FileLogger
	Config     GormLogConfig
	prefix     string
	fields     log.Fields
	preparaStr *string
}

type GormLogConfig struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	Level                     glog.LogLevel
	LogName                   string
	LogPath                   string
	Prefix                    string
}

func NewGormLogger(config GormLogConfig) *GormLog {
	if gormLogSingletons == nil {
		flevel := fileLogger.LEVEL(5 - int(config.Level))
		config := log.LogConfig{
			FileDir:  config.LogPath,
			FileName: config.LogName,
			LogScan:  log.DefaultLogConfig.LogScan,
			LogSeq:   log.DefaultLogConfig.LogSeq,
			Console:  false,
			Level:    flevel,
		}
		gormLogSingletons = log.NewLogInstance(config)
	}
	return &GormLog{Logger: gormLogSingletons, Config: config, prefix: config.Prefix}
}

func (gl *GormLog) logPreparation() string {
	if gl.preparaStr == nil {
		preparaStr := ""
		if gl.prefix != "" {
			preparaStr += "[" + gl.prefix + "]"
		}

		if gl.fields != nil {
			preparaStr += gl.fields.String() + " "
		}

		gl.preparaStr = &preparaStr
	}

	return *gl.preparaStr
}

func (gl *GormLog) LogMode(lv glog.LogLevel) glog.Interface {
	gl.Logger.SetLogLevel(fileLogger.LEVEL(5 - int(lv)))
	return gl
}

func (gl *GormLog) Info(ctx context.Context, format string, args ...interface{}) {
	gl.Logger.Info(gl.logPreparation()+format, args...)
}

func (gl *GormLog) Warn(ctx context.Context, format string, args ...interface{}) {
	gl.Logger.Warn(gl.logPreparation()+format, args...)
}

func (gl *GormLog) Error(ctx context.Context, format string, args ...interface{}) {
	gl.Logger.Error(gl.logPreparation()+format, args...)
}

func (gl *GormLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	rowstr := "-"
	if rows != -1 {
		rowstr = strconv.FormatInt(rows, 10)
	}
	if err != nil && (!errors.Is(err, glog.ErrRecordNotFound) || !gl.Config.IgnoreRecordNotFoundError) {
		gl.Error(ctx, "[%v][rows:%s][err:%s] %s", elapsed, rowstr, err.Error(), sql)
	} else if elapsed > gl.Config.SlowThreshold && gl.Config.SlowThreshold != 0 {
		gl.Warn(ctx, "[%v][rows:%s][%s:%v] %s", elapsed, rowstr, "SLOW SQL", gl.Config.SlowThreshold, sql)
	} else {
		gl.Info(ctx, "[%v][rows:%s] %s", elapsed, rowstr, sql)
	}
}

func (gl *GormLog) WithPrefix(prefix string) glog.Interface {
	config := gl.Config
	config.Prefix = prefix
	return NewGormLogger(config)
}

func (gl *GormLog) Prefix() string {
	return gl.prefix
}
