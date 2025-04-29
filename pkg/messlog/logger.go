package messlog

import (
	"fmt"
	"log"
)

type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "ERROR"
	case LogLevelWarn:
		return "WARN"
	case LogLevelInfo:
		return "INFO"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	level LogLevel
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{level: level}
}

func (l *Logger) log(msgLevel LogLevel, format string, args ...any) {
	if msgLevel <= l.level {
		log.Printf("[%s] %s", msgLevel.String(), fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Error(format string, args ...any) { l.log(LogLevelError, format, args...) }
func (l *Logger) Warn(format string, args ...any)  { l.log(LogLevelWarn, format, args...) }
func (l *Logger) Info(format string, args ...any)  { l.log(LogLevelInfo, format, args...) }
func (l *Logger) Debug(format string, args ...any) { l.log(LogLevelDebug, format, args...) }
func (l *Logger) Trace(format string, args ...any) { l.log(LogLevelTrace, format, args...) }
