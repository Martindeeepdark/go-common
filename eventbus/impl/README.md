# Eventbus 消息队列实现

## 当前状态

**接口已定义，具体实现待添加。**

已完成的：
- ✅ `eventbus.Producer` - 生产者接口
- ✅ `eventbus.ConsumerService` - 消费者服务接口
- ✅ `eventbus.ConsumerHandler` - 消息处理器接口
- ✅ `impl.NewConsumerService()` - 工厂方法
- ✅ `impl.NewProducer()` - 工厂方法
- ✅ 环境变量 `COMMON_MQ_TYPE` 控制队列类型

## 实现方式

### 方案 1: 在脚手架中实现（推荐）

在脚手架项目（如 scaffold-gin）中根据实际使用的消息队列添加实现：

```
scaffold-gin/
├── eventbus/
│   ├── nsq/          # NSQ 实现
│   │   ├── producer.go
│   │   └── consumer.go
│   └── kafka/        # Kafka 实现
│       ├── producer.go
│       └── consumer.go
```

```go
// scaffold-gin/eventbus/nsq/producer.go
package nsq

import (
	"context"
	"github.com/nsqio/go-nsq"
	"common/eventbus"
)

type producerImpl struct {
	producer *nsq.Producer
	topic    string
}

func NewProducer(nameServer, topic, group string) (eventbus.Producer, error) {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(nameServer, config)
	if err != nil {
		return nil, err
	}

	return &producerImpl{producer: p, topic: topic}, nil
}

func (p *producerImpl) Send(ctx context.Context, body []byte, opts ...eventbus.SendOpt) error {
	return p.producer.Publish(p.topic, body)
}

func (p *producerImpl) BatchSend(ctx context.Context, bodyArr [][]byte, opts ...eventbus.SendOpt) error {
	return p.producer.MultiPublish(p.topic, bodyArr)
}
```

### 方案 2: 在 common 包中添加依赖

如果需要在 common 包中直接使用，添加对应依赖：

```bash
# NSQ（最轻量）
go get github.com/nsqio/go-nsq

# Kafka
go get github.com/segmentio/kafka-go

# RocketMQ
go get github.com/apache/rocketmq-client-go

# Pulsar
go get github.com/apache/pulsar-client-go

# NATS
go get github.com/nats-io/nats.go
```

然后在 `eventbus/impl/` 下创建对应的实现目录。

## 参考实现

### NSQ 完整示例

```go
// eventbus/impl/nsq/producer.go
package nsq

import (
	"context"
	"fmt"
	"github.com/nsqio/go-nsq"
	"common/eventbus"
)

type producerImpl struct {
	producer *nsq.Producer
	topic    string
}

func NewProducer(nameServer, topic, group string) (eventbus.Producer, error) {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(nameServer, config)
	if err != nil {
		return nil, fmt.Errorf("create nsq producer failed: %w", err)
	}
	return &producerImpl{producer: p, topic: topic}, nil
}

func (p *producerImpl) Send(ctx context.Context, body []byte, opts ...eventbus.SendOpt) error {
	err := p.producer.Publish(p.topic, body)
	if err != nil {
		return fmt.Errorf("nsq publish failed: %w", err)
	}
	return nil
}

func (p *producerImpl) BatchSend(ctx context.Context, bodyArr [][]byte, opts ...eventbus.SendOpt) error {
	err := p.producer.MultiPublish(p.topic, bodyArr)
	if err != nil {
		return fmt.Errorf("nsq multi-publish failed: %w", err)
	}
	return nil
}
```

```go
// eventbus/impl/nsq/consumer.go
package nsq

import (
	"context"
	"fmt"
	"github.com/nsqio/go-nsq"
	"common/eventbus"
)

func RegisterConsumer(nameServer, topic, group string, handler eventbus.ConsumerHandler, opts ...eventbus.ConsumerOpt) error {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, group, config)
	if err != nil {
		return fmt.Errorf("create nsq consumer failed: %w", err)
	}

	consumer.AddHandler(&messageHandler{
		topic:   topic,
		group:   group,
		handler: handler,
	})

	return consumer.ConnectToNSQD(nameServer)
}

type messageHandler struct {
	topic   string
	group   string
	handler eventbus.ConsumerHandler
}

func (h *messageHandler) HandleMessage(m *nsq.Message) error {
	msg := &eventbus.Message{
		Topic: h.topic,
		Group: h.group,
		Body:  m.Body,
	}
	return h.handler.HandleMessage(context.Background(), msg)
}
```

## 使用示例

### 生产者

```go
import (
    "common/eventbus"
    "common/eventbus/impl"
    "context"
)

producer, err := impl.NewProducer("localhost:4150", "my-topic", "my-group")
if err != nil {
    panic(err)
}

// 发送单条消息
err = producer.Send(context.Background(), []byte("hello"))

// 批量发送
messages := [][]byte{
    []byte("msg1"),
    []byte("msg2"),
}
err = producer.BatchSend(context.Background(), messages, eventbus.WithShardingKey("key123"))
```

### 消费者

```go
import (
    "common/eventbus"
    "common/eventbus/impl"
)

type MyHandler struct{}

func (h *MyHandler) HandleMessage(ctx context.Context, msg *eventbus.Message) error {
    fmt.Printf("Received: Topic=%s, Body=%s\n", msg.Topic, string(msg.Body))
    return nil
}

func main() {
    consumerSvc := impl.NewConsumerService()
    eventbus.SetDefaultSVC(consumerSvc)

    handler := &MyHandler{}
    err := consumerSvc.RegisterConsumer("localhost:4150", "my-topic", "my-channel", handler)
    if err != nil {
        panic(err)
    }

    // 保持程序运行
    select {}
}
```

## 环境变量配置

```bash
# 设置使用的消息队列类型
export COMMON_MQ_TYPE=nsq  # 支持: nsq, kafka, rmq, pulsar, nats

# NSQ 配置
export NSQ_HOST=localhost:4150

# Kafka 配置
export KAFKA_BROKERS=localhost:9092

# RocketMQ 配置
export RMQ_NAMESERVER=localhost:9876

# Pulsar 配置
export PULSAR_SERVICE_URL=pulsar://localhost:6650

# NATS 配置
export NATS_URL=nats://localhost:4222
```

## 注意事项

1. **依赖管理**: 每种消息队列都有对应的客户端库依赖，按需添加
2. **编译标签**: 可以使用 build tag 控制不同 MQ 的编译
3. **连接管理**: 生产环境需要处理连接池、重连、优雅关闭等
4. **错误处理**: 建议实现重试机制和错误日志
5. **性能优化**: 批量发送、异步发送等高级特性

## 推荐方案

**对于脚手架项目**：
- 默认只添加 NSQ（最轻量，单机即可运行）
- 其他 MQ 根据实际需求添加
- 在脚手架文档中说明如何切换 MQ

**对于 common 包**：
- 保持接口定义
- 不强制添加任何 MQ 依赖
- 由脚手架选择具体的 MQ 实现
