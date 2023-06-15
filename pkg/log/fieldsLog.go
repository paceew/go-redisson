package log

import (
	"fmt"

	"github.com/aiwuTech/fileLogger"
)

type FieLog struct {
	Logger     *fileLogger.FileLogger
	prefix     string
	fields     Fields
	preparaStr *string
	level      Level
}

// NewFieLogger 是以default log单例作为instance,所以不仅会受制于本身的log level，也会受制于default log 的log level
func NewFieLogger(prefix string, fields Fields, level Level) *FieLog {
	return &FieLog{
		//以default log单例作为instance
		Logger: GetDefaultLogSingletons(),
		prefix: prefix,
		fields: fields,
		level:  level,
	}
}

func (f *FieLog) logPreparation() string {
	if f.preparaStr == nil {
		preparaStr := ""
		if f.prefix != "" {
			preparaStr += "[" + f.prefix + "]"
		}

		if f.fields != nil {
			preparaStr += f.fields.String()
		}

		preparaStr += " "
		f.preparaStr = &preparaStr
	}

	return *f.preparaStr
}

func (f *FieLog) Print(v ...interface{}) {
	f.Logger.Info(f.logPreparation() + fmt.Sprint(v...))
}

func (f *FieLog) Printf(format string, v ...interface{}) {
	f.Logger.Info(f.logPreparation()+format, v...)
}

func (f *FieLog) Trace(v ...interface{}) {
	if f.level <= TRACE {
		f.Logger.Trace(f.logPreparation() + fmt.Sprint(v...))
	}
}

func (f *FieLog) Tracef(format string, v ...interface{}) {
	if f.level <= TRACE {
		f.Logger.Trace(f.logPreparation()+format, v...)
	}
}

func (f *FieLog) Debug(v ...interface{}) {
	if f.level <= DEBUG {
		f.Logger.Trace(f.logPreparation() + fmt.Sprint(v...))
	}
}

func (f *FieLog) Debugf(format string, v ...interface{}) {
	if f.level <= DEBUG {
		f.Logger.Trace(f.logPreparation()+format, v...)
	}
}

func (f *FieLog) Info(v ...interface{}) {
	if f.level <= INFO {
		f.Logger.Info(f.logPreparation() + fmt.Sprint(v...))
	}

}

func (f *FieLog) Infof(format string, v ...interface{}) {
	if f.level <= INFO {
		f.Logger.Info(f.logPreparation()+format, v...)
	}
}

func (f *FieLog) Warn(v ...interface{}) {
	if f.level <= WARN {
		f.Logger.Warn(f.logPreparation() + fmt.Sprint(v...))
	}
}

func (f *FieLog) Warnf(format string, v ...interface{}) {
	if f.level <= WARN {
		f.Logger.Warn(f.logPreparation()+format, v...)
	}
}

func (f *FieLog) Error(v ...interface{}) {
	if f.level <= ERROR {
		f.Logger.Error(f.logPreparation() + fmt.Sprint(v...))
	}
}

func (f *FieLog) Errorf(format string, v ...interface{}) {
	if f.level <= ERROR {
		f.Logger.Error(f.logPreparation()+format, v...)
	}
}

func (f *FieLog) WithPrefix(prefix string) FieldsLogger {
	return NewFieLogger(prefix, f.fields, f.level)
}

func (f *FieLog) Prefix() string {
	return f.prefix
}

func (f *FieLog) WithFields(fields Fields) FieldsLogger {
	return NewFieLogger(f.prefix, f.fields.WithFields(fields), f.level)
}

func (f *FieLog) Fields() Fields {
	return f.fields
}

func (f *FieLog) SetLevel(level Level) {
	f.level = level
}

/***************************  Empty Implement  *****************************/
type EmptyLog struct {
}

func NewEmptyLogger() *EmptyLog {
	return &EmptyLog{}
}

func (f *EmptyLog) Print(v ...interface{}) {
}

func (f *EmptyLog) Printf(format string, v ...interface{}) {
}

func (f *EmptyLog) Trace(v ...interface{}) {
}

func (f *EmptyLog) Tracef(format string, v ...interface{}) {
}

func (f *EmptyLog) Debug(v ...interface{}) {
}

func (f *EmptyLog) Debugf(format string, v ...interface{}) {
}

func (f *EmptyLog) Info(v ...interface{}) {
}

func (f *EmptyLog) Infof(format string, v ...interface{}) {
}

func (f *EmptyLog) Warn(v ...interface{}) {
}

func (f *EmptyLog) Warnf(format string, v ...interface{}) {
}

func (f *EmptyLog) Error(v ...interface{}) {
}

func (f *EmptyLog) Errorf(format string, v ...interface{}) {
}

func (f *EmptyLog) WithPrefix(prefix string) FieldsLogger {
	return NewEmptyLogger()
}

func (f *EmptyLog) Prefix() string {
	return ""
}

func (f *EmptyLog) WithFields(fields Fields) FieldsLogger {
	return NewEmptyLogger()
}

func (f *EmptyLog) Fields() Fields {
	return nil
}

func (f *EmptyLog) SetLevel(level Level) {
}
