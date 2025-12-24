# Common - Go 通用工具库

一个轻量、通用、框架无关的 Go 工具库，为脚手架提供基础设施支持。

## ✨ 特性

- ✅ **错误处理** (errorx) - 集中式错误码管理，支持堆栈跟踪
- ✅ **日志系统** (logs) - 基于 zap 的高性能结构化日志
- ✅ **上下文缓存** (ctxcache) - 线程安全的上下文缓存
- ✅ **事件总线** (eventbus) - 本地事件总线 + 消息队列接口（支持 NSQ/Kafka/RabbitMQ/Pulsar/NATS）
- ✅ **数据库抽象** (database) - SQL 构建器、事务管理、分页工具
- ✅ **缓存抽象** (cache) - 统一的缓存接口（简单缓存、Hash、List、Set、有序集合）
- ✅ **配置管理** (config) - YAML 配置文件加载
- ✅ **任务组** (taskgroup) - Goroutine 管理和错误处理
- ✅ **语言工具** (lang) - 丰富的类型转换、切片、Map 等工具函数
- ✅ **Mock 支持** - 集成 gomock 用于单元测试

## 🎯 设计理念

### 职责分工

**Common 包**（当前项目）：
- ✅ 框架无关的通用工具
- ✅ 定义标准接口（database、cache、eventbus）
- ✅ 提供通用工具函数

**脚手架**（独立项目）：
- ✅ 框架特定的实现（Gin、Echo、Fiber）
- ✅ 具体数据库驱动（GORM + MySQL、pgx + PostgreSQL）
- ✅ 具体缓存实现（go-redis、freecache）
- ✅ HTTP 中间件和路由增强

## 📦 安装

```bash
go get common
```

## 🚀 快速开始

### 错误处理 (errorx)

```go
import "common/errorx"

errorx.Register(1001, "用户不存在")
err := errorx.New(1001, errorx.KV("user_id", "123"))
```

### 日志系统 (logs)

```go
import "common/logs"

logs.Init(logs.LevelInfo, logs.WithDevelopment())
logs.Info("Server started")
logs.CtxInfof(ctx, "Processing request")
```

### 数据库工具 (database)

#### SQL 构建器

```go
import "common/database/sql"

builder := sql.New().
    Select("id", "name", "email").
    From("users").
    Where("status = ?", "active").
    Where("age > ?", 18).
    OrderByDesc("created_at").
    Limit(10)

query, args := builder.Build()
// SELECT id, name, email FROM users WHERE status = ? AND age > ? ORDER BY created_at DESC LIMIT 10
```

#### 分页工具

```go
import "common/database/sql"

pager := sql.NewPager(1, 20) // 第1页，每页20条
pager.SetTotal(100)          // 总共100条记录

offset := pager.Offset()     // 0
limit := pager.Limit()       // 20

hasNext := pager.HasNext()   // true
hasPrev := pager.HasPrev()   // false
```

#### 事务管理

```go
import "common/database/sql"

// 自动事务管理
err := sql.TxManager.ExecuteInTx(ctx, db, func(ctx context.Context) error {
    // 在事务中执行操作
    // 如果返回错误，会自动回滚
    // 否则自动提交
    return nil
})

// 从上下文获取事务
tx := sql.TxManager.MustFromContext(ctx)
tx.Exec(ctx, "UPDATE users SET name = ?", "John")
```

### 缓存接口 (cache)

```go
import "common/cache/defs"

// 脚手架实现这个接口
type MyCache struct{}

func (c *MyCache) Get(ctx context.Context, key string, dest interface{}) error {
    // 实现获取逻辑
    return nil
}

func (c *MyCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    // 实现设置逻辑
    return nil
}

// ... 其他接口方法
```

**支持的缓存类型**：
- `Cache` - 简单 KV 缓存
- `HashCache` - Hash 结构（类似 Redis HGET/HSET）
- `ListCache` - 列表结构（类似 Redis LPUSH/RPOP）
- `SetCache` - 集合结构（类似 Redis SADD/SMEMBERS）
- `SortedSetCache` - 有序集合（类似 Redis ZADD/ZRANGE）
- `Lock` - 分布式锁

### 事件总线 (eventbus)

#### 本地事件总线

```go
import "common/eventbus"

bus := eventbus.NewEventBus()
bus.Subscribe("user.created", func(userID string) {
    fmt.Println("User created:", userID)
})
bus.Publish("user.created", "123")
```

#### 消息队列接口

```go
import (
    "common/eventbus"
    "common/eventbus/impl"
)

consumerSvc := impl.NewConsumerService()
eventbus.SetDefaultSVC(consumerSvc)

producer, _ := impl.NewProducer("localhost:4150", "topic", "group", 3)
producer.Send(ctx, []byte("message"))
```

### 语言工具 (lang)

```go
import (
    "common/lang/ptr"
    "common/lang/slices"
    "common/lang/conv"
)

p := ptr.ToPtr("hello")
s := ptr.Deref(p)

nums := []int{1, 2, 3, 4, 5}
even := slices.Filter(nums, func(n int) bool { return n%2 == 0 })

str := conv.ToString(123)
num := conv.ToInt64("456")
```

## 📁 项目结构

```
common/
├── errorx/           # 错误处理系统
│   └── internal/
├── logs/             # 日志系统
├── ctxcache/         # 上下文缓存
├── eventbus/         # 事件总线
│   ├── impl/         # 消息队列实现（NSQ/Kafka/RMQ/Pulsar/NATS）
│   ├── local.go      # 本地事件总线
│   └── eventbus.go   # 消息队列接口
├── database/         # 数据库抽象
│   ├── defs/         # 数据库接口定义
│   └── sql/          # SQL 工具（构建器、分页、事务）
├── cache/            # 缓存抽象
│   └── defs/         # 缓存接口定义
├── config/           # 配置管理
├── taskgroup/        # 任务组管理
├── lang/             # 语言工具
│   ├── conv/         # 类型转换
│   ├── ptr/          # 指针工具
│   ├── slices/       # 切片工具
│   ├── maps/         # Map 工具
│   ├── ternary/      # 三元操作
│   └── crypto/       # 加密工具
└── examples/         # 使用示例
```

## 🔧 模块说明

### database - 数据库抽象

**defs 包**：定义数据库接口
- `Database` - 基础数据库接口
- `SQLDatabase` - SQL 数据库操作接口
- `Transaction` - 事务接口
- `Rows`、`Result` - 查询结果接口

**sql 包**：通用 SQL 工具
- `Builder` - 链式 SQL 构建器
- `Pager` - 分页工具
- `TxManager` - 事务管理器

### cache - 缓存抽象

**defs 包**：定义缓存接口
- `Cache` - 简单缓存接口
- `HashCache` - Hash 操作接口
- `ListCache` - 列表操作接口
- `SetCache` - 集合操作接口
- `SortedSetCache` - 有序集合接口
- `Lock` - 分布式锁接口
- `StatsProvider` - 缓存统计接口

**脚手架实现示例**：
- Gin + GORM + MySQL + go-redis
- Gin + pgx + PostgreSQL + freecache
- Echo + GORM + PostgreSQL + bigcache

## 📝 依赖

- Go 1.25+
- go.uber.org/mock v0.4.0
- go.uber.org/zap v1.27.0
- gopkg.in/yaml.v3 v3.0.1

## 🎯 使用原则

1. **框架无关** - Common 包不依赖任何 Web 框架
2. **接口优先** - 定义接口，由脚手架实现
3. **通用工具** - 提供可在任何项目中使用的工具
4. **易于测试** - 集成 gomock，方便单元测试

## 🔄 脚手架集成

脚手架应该：

1. **实现接口** - 实现 common 定义的 database、cache 接口
2. **添加框架支持** - Gin、Echo 等特定框架的工具
3. **集成 ORM** - GORM、pgx 等 ORM 框架
4. **提供中间件** - HTTP 中间件、路由增强
5. **配置管理** - 开发环境、生产环境配置

## 📄 License

MIT License
