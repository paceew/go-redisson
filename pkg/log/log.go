package log

import (
	"github.com/aiwuTech/fileLogger"
	"github.com/paceew/go-redisson/pkg/config"
)

var (
	log              *fileLogger.FileLogger
	DefaultLogConfig = &LogConfig{"./", "sys.log", "", 300, 1000, false, fileLogger.TRACE}
)

type LogConfig struct {
	FileDir  string
	FileName string
	Prefix   string
	LogScan  int64
	LogSeq   int
	Console  bool
	Level    fileLogger.LEVEL
}

// logScan:日志检查时间 ，logSeq:日志通道buf数
func SetDefaultLogConfig(fileDir, fileName, prefix string, logScan int64, logSeq int, level Level, console bool) {
	DefaultLogConfig.FileDir = fileDir
	DefaultLogConfig.FileName = fileName
	DefaultLogConfig.Prefix = prefix
	DefaultLogConfig.LogScan = logScan
	DefaultLogConfig.LogSeq = logSeq
	DefaultLogConfig.Console = console
	DefaultLogConfig.Level = fileLogger.LEVEL(level)
}

func initDefaultLog(config LogConfig) {
	log = fileLogger.NewDailyLogger(config.FileDir, config.FileName, config.Prefix, config.LogScan, config.LogSeq)
	log.SetLogLevel(config.Level)
	log.SetLogConsole(config.Console)
}

// GetDefaultLogSingletons 获取default log单例
func GetDefaultLogSingletons() *fileLogger.FileLogger {
	if log == nil {
		initDefaultLog(*DefaultLogConfig)
	}

	return log
}

func NewLogInstance(config LogConfig) *fileLogger.FileLogger {
	LogInstance := fileLogger.NewDailyLogger(config.FileDir, config.FileName, config.Prefix, config.LogScan, config.LogSeq)
	LogInstance.SetLogLevel(config.Level)
	LogInstance.SetLogConsole(config.Console)
	return LogInstance
}

func AutoInitLoggerWithConfig() {
	logPath := config.VipCfg().GetStringWithDefault("logger.log_path", "./")
	logFileName := config.VipCfg().GetStringWithDefault("logger.log_name", "sys.log")
	loglevel := Level(config.VipCfg().GetIntWithDefault("logger.level", 1))
	if loglevel > OFF {
		loglevel = INFO
	}
	SetDefaultLogConfig(logPath, logFileName, "", 300, 1000, loglevel, false)
}
