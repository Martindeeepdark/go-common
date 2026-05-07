# Common（Go 通用工具库）Code Wiki

## 1. 项目概览

**定位**
- 该仓库是一个 Go module（`module common`），提供框架无关的基础设施能力与通用工具函数，主要面向“脚手架/业务项目”复用。
- 代码以“接口定义 + 通用实现（少量）+ 使用示例”的形态组织，其中 `cache/defs`、`database/defs`、`eventbus` 的 MQ 部分偏向抽象接口，具体落地实现通常由脚手架/业务侧提供。

**入口**
- 仓库本身是库项目；可执行入口用于演示： [basic_usage.go](file:///workspace/examples/basic_usage.go#L1-L195)

**外部依赖（go.mod）**
- 日志：`go.uber.org/zap`（[go.mod](file:///workspace/go.mod#L1-L15)）
- 配置：`gopkg.in/yaml.v3`（[go.mod](file:///workspace/go.mod#L1-L15)）
- 测试：`github.com/stretchr/testify`（[go.mod](file:///workspace/go.mod#L1-L15)）

**Go 版本**
- `go.mod` 声明：`go 1.24.0`（[go.mod](file:///workspace/go.mod#L1-L4)）
- CI（GitHub Actions）使用：Go 1.21（[test.yml](file:///workspace/.github/workflows/test.yml#L14-L19)）

## 2. 仓库结构

来自 [README.md](file:///workspace/README.md#L189-L216) 的结构基础上补充关键文件：

```
common/
├── cache/
│   └── defs/                 # 缓存接口定义（KV/Hash/List/Set/ZSet/Lock/Stats）
├── config/                   # YAML 配置加载与全局默认实例
├── ctxcache/                 # 基于 context 的线程安全缓存（sync.Map）
├── database/
│   ├── defs/                 # 数据库接口定义（Database/SQLDatabase/Tx/Rows/Result）
│   └── sql/                  # SQL Builder、分页器、事务上下文管理（TxManager）
├── errorx/
│   └── internal/             # 错误码注册、StatusError 实现、堆栈采集
├── eventbus/
│   ├── impl/                 # MQ 工厂/占位实现（由脚手架提供具体实现）
│   └── local.go              # 本地事件总线实现（反射调用 + 异步 WaitGroup）
├── lang/
│   ├── conv/                 # 通用类型转换
│   ├── crypto/               # md5 等工具
│   ├── maps/                 # map 工具
│   ├── ptr/                  # 指针工具
│   ├── slices/               # slice 工具（泛型）
│   └── ternary/              # “三元/短路”工具（泛型）
├── logs/                     # 基于 zap 的全局日志器封装
├── taskgroup/                # goroutine 任务组（错误收集 + ctx cancel）
├── examples/                 # 示例入口 main
└── .github/workflows/test.yml # CI：测试 + 覆盖率 + lint
```

## 3. 整体架构与设计原则

**3.1 分层与边界**
- **纯工具包**：`lang/*` 提供无副作用的泛型/函数式工具。
- **基础设施能力（通用实现）**：`errorx`、`logs`、`config`、`ctxcache`、`taskgroup`、`eventbus/local`、`database/sql`。
- **抽象接口层（供外部实现）**：`cache/defs`、`database/defs`、`eventbus`（Producer/ConsumerService/ConsumerHandler）。
- **示例层**：`examples` 仅用于展示 API 的组合用法。

**3.2 全局单例（Global Defaults）**
- 日志全局 logger：`logs.L()` 背后是 `globalLogger`（[zap.go](file:///workspace/logs/zap.go#L12-L131)）
- 配置默认实例：`config` 包级 `defaultConfig`（[config.go](file:///workspace/config/config.go#L262-L293)）
- MQ 消费服务默认实例：`eventbus.defaultSVC`（[eventbus.go](file:///workspace/eventbus/eventbus.go#L62-L73)）
- 事务管理器：`sql.TxManager`（[tx.go](file:///workspace/database/sql/tx.go#L16-L21)）

这些全局单例简化集成成本，但也意味着：
- 测试中建议显式初始化/重置（例如 `logs.Init(...)`、`config.Clear()`）。
- 多实例隔离场景（多租户/多 logger）需要业务侧自行封装或避免使用默认实例。

**3.3 包级依赖关系（仅列出仓库内部引用）**
- `common/errorx` → `common/errorx/internal`（[error.go](file:///workspace/errorx/error.go#L1-L91)）
- `common/database/sql` → `common/database/defs`（[tx.go](file:///workspace/database/sql/tx.go#L3-L9)）
- `common/eventbus/impl` → `common/eventbus`（[impl.go](file:///workspace/eventbus/impl/impl.go#L3-L8)）
- `common/examples` → 依赖所有功能包（[basic_usage.go](file:///workspace/examples/basic_usage.go#L3-L16)）

外部依赖集中在：
- `common/logs` → `zap`（[zap.go](file:///workspace/logs/zap.go#L8-L10)）
- `common/config` → `yaml.v3`（[config.go](file:///workspace/config/config.go#L3-L9)）

## 4. 模块 Wiki

### 4.1 errorx：错误码 + 可追踪错误（StatusError）

**目标**
- 为业务提供“错误码（Code）+ 文本（Msg）+ 扩展字段（Extra）+ 可选堆栈（stack）”的统一错误对象。

**对外 API（common/errorx）**
- `type StatusError`：统一错误接口（[error.go](file:///workspace/errorx/error.go#L13-L19)）
- `Register(code, msg, ...)`：注册错误码与模板消息（[error.go](file:///workspace/errorx/error.go#L78-L81)）
- `New(code, options...)`：创建带堆栈的错误（[error.go](file:///workspace/errorx/error.go#L39-L43)）
- `WrapByCode(err, code, options...)`：包装已有错误并附加错误码与（必要时）堆栈（[error.go](file:///workspace/errorx/error.go#L45-L53)）
- `KV/KVf`：用于消息模板占位符替换（如 `{field}`）（[error.go](file:///workspace/errorx/error.go#L24-L32)）
- `Extra`：附加扩展信息键值（[error.go](file:///workspace/errorx/error.go#L34-L37)）
- `ErrorWithoutStack(err)`：去掉 `stack=` 之后的内容，便于对外输出（[error.go](file:///workspace/errorx/error.go#L65-L76)）

**实现要点（common/errorx/internal）**
- 错误码注册表：`CodeDefinitions`（[register.go](file:///workspace/errorx/internal/register.go#L8-L11)）
- `withStatus` 实现了 `error`，并支持 `errors.Is/As/Unwrap`（[status.go](file:///workspace/errorx/internal/status.go#L24-L97)）
- 堆栈采集：`runtime.Stack`（[stack.go](file:///workspace/errorx/internal/stack.go#L12-L36)）
- `WrapByCode` 会检测被包装 error 是否已有 `StackTrace()`，避免重复采集（[status.go](file:///workspace/errorx/internal/status.go#L159-L167)）

**典型用法**
- 见示例：[basic_usage.go](file:///workspace/examples/basic_usage.go#L19-L45)

### 4.2 logs：基于 zap 的全局结构化日志

**目标**
- 提供统一的 `FullLogger` 接口（级别日志 + 格式化日志 + 带 context 的格式化日志 + 控制接口），并用 zap 实现默认落地。

**关键类型**
- `type FullLogger`：组合接口（[logger.go](file:///workspace/logs/logger.go#L45-L51)）
- `type Level`：日志级别枚举（[logger.go](file:///workspace/logs/logger.go#L53-L67)）
- `ParseLevel(string) Level`：解析字符串级别（[logger.go](file:///workspace/logs/logger.go#L91-L110)）

**关键函数**
- `Init(level, opts...)`：初始化全局 logger（[zap.go](file:///workspace/logs/zap.go#L23-L81)）
- `WithDevelopment()`：开发模式（彩色 level 等）（[zap.go](file:///workspace/logs/zap.go#L117-L122)）
- `L()`：获取全局 logger（未 init 会默认 `LevelInfo` 初始化）（[zap.go](file:///workspace/logs/zap.go#L124-L131)）
- `Sync()`：flush（[zap.go](file:///workspace/logs/zap.go#L249-L257)）

**注意点**
- `SetOutput` 目前是空实现占位（[zap.go](file:///workspace/logs/zap.go#L243-L247)），如果需要动态切换输出，需要业务侧扩展实现方式（例如重建 zap core）。

**典型用法**
- 见示例：[basic_usage.go](file:///workspace/examples/basic_usage.go#L47-L70)

### 4.3 config：YAML 配置与线程安全读写

**能力范围**
- 读 YAML 文件 → `map[string]interface{}`，提供线程安全的 `Get/Set/Keys/Merge` 等操作，并带有默认全局实例。

**关键类型与函数**
- `type Config`：内部使用 `sync.RWMutex` 保护 `data`（[config.go](file:///workspace/config/config.go#L11-L15)）
- `LoadFromFile(path)`：从 YAML 文件加载（[config.go](file:///workspace/config/config.go#L24-L41)）
- `GetString/GetInt/GetBool` 与对应 `OrDefault`（[config.go](file:///workspace/config/config.go#L99-L187)）
- 默认实例：`LoadDefault/Get/Set/...`（[config.go](file:///workspace/config/config.go#L262-L293)）

**限制**
- `LoadFromEnv(prefix)` 中的 `parseEnvVar` 目前是未完成占位实现（[config.go](file:///workspace/config/config.go#L43-L64)），如果依赖环境变量加载，需要自行补齐解析逻辑或在业务侧实现。

**典型用法**
- 见示例：[basic_usage.go](file:///workspace/examples/basic_usage.go#L160-L184)

### 4.4 ctxcache：绑定在 context 上的临时缓存（sync.Map）

**目标**
- 为请求链路（context 生命周期）提供线程安全的临时 KV 缓存，适合存放一次请求中重复使用的中间结果（如解析后的 token、DB 查询的辅助结果等）。

**核心 API**
- `Init(ctx)`：在 context 中挂载 `*sync.Map`（[ctx_cache.go](file:///workspace/ctxcache/ctx_cache.go#L10-L13)）
- `Get[T]/Store/HasKey/LoadOrStore/LoadAndDelete/Delete/Keys/Clear`（[ctx_cache.go](file:///workspace/ctxcache/ctx_cache.go#L15-L123)）

**使用注意**
- 只有调用 `Init` 之后缓存相关函数才会生效（否则会返回未命中/不写入）。

### 4.5 taskgroup：并发任务编排（错误收集 + ctx cancel）

**目标**
- 在并发执行多个任务时：收集错误、在出现错误时触发取消（`context.WithCancel`）、并提供 `Wait/WaitAll` 的汇总接口。

**关键类型与函数**
- `New(ctx)`：创建任务组并生成可取消的 context（[taskgroup.go](file:///workspace/taskgroup/taskgroup.go#L17-L25)）
- `Go(fn)` / `GoWithContext(fn)`：启动 goroutine，任务失败会 `cancel()`（[taskgroup.go](file:///workspace/taskgroup/taskgroup.go#L27-L61)）
- `Wait()`：返回第一个错误（[taskgroup.go](file:///workspace/taskgroup/taskgroup.go#L63-L72)）
- `WaitAll()`：返回所有错误（[taskgroup.go](file:///workspace/taskgroup/taskgroup.go#L74-L78)）
- `Parallel/Serial`：便捷并发/串行编排（[taskgroup.go](file:///workspace/taskgroup/taskgroup.go#L108-L154)）

### 4.6 eventbus：本地事件总线 + MQ 接口抽象

**两类能力**
- 本地事件总线：`eventbus.NewEventBus()` 返回 `Bus`，支持同步/异步/Once 订阅与发布（[local.go](file:///workspace/eventbus/local.go#L9-L234)）
- MQ 抽象接口：`Producer/ConsumerService/ConsumerHandler` 与默认 `defaultSVC`（[eventbus.go](file:///workspace/eventbus/eventbus.go#L5-L73)）

**本地事件总线关键机制**
- 订阅函数必须是函数类型，通过 `reflect.Value.Call` 触发（[local.go](file:///workspace/eventbus/local.go#L66-L71), [local.go](file:///workspace/eventbus/local.go#L174-L177)）
- 异步订阅会用 `WaitGroup` 跟踪，需要调用 `WaitAsync()` 等待完成（[local.go](file:///workspace/eventbus/local.go#L167-L170), [local.go](file:///workspace/eventbus/local.go#L231-L234)）
- `transactional=true` 的异步 handler 会对该 handler 加锁串行执行（[local.go](file:///workspace/eventbus/local.go#L179-L187)）

**MQ impl 现状**
- `eventbus/impl` 当前仅提供“工厂占位”，会基于 `COMMON_MQ_TYPE` 返回“需要脚手架提供实现”的错误（[impl.go](file:///workspace/eventbus/impl/impl.go#L10-L53)）
- 具体实现方式参考：[eventbus/impl/README.md](file:///workspace/eventbus/impl/README.md)

### 4.7 database：数据库接口抽象 + SQL 工具

**接口层（database/defs）**
- `Database`：连接生命周期与健康检查（[database.go](file:///workspace/database/defs/database.go#L7-L15)）
- `SQLDatabase`：Query/Exec/BeginTx（[database.go](file:///workspace/database/defs/database.go#L17-L27)）
- `Transaction/Rows/Result`：抽象事务与结果（[database.go](file:///workspace/database/defs/database.go#L29-L60)）
- `TxOptions/IsolationLevel`：事务选项（[database.go](file:///workspace/database/defs/database.go#L62-L82)）

**通用工具层（database/sql）**
- SQL Builder：`sql.New()` + `Select/From/Where/Join/OrderBy.../Build()`（[builder.go](file:///workspace/database/sql/builder.go#L8-L181)）
- Pager：`NewPager/SetTotal/Offset/Limit/HasNext/HasPrev`（[pager.go](file:///workspace/database/sql/pager.go#L8-L92)）
- TxManager：在 context 中存取事务对象，并提供 `ExecuteInTx` 自动提交/回滚（[tx.go](file:///workspace/database/sql/tx.go#L16-L118)）

### 4.8 cache：缓存接口抽象

**目标**
- 抽象出通用缓存操作能力，便于脚手架在不同缓存实现间切换（Redis/freecache/bigcache 等）。

**接口集合（cache/defs）**
- `Cache`：基础 KV（[cache.go](file:///workspace/cache/defs/cache.go#L8-L20)）
- `HashCache/ListCache/SetCache/SortedSetCache`：类 Redis 结构操作（[cache.go](file:///workspace/cache/defs/cache.go#L22-L84)）
- `Lock`：分布式锁（[cache.go](file:///workspace/cache/defs/cache.go#L86-L94)）
- `StatsProvider`：统计与 key 扫描（[cache.go](file:///workspace/cache/defs/cache.go#L105-L111)）

### 4.9 lang：通用语言层工具

**conv（类型转换）**
- `ToString/ToInt/ToInt64/ToFloat64/ToBool/ToTime` 等（[to.go](file:///workspace/lang/conv/to.go#L10-L224)）

**ptr（指针工具）**
- `ToPtr/Deref/DerefOrDefault` + 常用基础类型 `ToInt64/ToString/...`（[ptr.go](file:///workspace/lang/ptr/ptr.go#L5-L127)）

**slices（泛型 slice 工具）**
- 包含 `Contains/Filter/Map/Reduce/Unique/Chunk/Shuffle/GroupBy/...`（[slices.go](file:///workspace/lang/slices/slices.go#L8-L365)）

**maps（泛型 map 工具）**
- `Keys/Values/Clone/Merge/Filter/MapValues/Equal/...`（[maps.go](file:///workspace/lang/maps/maps.go#L3-L240)）

**ternary（短路/三元工具）**
- `T/Lazy/Coalesce/Default/...`（[ternary.go](file:///workspace/lang/ternary/ternary.go#L3-L103)）

**crypto**
- `MD5/MD5Bytes`（[md5.go](file:///workspace/lang/crypto/md5.go#L8-L20)）

## 5. 运行、测试与 CI

### 5.1 本地运行示例

```bash
go run ./examples/basic_usage.go
```

### 5.2 运行测试

```bash
go test ./...
```

建议与 CI 一致的命令（含竞态与覆盖率）：

```bash
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
```

### 5.3 CI 工作流

CI 定义见：[test.yml](file:///workspace/.github/workflows/test.yml#L1-L61)
- `go test -v -race -coverprofile=coverage.out -covermode=atomic ./...`
- 生成 `coverage.html`
- `golangci-lint`（通过 GitHub Action）

### 5.4 覆盖率现状

覆盖率报告见：[TEST_COVERAGE.md](file:///workspace/TEST_COVERAGE.md#L1-L153)
- 总体覆盖率（文档记录）：55.1%
- 高覆盖模块：`lang/*`（大部分 100%）、`errorx`（100%）、`ctxcache/taskgroup/config`（90%+）
- 待补齐测试：`logs/eventbus/database/sql/errorx/internal/lang/crypto`（文档记录为 0%）

## 6. 扩展与集成指南（面向脚手架/业务项目）

**6.1 替换/注入基础设施**
- 日志：在应用启动阶段调用一次 `logs.Init(...)`，不要在业务链路中重复 init。
- 配置：可以用 `config.Load(path)` 或维护自有实例；如果使用默认实例，尽量集中初始化（`LoadDefault`）。
- MQ：按 [eventbus/impl/README.md](file:///workspace/eventbus/impl/README.md) 在脚手架侧实现 `Producer/ConsumerService`，再通过 `eventbus.SetDefaultSVC(...)` 注入默认消费者服务。
- 数据库：实现 `database/defs.SQLDatabase` 与 `Transaction`，即可复用 `database/sql.TxManager.ExecuteInTx` 与 `Builder/Pager`。
- 缓存：实现 `cache/defs` 中的接口族；建议按需实现（只用 KV 就只实现 `Cache`），避免过度实现成本。

**6.2 常见集成模式**
- “业务包只依赖 defs 接口，基础设施包负责实现并注入”：降低耦合、方便替换实现。
- “在 context 中携带事务/临时缓存”：配合 `TxManager.WithContext` 与 `ctxcache.Init`，把请求级状态收敛在 context 生命周期内。
