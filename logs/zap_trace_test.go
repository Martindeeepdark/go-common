package logs

import (
	"context"
	"testing"
)

func TestZapLogger_WithTraceID_Integration(t *testing.T) {
	// 使用默认初始化（输出到 stdout）
	err := Init(LevelInfo)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试：带 traceId 的 context
	ctx := WithTraceID(context.Background(), "test-trace-123")

	// 验证所有 Ctx* 方法都能正常调用且不 panic
	t.Run("CtxTracef", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CtxTracef panicked: %v", r)
			}
		}()
		L().CtxTracef(ctx, "trace message traceId=%s", "test-trace-123")
	})

	t.Run("CtxDebugf", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CtxDebugf panicked: %v", r)
			}
		}()
		L().CtxDebugf(ctx, "debug message")
	})

	t.Run("CtxInfof", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CtxInfof panicked: %v", r)
			}
		}()
		L().CtxInfof(ctx, "info message")
	})

	t.Run("CtxNoticef", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CtxNoticef panicked: %v", r)
			}
		}()
		L().CtxNoticef(ctx, "notice message")
	})

	t.Run("CtxWarnf", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CtxWarnf panicked: %v", r)
			}
		}()
		L().CtxWarnf(ctx, "warn message")
	})

	t.Run("CtxErrorf", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CtxErrorf panicked: %v", r)
			}
		}()
		L().CtxErrorf(ctx, "error message")
	})
}

func TestZapLogger_WithoutTraceID_Integration(t *testing.T) {
	// 使用默认初始化
	err := Init(LevelInfo)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试：不带 traceId 的 context
	ctx := context.Background()

	// 验证不带 traceId 时方法也能正常工作
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logging without traceId panicked: %v", r)
		}
	}()

	L().CtxInfof(ctx, "message without traceId")
}

func TestZapLogger_NilContext(t *testing.T) {
	// 使用默认初始化
	err := Init(LevelInfo)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 测试：nil context
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logging with nil context panicked: %v", r)
		}
	}()

	L().CtxInfof(nil, "message with nil context")
}
