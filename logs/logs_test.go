package logs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected string
	}{
		{"Trace", LevelTrace, "[Trace]"},
		{"Debug", LevelDebug, "[Debug]"},
		{"Info", LevelInfo, "[Info]"},
		{"Notice", LevelNotice, "[Notice]"},
		{"Warn", LevelWarn, "[Warn]"},
		{"Error", LevelError, "[Error]"},
		{"Fatal", LevelFatal, "[Fatal]"},
		{"Unknown", Level(999), "[?Unknown]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.String())
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Level
	}{
		{"trace lowercase", "trace", LevelTrace},
		{"trace uppercase", "TRACE", LevelTrace},
		{"debug lowercase", "debug", LevelDebug},
		{"debug uppercase", "DEBUG", LevelDebug},
		{"info lowercase", "info", LevelInfo},
		{"info uppercase", "INFO", LevelInfo},
		{"notice lowercase", "notice", LevelNotice},
		{"notice uppercase", "NOTICE", LevelNotice},
		{"warn lowercase", "warn", LevelWarn},
		{"warn uppercase", "WARN", LevelWarn},
		{"warning lowercase", "warning", LevelWarn},
		{"warning uppercase", "WARNING", LevelWarn},
		{"error lowercase", "error", LevelError},
		{"error uppercase", "ERROR", LevelError},
		{"fatal lowercase", "fatal", LevelFatal},
		{"fatal uppercase", "FATAL", LevelFatal},
		{"invalid defaults to info", "invalid", LevelInfo},
		{"empty defaults to info", "", LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ParseLevel(tt.input))
		})
	}
}

func TestInit(t *testing.T) {
	t.Run("initialize with default config", func(t *testing.T) {
		err := Init(LevelInfo)
		assert.NoError(t, err)

		logger := L()
		assert.NotNil(t, logger)
	})

	t.Run("initialize with development mode", func(t *testing.T) {
		err := Init(LevelDebug, WithDevelopment())
		assert.NoError(t, err)

		logger := L()
		assert.NotNil(t, logger)
	})
}

func TestLoggerLevels(t *testing.T) {
	t.Run("log at different levels", func(t *testing.T) {
		err := Init(LevelInfo)
		assert.NoError(t, err)

		logger := L()

		// These should not panic
		logger.Trace("trace message")
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Notice("notice message")
		logger.Warn("warn message")
		logger.Error("error message")
	})
}

func TestLoggerFormatting(t *testing.T) {
	t.Run("formatted logging", func(t *testing.T) {
		err := Init(LevelDebug)
		assert.NoError(t, err)

		logger := L()

		// These should not panic
		logger.Infof("formatted %s", "message")
		logger.Debugf("debug %d", 42)
	})
}

func TestCtxLogger(t *testing.T) {
	t.Run("context logger methods", func(t *testing.T) {
		err := Init(LevelInfo)
		assert.NoError(t, err)

		ctx := context.Background()
		logger := L()

		// These should not panic
		logger.CtxInfof(ctx, "context message: %s", "test")
		logger.CtxErrorf(ctx, "error: %v", "test error")
	})
}

func TestSetLevel(t *testing.T) {
	t.Run("set log level dynamically", func(t *testing.T) {
		err := Init(LevelInfo)
		assert.NoError(t, err)

		logger := L()

		// Set to Error level
		logger.SetLevel(LevelError)
		logger.SetLevel(LevelDebug)
	})
}

func TestSync(t *testing.T) {
	t.Run("sync logger", func(t *testing.T) {
		err := Init(LevelInfo)
		assert.NoError(t, err)

		err = Sync()
		assert.NoError(t, err)
	})
}

func TestGlobalHelperFunctions(t *testing.T) {
	t.Run("global logging functions", func(t *testing.T) {
		err := Init(LevelInfo)
		assert.NoError(t, err)

		// These should not panic
		Trace("trace")
		Debug("debug")
		Info("info")
		Notice("notice")
		Warn("warn")
		Error("error")
	})

	t.Run("global formatted functions", func(t *testing.T) {
		err := Init(LevelInfo)
		assert.NoError(t, err)

		Infof("test %s", "formatted")
	})

	t.Run("global context functions", func(t *testing.T) {
		err := Init(LevelInfo)
		assert.NoError(t, err)

		ctx := context.Background()
		CtxInfof(ctx, "ctx %s", "message")
	})
}

func TestMultipleLoggers(t *testing.T) {
	t.Run("reinitialize logger", func(t *testing.T) {
		// First initialization
		err := Init(LevelInfo)
		assert.NoError(t, err)
		Info("first logger")

		// Second initialization
		err = Init(LevelDebug)
		assert.NoError(t, err)
		Debug("second logger")
	})
}

func TestLazyInitialization(t *testing.T) {
	t.Run("auto-initialize if not initialized", func(t *testing.T) {
		// Reset global logger
		globalLogger = nil

		logger := L()
		assert.NotNil(t, logger)
	})
}
