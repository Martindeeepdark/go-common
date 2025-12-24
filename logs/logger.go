package logs

import "context"

// FormatLogger is a logs interface that output logs with a format.
type FormatLogger interface {
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Noticef(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// Logger is a logs interface that provides logging function with levels.
type Logger interface {
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Notice(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

// CtxLogger is a logs interface that accepts a context argument and output
// logs with a format.
type CtxLogger interface {
	CtxTracef(ctx context.Context, format string, v ...interface{})
	CtxDebugf(ctx context.Context, format string, v ...interface{})
	CtxInfof(ctx context.Context, format string, v ...interface{})
	CtxNoticef(ctx context.Context, format string, v ...interface{})
	CtxWarnf(ctx context.Context, format string, v ...interface{})
	CtxErrorf(ctx context.Context, format string, v ...interface{})
	CtxFatalf(ctx context.Context, format string, v ...interface{})
}

// Control provides methods to config a logs.
type Control interface {
	SetLevel(Level)
	SetOutput(output interface{}) // io.Writer or zap.WriteSyncer
}

// FullLogger is the combination of Logger, FormatLogger, CtxLogger and Control.
type FullLogger interface {
	Logger
	FormatLogger
	CtxLogger
	Control
}

// Level defines the priority of a log message.
// When a logs is configured with a level, any log message with a lower
// log level (smaller by integer comparison) will not be output.
type Level int

// The levels of logs.
const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
)

// String returns the string representation of the level.
func (lv Level) String() string {
	switch lv {
	case LevelTrace:
		return "[Trace]"
	case LevelDebug:
		return "[Debug]"
	case LevelInfo:
		return "[Info]"
	case LevelNotice:
		return "[Notice]"
	case LevelWarn:
		return "[Warn]"
	case LevelError:
		return "[Error]"
	case LevelFatal:
		return "[Fatal]"
	default:
		return "[?Unknown]"
	}
}

// ParseLevel parses a string to Level.
func ParseLevel(level string) Level {
	switch level {
	case "trace", "TRACE":
		return LevelTrace
	case "debug", "DEBUG":
		return LevelDebug
	case "info", "INFO":
		return LevelInfo
	case "notice", "NOTICE":
		return LevelNotice
	case "warn", "WARN", "warning", "WARNING":
		return LevelWarn
	case "error", "ERROR":
		return LevelError
	case "fatal", "FATAL":
		return LevelFatal
	default:
		return LevelInfo
	}
}
