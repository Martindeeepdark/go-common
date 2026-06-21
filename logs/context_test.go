package logs

import (
	"context"
	"testing"
)

func TestWithTraceID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		traceID  string
		expected string
	}{
		{
			name:     "nil context",
			ctx:      nil,
			traceID:  "trace123",
			expected: "trace123",
		},
		{
			name:     "background context",
			ctx:      context.Background(),
			traceID:  "trace456",
			expected: "trace456",
		},
		{
			name:     "empty traceID",
			ctx:      context.Background(),
			traceID:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithTraceID(tt.ctx, tt.traceID)
			got := TraceID(ctx)
			if got != tt.expected {
				t.Errorf("TraceID() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTraceID_NoValue(t *testing.T) {
	ctx := context.Background()
	got := TraceID(ctx)
	if got != "" {
		t.Errorf("TraceID() = %v, want empty string", got)
	}
}

func TestTraceID_NilContext(t *testing.T) {
	got := TraceID(nil)
	if got != "" {
		t.Errorf("TraceID(nil) = %v, want empty string", got)
	}
}
