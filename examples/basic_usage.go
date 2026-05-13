package main

import (
	"github.com/Martindeeepdark/go-common/config"
	"github.com/Martindeeepdark/go-common/errorx"
	"github.com/Martindeeepdark/go-common/eventbus"
	"github.com/Martindeeepdark/go-common/lang/conv"
	"github.com/Martindeeepdark/go-common/lang/ptr"
	"github.com/Martindeeepdark/go-common/lang/slices"
	"github.com/Martindeeepdark/go-common/logs"
	"github.com/Martindeeepdark/go-common/taskgroup"
	"context"
	"errors"
	"fmt"
	"time"
)

// 错误处理示例
func errorExample() {
	fmt.Println("=== 错误处理示例 ===")

	// 注册错误码
	errorx.Register(1001, "用户不存在")
	errorx.Register(1002, "参数错误: {field}")

	// 创建错误
	err := errorx.New(1001)
	fmt.Printf("Error: %v\n", err)

	// 带参数的错误
	err2 := errorx.New(1002, errorx.KV("field", "email"))
	fmt.Printf("Error with param: %v\n", err2)

	// 错误包装
	wrappedErr := errorx.WrapByCode(err2, 1001)
	fmt.Printf("Wrapped error: %v\n", wrappedErr)

	// 检查错误类型
	var se errorx.StatusError
	if errors.As(wrappedErr, &se) {
		fmt.Printf("Code: %d, Msg: %s\n", se.Code(), se.Msg())
	}

	fmt.Println()
}

// 日志示例
func logExample() {
	fmt.Println("=== 日志示例 ===")

	// 初始化日志
	logs.Init(logs.LevelDebug, logs.WithDevelopment())

	// 各种级别的日志
	logs.Trace("This is a trace message")
	logs.Debug("This is a debug message")
	logs.Info("This is an info message")
	logs.Notice("This is a notice message")
	logs.Warn("This is a warning message")
	logs.Error("This is an error message")

	// 格式化日志
	logs.Infof("User %s logged in at %s", "john", time.Now().Format("2006-01-02 15:04:05"))

	// 带上下文的日志
	ctx := context.Background()
	logs.CtxInfof(ctx, "Processing request")

	fmt.Println()
}

// 事件总线示例
func eventbusExample() {
	fmt.Println("=== 事件总线示例 ===")

	bus := eventbus.NewEventBus()

	// 订阅事件
	err := bus.Subscribe("user.created", func(userID string) {
		logs.Infof("User created event received: %s", userID)
	})
	if err != nil {
		return
	}

	// 异步订阅
	err = bus.SubscribeAsync("user.deleted", func(userID string) {
		logs.Infof("User deleted event received: %s", userID)
	}, false)
	if err != nil {
		return
	}

	// 发布事件
	bus.Publish("user.created", "12345")
	bus.Publish("user.deleted", "67890")

	// 等待异步任务完成
	bus.WaitAsync()

	fmt.Println()
}

// 任务组示例
func taskgroupExample() {
	fmt.Println("=== 任务组示例 ===")

	ctx := context.Background()
	g := taskgroup.New(ctx)

	// 添加任务
	g.Go(func() error {
		logs.Info("Task 1 started")
		time.Sleep(100 * time.Millisecond)
		logs.Info("Task 1 completed")
		return nil
	})

	g.Go(func() error {
		logs.Info("Task 2 started")
		time.Sleep(50 * time.Millisecond)
		logs.Info("Task 2 completed")
		return nil
	})

	// 等待所有任务完成
	if err := g.Wait(); err != nil {
		logs.Errorf("Task failed: %v", err)
	}

	fmt.Println()
}

// 语言工具示例
func langExample() {
	fmt.Println("=== 语言工具示例 ===")

	// 指针工具
	name := ptr.ToPtr("Alice")
	fmt.Printf("Pointer: %v, Value: %s\n", name, ptr.Deref(name))

	// 切片工具
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("Numbers: %v\n", numbers)
	fmt.Printf("Contains 3: %v\n", slices.Contains(numbers, 3))
	fmt.Printf("Even numbers: %v\n", slices.Filter(numbers, func(n int) bool {
		return n%2 == 0
	}))
	fmt.Printf("Doubled: %v\n", slices.Map(numbers, func(n int) int {
		return n * 2
	}))

	// 类型转换
	fmt.Printf("String to int: %d\n", conv.ToInt64("123"))
	fmt.Printf("Int to string: %s\n", conv.ToString(456))

	fmt.Println()
}

// 配置管理示例
func configExample() {
	fmt.Println("=== 配置管理示例 ===")

	cfg := config.New()

	// 设置配置值
	cfg.Set("server.port", "8080")
	cfg.Set("server.debug", true)
	cfg.Set("server.timeout", 30)

	// 获取配置值
	port := cfg.GetStringOrDefault("server.port", "8080")
	debug := cfg.GetBoolOrDefault("server.debug", false)
	timeout := cfg.GetIntOrDefault("server.timeout", 30)

	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Debug: %v\n", debug)
	fmt.Printf("Timeout: %d\n", timeout)

	// 检查配置存在
	fmt.Printf("Has server.host: %v\n", cfg.Has("server.host"))

	fmt.Println()
}

func main() {
	errorExample()
	logExample()
	eventbusExample()
	taskgroupExample()
	langExample()
	configExample()

	fmt.Println("All examples completed!")
}
