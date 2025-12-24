package eventbus

import (
	"fmt"
	"reflect"
	"sync"
)

// BusSubscriber defines subscription-related bus behavior
type BusSubscriber interface {
	Subscribe(topic string, fn interface{}) error
	SubscribeAsync(topic string, fn interface{}, transactional bool) error
	SubscribeOnce(topic string, fn interface{}) error
	SubscribeOnceAsync(topic string, fn interface{}) error
	Unsubscribe(topic string, handler interface{}) error
}

// BusPublisher defines publish-related bus behavior
type BusPublisher interface {
	Publish(topic string, args ...interface{})
}

// BusController defines bus control behavior
type BusController interface {
	HasCallback(topic string) bool
	WaitAsync()
}

// Bus combines all bus behaviors
type Bus interface {
	BusController
	BusSubscriber
	BusPublisher
}

// EventBus is the event bus for storing handlers and callback functions
type EventBus struct {
	handlers map[string][]*eventHandler
	lock     sync.Mutex
	wg       sync.WaitGroup
}

type eventHandler struct {
	callBack      reflect.Value
	flagOnce      bool
	async         bool
	transactional bool
	sync.Mutex
}

// NewEventBus creates and returns a new EventBus instance
func NewEventBus() Bus {
	b := &EventBus{
		handlers: make(map[string][]*eventHandler),
		lock:     sync.Mutex{},
		wg:       sync.WaitGroup{},
	}
	return Bus(b)
}

// doSubscribe handles subscription logic, used by public subscription functions
func (bus *EventBus) doSubscribe(topic string, fn interface{}, handler *eventHandler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not reflect.Func type", reflect.TypeOf(fn).Kind())
	}

	bus.handlers[topic] = append(bus.handlers[topic], handler)
	return nil
}

// Subscribe subscribes to a topic synchronously
// Returns error if fn is not a function type
func (bus *EventBus) Subscribe(topic string, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		callBack:      reflect.ValueOf(fn),
		flagOnce:      false,
		async:         false,
		transactional: false,
		Mutex:         sync.Mutex{},
	})
}

// SubscribeAsync subscribes to a topic asynchronously
// Transactional determines whether subsequent callbacks for the topic run serially (true) or concurrently (false)
// Returns error if fn is not a function type
func (bus *EventBus) SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		callBack:      reflect.ValueOf(fn),
		flagOnce:      false,
		async:         true,
		transactional: transactional,
		Mutex:         sync.Mutex{},
	})
}

// SubscribeOnce subscribes to a topic, the handler will be removed after execution
// Returns error if fn is not a function type
func (bus *EventBus) SubscribeOnce(topic string, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		callBack:      reflect.ValueOf(fn),
		flagOnce:      true,
		async:         false,
		transactional: false,
		Mutex:         sync.Mutex{},
	})
}

// SubscribeOnceAsync subscribes to a topic asynchronously, the handler will be removed after execution
// Returns error if fn is not a function type
func (bus *EventBus) SubscribeOnceAsync(topic string, fn interface{}) error {
	return bus.doSubscribe(topic, fn, &eventHandler{
		callBack:      reflect.ValueOf(fn),
		flagOnce:      true,
		async:         true,
		transactional: false,
		Mutex:         sync.Mutex{},
	})
}

// HasCallback returns true if the topic has any subscribed callback functions
func (bus *EventBus) HasCallback(topic string) bool {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if _, ok := bus.handlers[topic]; ok {
		return len(bus.handlers[topic]) > 0
	}
	return false
}

// Unsubscribe removes the callback function for a topic
// Returns error if the topic has no subscribed callback functions
func (bus *EventBus) Unsubscribe(topic string, handler interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if _, ok := bus.handlers[topic]; ok && len(bus.handlers[topic]) > 0 {
		bus.removeHandler(topic, bus.findHandlerIdx(topic, reflect.ValueOf(handler)))
		return nil
	}
	return fmt.Errorf("topic %s does not exist", topic)
}

// Publish executes the callback functions defined for the topic,
// any additional arguments will be passed to the callback functions
func (bus *EventBus) Publish(topic string, args ...interface{}) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if handlers, ok := bus.handlers[topic]; ok && len(handlers) > 0 {
		// Create a copy of handlers slice as it may be modified during iteration
		// by removeHandler and Unsubscribe
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
	bus.handlers[topic][l-1] = nil
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

// WaitAsync waits for all async callback functions to complete
func (bus *EventBus) WaitAsync() {
	bus.wg.Wait()
}
