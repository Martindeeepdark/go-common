package logs

import "context"

// TraceIDKey 是存储 traceId 的标准 context key
// 所有项目使用此 key 存储 traceId，确保 logs 包能正确提取
type TraceIDKey struct{}

// WithTraceID 将 traceId 存入 context
// 返回新的 context，不修改原 context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, TraceIDKey{}, traceID)
}

// TraceID 从 context 提取 traceId
// 如果 context 为 nil 或不包含 traceId，返回空字符串
func TraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(TraceIDKey{}).(string); ok {
		return traceID
	}
	return ""
}
