package EventBus

import (
	"fmt"
	"reflect"
	"sync"
)

// BusSubscriber 定义订阅相关的总线行为
type BusSubscriber interface {
	Subscribe(topic string, fn interface{}) error
	SubscribeAsync(topic string, fn interface{}, transactional bool) error
	SubscribeOnce(topic string, fn interface{}) error
	SubscribeOnceAsync(topic string, fn interface{}) error
	Unsubscribe(topic string, handler interface{}) error
}

// BusPublisher 定义发布相关的总线行为
type BusPublisher interface {
	Publish(topic string, args ...interface{})
}

// BusController 定义总线控制行为
type BusController interface {
	HasCallback(topic string) bool
	WaitAsync()
}

// Bus 包含全局（订阅、发布、controller）总线行为
type Bus interface {
	BusController
	BusSubscriber
	BusPublisher
}

// EventBus 事件总线 - 用于存储handler和回调函数的容器
type EventBus struct {
	handlers map[string][]*eventHandler
	lock     sync.Mutex // 映射的锁
	wg       sync.WaitGroup
}

type eventHandler struct {
	callBack      reflect.Value
	flagOnce      bool
	async         bool
	transactional bool
	sync.Mutex    // 事件处理器的锁 - 用于串行运行异步回调函数
}

// NewEventBus 返回一个新的EventBus实例
func NewEventBus() Bus {
	b := &EventBus{
		make(map[string][]*eventHandler),
		sync.Mutex{},
		sync.WaitGroup{},
	}
	return Bus(b)
}

// doSubscribe 处理订阅逻辑，被公共订阅函数使用
func (bus *EventBus) doSubscribe(topic string, fn interface{}, handler *eventHandler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s 不是 reflect.Func 类型", reflect.TypeOf(fn).Kind())
	}
	bus.handlers[topic] = append(bus.handlers[topic], handler)
	return nil
}

// Subscribe 订阅一个主题
// 如果 `fn` 不是函数类型则返回错误
func (bus *EventBus) Subscribe(topic string, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn), false, false, false, sync.Mutex{},
	})
}

// SubscribeAsync 异步订阅一个主题
// Transactional 决定主题的后续回调是串行运行（true）还是并发运行（false）
// 如果 `fn` 不是函数类型则返回错误
func (bus *EventBus) SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn), false, true, transactional, sync.Mutex{},
	})
}

// SubscribeOnce 订阅一个主题，处理器将在执行后被移除
// 如果 `fn` 不是函数类型则返回错误
func (bus *EventBus) SubscribeOnce(topic string, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn), true, false, false, sync.Mutex{},
	})
}

// SubscribeOnceAsync 异步订阅一个主题，处理器将在执行后被移除
// 如果 `fn` 不是函数类型则返回错误
func (bus *EventBus) SubscribeOnceAsync(topic string, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		reflect.ValueOf(fn), true, true, false, sync.Mutex{},
	})
}

// HasCallback 如果主题存在任何订阅的回调函数则返回true
func (bus *EventBus) HasCallback(topic string) bool {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	_, ok := bus.handlers[topic]
	if ok {
		return len(bus.handlers[topic]) > 0
	}
	return false
}

// Unsubscribe 移除主题的回调函数
// 如果主题没有任何订阅的回调函数则返回错误
func (bus *EventBus) Unsubscribe(topic string, handler interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if _, ok := bus.handlers[topic]; ok && len(bus.handlers[topic]) > 0 {
		bus.removeHandler(topic, bus.findHandlerIdx(topic, reflect.ValueOf(handler)))
		return nil
	}
	return fmt.Errorf("主题 %s 不存在", topic)
}

// Publish 执行为topic定义的回调函数，任何额外的参数将传递给回调函数
func (bus *EventBus) Publish(topic string, args ...interface{}) {
	bus.lock.Lock() // 如果未找到handler或在setUpPublish后总是会解锁
	defer bus.lock.Unlock()
	if handlers, ok := bus.handlers[topic]; ok && 0 < len(handlers) {
		// 处理器切片可能在迭代过程中被removeHandler和Unsubscribe修改，
		// 因此创建一个副本并迭代该副本
		copyHandlers := make([]*eventHandler, len(handlers))
		copy(copyHandlers, handlers)
		for i, handler := range copyHandlers {
			if handler.flagOnce {
				bus.removeHandler(topic, i)
			}
			if !handler.async {
				bus.doPublish(handler, topic, args...)
			} else {
				bus.wg.Add(1)
				if handler.transactional {
					bus.lock.Unlock()
					handler.Lock()
					bus.lock.Lock()
				}
				go bus.doPublishAsync(handler, topic, args...)
			}
		}
	}
}

func (bus *EventBus) doPublish(handler *eventHandler, topic string, args ...interface{}) {
	passedArguments := bus.setUpPublish(handler, args...)
	handler.callBack.Call(passedArguments)
}

func (bus *EventBus) doPublishAsync(handler *eventHandler, topic string, args ...interface{}) {
	defer bus.wg.Done()
	if handler.transactional {
		defer handler.Unlock()
	}
	bus.doPublish(handler, topic, args...)
}

func (bus *EventBus) removeHandler(topic string, idx int) {
	if _, ok := bus.handlers[topic]; !ok {
		return
	}
	l := len(bus.handlers[topic])

	if !(0 <= idx && idx < l) {
		return
	}

	copy(bus.handlers[topic][idx:], bus.handlers[topic][idx+1:])
	bus.handlers[topic][l-1] = nil // 或者 T 类型的零值
	bus.handlers[topic] = bus.handlers[topic][:l-1]
}

func (bus *EventBus) findHandlerIdx(topic string, callback reflect.Value) int {
	if _, ok := bus.handlers[topic]; ok {
		for idx, handler := range bus.handlers[topic] {
			if handler.callBack.Type() == callback.Type() &&
				handler.callBack.Pointer() == callback.Pointer() {
				return idx
			}
		}
	}
	return -1
}

func (bus *EventBus) setUpPublish(callback *eventHandler, args ...interface{}) []reflect.Value {
	funcType := callback.callBack.Type()
	passedArguments := make([]reflect.Value, len(args))
	for i, v := range args {
		if v == nil {
			passedArguments[i] = reflect.New(funcType.In(i)).Elem()
		} else {
			passedArguments[i] = reflect.ValueOf(v)
		}
	}

	return passedArguments
}

// WaitAsync 等待所有异步回调函数完成
func (bus *EventBus) WaitAsync() {
	bus.wg.Wait()
}
