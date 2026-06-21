package logs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger FullLogger
)

// zapLogger implements FullLogger interface using zap
type zapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	level  zap.AtomicLevel
}

// Init initializes the global logger with zap
func Init(level Level, opts ...Option) error {
	zapLevel := zapcore.Level(level)
	atomicLevel := zap.NewAtomicLevelAt(zapLevel)

	// Default encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Apply options
	config := &Config{
		Level:         level,
		EncoderConfig: encoderConfig,
		AtomicLevel:   atomicLevel,
	}

	for _, opt := range opts {
		opt(config)
	}

	// Set default encoder if not set
	if config.Encoder == nil {
		config.Encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Set default output if not set
	if config.Output == nil {
		config.Output = zapcore.AddSync(zapcore.Lock(os.Stdout))
	}

	// Build core
	core := zapcore.NewCore(
		config.Encoder,
		config.Output,
		atomicLevel,
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	z := &zapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
		level:  atomicLevel,
	}

	globalLogger = z
	return nil
}

// Config for logger initialization
type Config struct {
	Level         Level
	Output        zapcore.WriteSyncer
	Encoder       zapcore.Encoder
	EncoderConfig zapcore.EncoderConfig
	AtomicLevel   zap.AtomicLevel
}

// Option for logger configuration
type Option func(*Config)

// WithOutput sets the output destination
func WithOutput(output zapcore.WriteSyncer) Option {
	return func(c *Config) {
		c.Output = output
	}
}

// WithEncoder sets the encoder
func WithEncoder(encoder zapcore.Encoder) Option {
	return func(c *Config) {
		c.Encoder = encoder
	}
}

// WithEncoderConfig sets the encoder config
func WithEncoderConfig(encoderConfig zapcore.EncoderConfig) Option {
	return func(c *Config) {
		c.EncoderConfig = encoderConfig
		c.Encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
}

// WithDevelopment enables development mode (with debug logging)
func WithDevelopment() Option {
	return func(c *Config) {
		c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
}

// L returns the global logger instance
func L() FullLogger {
	if globalLogger == nil {
		// Initialize with default config if not initialized
		_ = Init(LevelInfo)
	}
	return globalLogger
}

// Trace logs a message at TraceLevel
func (z *zapLogger) Trace(v ...interface{}) {
	z.sugar.Debug(fmt.Sprint(v...))
}

// Debug logs a message at DebugLevel
func (z *zapLogger) Debug(v ...interface{}) {
	z.sugar.Debug(v...)
}

// Info logs a message at InfoLevel
func (z *zapLogger) Info(v ...interface{}) {
	z.sugar.Info(v...)
}

// Notice logs a message at NoticeLevel (mapped to Warn in zap)
func (z *zapLogger) Notice(v ...interface{}) {
	z.sugar.Warn(v...)
}

// Warn logs a message at WarnLevel
func (z *zapLogger) Warn(v ...interface{}) {
	z.sugar.Warn(v...)
}

// Error logs a message at ErrorLevel
func (z *zapLogger) Error(v ...interface{}) {
	z.sugar.Error(v...)
}

// Fatal logs a message at FatalLevel and exits
func (z *zapLogger) Fatal(v ...interface{}) {
	z.sugar.Fatal(v...)
}

// Tracef logs a formatted message at TraceLevel
func (z *zapLogger) Tracef(format string, v ...interface{}) {
	z.sugar.Debugf(format, v...)
}

// Debugf logs a formatted message at DebugLevel
func (z *zapLogger) Debugf(format string, v ...interface{}) {
	z.sugar.Debugf(format, v...)
}

// Infof logs a formatted message at InfoLevel
func (z *zapLogger) Infof(format string, v ...interface{}) {
	z.sugar.Infof(format, v...)
}

// Noticef logs a formatted message at NoticeLevel
func (z *zapLogger) Noticef(format string, v ...interface{}) {
	z.sugar.Warnf(format, v...)
}

// Warnf logs a formatted message at WarnLevel
func (z *zapLogger) Warnf(format string, v ...interface{}) {
	z.sugar.Warnf(format, v...)
}

// Errorf logs a formatted message at ErrorLevel
func (z *zapLogger) Errorf(format string, v ...interface{}) {
	z.sugar.Errorf(format, v...)
}

// Fatalf logs a formatted message at FatalLevel and exits
func (z *zapLogger) Fatalf(format string, v ...interface{}) {
	z.sugar.Fatalf(format, v...)
}

// withTraceID 从 context 中提取 traceId 并返回带 traceId 字段的 logger
// 如果 context 为 nil 或没有 traceId，返回原 logger
func (z *zapLogger) withTraceID(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return z.sugar
	}

	traceID := TraceID(ctx)
	if traceID == "" {
		return z.sugar
	}

	// 返回带 traceId 字段的子 logger
	return z.sugar.With("traceId", traceID)
}

// CtxTracef logs a formatted message with context at TraceLevel
func (z *zapLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	z.withTraceID(ctx).Debugf(format, v...)
}

// CtxDebugf logs a formatted message with context at DebugLevel
func (z *zapLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	z.withTraceID(ctx).Debugf(format, v...)
}

// CtxInfof logs a formatted message with context at InfoLevel
func (z *zapLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	z.withTraceID(ctx).Infof(format, v...)
}

// CtxNoticef logs a formatted message with context at NoticeLevel
func (z *zapLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	z.withTraceID(ctx).Warnf(format, v...)
}

// CtxWarnf logs a formatted message with context at WarnLevel
func (z *zapLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	z.withTraceID(ctx).Warnf(format, v...)
}

// CtxErrorf logs a formatted message with context at ErrorLevel
func (z *zapLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	z.withTraceID(ctx).Errorf(format, v...)
}

// CtxFatalf logs a formatted message with context at FatalLevel and exits
func (z *zapLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	z.withTraceID(ctx).Fatalf(format, v...)
}

// SetLevel sets the logging level
func (z *zapLogger) SetLevel(level Level) {
	z.level.SetLevel(zapcore.Level(level))
}

// SetOutput sets the output destination
func (z *zapLogger) SetOutput(output interface{}) {
	// This would require recreating the logger, simplified for now
	// In production, you might want to implement this properly
}

// Sync flushes any buffered log entries
func Sync() error {
	if globalLogger != nil {
		if z, ok := globalLogger.(*zapLogger); ok {
			if err := z.logger.Sync(); err != nil {
				if isIgnorableSyncError(err) {
					return nil
				}
				return err
			}
			return nil
		}
	}
	return nil
}

func isIgnorableSyncError(err error) bool {
	if err == nil {
		return true
	}
	if errors.Is(err, syscall.EINVAL) {
		return true
	}
	msg := err.Error()
	if strings.Contains(msg, "invalid argument") {
		return true
	}
	return false
}

// Helper functions for quick access

func Trace(v ...interface{}) {
	L().Trace(v...)
}

func Debug(v ...interface{}) {
	L().Debug(v...)
}

func Info(v ...interface{}) {
	L().Info(v...)
}

func Notice(v ...interface{}) {
	L().Notice(v...)
}

func Warn(v ...interface{}) {
	L().Warn(v...)
}

func Error(v ...interface{}) {
	L().Error(v...)
}

func Fatal(v ...interface{}) {
	L().Fatal(v...)
}

func Tracef(format string, v ...interface{}) {
	L().Tracef(format, v...)
}

func Debugf(format string, v ...interface{}) {
	L().Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	L().Infof(format, v...)
}

func Noticef(format string, v ...interface{}) {
	L().Noticef(format, v...)
}

func Warnf(format string, v ...interface{}) {
	L().Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	L().Errorf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	L().Fatalf(format, v...)
}

func CtxTracef(ctx context.Context, format string, v ...interface{}) {
	L().CtxTracef(ctx, format, v...)
}

func CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	L().CtxDebugf(ctx, format, v...)
}

func CtxInfof(ctx context.Context, format string, v ...interface{}) {
	L().CtxInfof(ctx, format, v...)
}

func CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	L().CtxNoticef(ctx, format, v...)
}

func CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	L().CtxWarnf(ctx, format, v...)
}

func CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	L().CtxErrorf(ctx, format, v...)
}

func CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	L().CtxFatalf(ctx, format, v...)
}

// GetSlogLogger 返回一个 *slog.Logger 用于向后兼容
// 这是一个适配层，底层仍使用 zap，但提供 slog 接口给需要的基础设施代码
func GetSlogLogger() *slog.Logger {
	if globalLogger == nil {
		_ = Init(LevelInfo)
	}

	z, ok := globalLogger.(*zapLogger)
	if !ok {
		// 降级：返回默认 slog logger
		return slog.Default()
	}

	// 创建 slog.Handler 适配器包装 zap
	handler := &zapSlogHandler{sugar: z.sugar}
	return slog.New(handler)
}

// zapSlogHandler 是 slog.Handler 的实现，底层使用 zap SugaredLogger
type zapSlogHandler struct {
	sugar *zap.SugaredLogger
}

func (h *zapSlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	// 始终返回 true，让 zap 的 level 控制来决定
	return true
}

func (h *zapSlogHandler) Handle(_ context.Context, record slog.Record) error {
	// 收集所有 attributes
	var args []interface{}
	record.Attrs(func(attr slog.Attr) bool {
		args = append(args, attr.Key, attr.Value.Any())
		return true
	})

	// 根据 level 调用对应的 zap 方法
	switch record.Level {
	case slog.LevelDebug:
		h.sugar.Debugw(record.Message, args...)
	case slog.LevelInfo:
		h.sugar.Infow(record.Message, args...)
	case slog.LevelWarn:
		h.sugar.Warnw(record.Message, args...)
	case slog.LevelError:
		h.sugar.Errorw(record.Message, args...)
	default:
		h.sugar.Infow(record.Message, args...)
	}
	return nil
}

func (h *zapSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// 将 attrs 转换为 zap fields
	var args []interface{}
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value.Any())
	}
	return &zapSlogHandler{sugar: h.sugar.With(args...)}
}

func (h *zapSlogHandler) WithGroup(name string) slog.Handler {
	// zap 不直接支持 group，简单实现：将 group name 作为前缀
	return &zapSlogHandler{sugar: h.sugar.Named(name)}
}
