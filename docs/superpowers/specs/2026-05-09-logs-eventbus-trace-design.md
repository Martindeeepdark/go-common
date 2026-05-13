# Logs + Trace 打通 & EventBus + Trace 集成

## Context

common 包已有 `logs`（zap 封装）和 `trace/otel`（OpenTelemetry），但两者互不感知：
- `CtxXxx(ctx, ...)` 方法接收 ctx 但完全没使用
- eventbus 不支持 context，handler 里无法串联 trace

目标：让日志自动携带 traceID，让 eventbus handler 能继承上游 trace 链路。

## 一、Logs + Trace

### 改动文件
- `logs/zap.go` — 修改 CtxXxx 方法实现

### 实现

新增内部 helper：

```go
func (z *zapLogger) ctxFields(ctx context.Context) []zap.Field {
    var fields []zap.Field
    if ctx != nil {
        span := trace.SpanFromContext(ctx)
        if sc := span.SpanContext(); sc.HasTraceID() {
            fields = append(fields, zap.String("traceID", sc.TraceID().String()))
        }
        if sc.HasSpanID() {
            fields = append(fields, zap.String("spanID", sc.SpanID().String()))
        }
    }
    return fields
}
```

所有 CtxXxx 方法改为：

```go
func (z *zapLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
    if fields := z.ctxFields(ctx); len(fields) > 0 {
        z.sugar.With(fields...).Infof(format, v...)
        return
    }
    z.sugar.Infof(format, v...)
}
```

### 依赖方向
- `logs` 包 import `go.opentelemetry.io/otel/trace`（只依赖 OTEL trace API，不依赖 common/trace）
- `trace/otel` 不 import `logs`（保持单向）

### 行为
- `Info()` / `Infof()` — 不变，无 traceID
- `CtxInfof(ctx, ...)` — ctx 有 span → 自动带 traceID/spanID；ctx 无 span → 退化为普通 Infof
- 接口签名零改动，调用方零修改

### 日志输出

```json
{"time":"2026-05-09T12:00:00Z","level":"INFO","caller":"service/order.go:42","traceID":"a1b2c3d4e5f6...","spanID":"7890...","msg":"order 42 created"}
```

## 二、EventBus + Trace

### 改动文件
- `eventbus/local.go` — 新增 PublishContext，修改反射调用逻辑

### PublishContext

```go
func (bus *EventBus) PublishContext(ctx context.Context, topic string, args ...interface{}) {
    // 创建 child span
    var span trace.Span
    if ctx != nil {
        tp := otel.GetTracerProvider().Tracer("eventbus")
        ctx, span = tp.Start(ctx, "eventbus: "+topic)
        defer span.End()
    }

    bus.lock.Lock()
    defer bus.lock.Unlock()

    if handlers, ok := bus.handlers[topic]; ok && len(handlers) > 0 {
        copyHandlers := make([]*eventHandler, len(handlers))
        copy(copyHandlers, handlers)

        for i, handler := range copyHandlers {
            if handler.flagOnce {
                bus.removeHandler(topic, i)
            }
            if !handler.async {
                bus.doPublishWithContext(ctx, handler, topic, args...)
            } else {
                bus.wg.Add(1)
                go bus.doPublishAsyncWithContext(ctx, handler, topic, args...)
            }
        }
    }
}
```

### Handler 自动识别 ctx

反射调用时检查 handler 第一个参数类型：

```go
func (bus *EventBus) setUpPublishWithContext(ctx context.Context, handler *eventHandler, args ...interface{}) []reflect.Value {
    funcType := handler.callBack.Type()
    needsCtx := funcType.NumIn() > 0 && funcType.In(0) == reflect.TypeOf((*context.Context)(nil)).Elem()

    if needsCtx {
        passedArguments := make([]reflect.Value, 0, 1+len(args))
        passedArguments = append(passedArguments, reflect.ValueOf(ctx))
        for i, v := range args {
            if v == nil {
                passedArguments = append(passedArguments, reflect.New(funcType.In(i+1)).Elem())
            } else {
                passedArguments = append(passedArguments, reflect.ValueOf(v))
            }
        }
        return passedArguments
    }
    // 旧 handler，不带 ctx
    return bus.setUpPublish(handler, args...)
}
```

### 行为

```go
// 旧 handler — 继续工作
bus.Subscribe("order.created", func(order Order) { ... })

// 新 handler — 自动接收 ctx
bus.Subscribe("order.created", func(ctx context.Context, order Order) {
    logs.CtxInfof(ctx, "handling order")  // traceID 自动带上
})
```

### 依赖
- `eventbus` import `context`, `reflect`, `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/trace`
- 不依赖 `logs` 包

## 三、不做的事

- database/sql 不加 trace
- 不改 `Info()` 系列为 ctx 版本
- 不改 `eventbus/eventbus.go` 的 Producer/Consumer 接口（那是 MQ 层的）

## 四、测试计划

| 包 | 测试内容 |
|---|---|
| logs | CtxXxx 方法带 span 的 ctx → 日志输出含 traceID |
| logs | CtxXxx 方法不带 span 的 ctx → 日志输出无 traceID |
| logs | 非 Ctx 方法 → 行为不变 |
| eventbus | PublishContext + ctx handler → handler 收到 ctx 含 traceID |
| eventbus | PublishContext + 旧 handler → 旧 handler 正常工作 |
| eventbus | PublishContext + async handler → ctx 正确传递 |
