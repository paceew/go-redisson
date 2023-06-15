package log

import (
	"fmt"
	"strings"
)

type Level byte

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	OFF
)

type FieldsLogger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	WithPrefix(prefix string) FieldsLogger
	Prefix() string

	WithFields(fields Fields) FieldsLogger
	Fields() Fields

	SetLevel(level Level)
}

type Fields map[string]interface{}

func (f Fields) String() string {
	var builder strings.Builder
	for k, v := range f {
		builder.WriteString(fmt.Sprintf("%s=%+v;", k, v))
	}
	output := builder.String()
	length := len(output)
	if length > 0 {
		output = "[" + output[:length-1] + "] "
	}
	return output
}

func (f Fields) WithFields(newFields Fields) Fields {
	allFields := make(Fields, len(newFields)+len(f))

	for k, v := range f {
		allFields[k] = v
	}

	for k, v := range newFields {
		allFields[k] = v
	}

	return allFields
}
