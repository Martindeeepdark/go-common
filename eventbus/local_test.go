package eventbus

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEventBus(t *testing.T) {
	bus := NewEventBus()
	assert.NotNil(t, bus)
}

func TestSubscribe(t *testing.T) {
	t.Run("subscribe to topic", func(t *testing.T) {
		bus := NewEventBus()
		called := false

		err := bus.Subscribe("topic1", func(msg string) {
			called = true
			assert.Equal(t, "test", msg)
		})
		assert.NoError(t, err)

		bus.Publish("topic1", "test")
		assert.True(t, called)
	})

	t.Run("subscribe with non-function returns error", func(t *testing.T) {
		bus := NewEventBus()
		err := bus.Subscribe("topic1", "not a function")
		assert.Error(t, err)
	})

	t.Run("multiple subscribers to same topic", func(t *testing.T) {
		bus := NewEventBus()
		count := 0
		var mu sync.Mutex

		handler1 := func(msg string) {
			mu.Lock()
			count++
			mu.Unlock()
		}
		handler2 := func(msg string) {
			mu.Lock()
			count++
			mu.Unlock()
		}

		bus.Subscribe("topic1", handler1)
		bus.Subscribe("topic1", handler2)

		bus.Publish("topic1", "test")
		time.Sleep(10 * time.Millisecond)

		assert.Equal(t, 2, count)
	})
}

func TestSubscribeAsync(t *testing.T) {
	t.Run("async subscription", func(t *testing.T) {
		bus := NewEventBus()
		called := false
		var mu sync.Mutex

		err := bus.SubscribeAsync("topic1", func(msg string) {
			mu.Lock()
			called = true
			mu.Unlock()
		}, false)
		assert.NoError(t, err)

		bus.Publish("topic1", "test")
		bus.WaitAsync()

		mu.Lock()
		assert.True(t, called)
		mu.Unlock()
	})

	t.Run("transactional async subscription", func(t *testing.T) {
		bus := NewEventBus()
		count := 0
		var mu sync.Mutex

		for i := 0; i < 3; i++ {
			err := bus.SubscribeAsync("topic1", func(msg int) {
				mu.Lock()
				count++
				mu.Unlock()
			}, true) // transactional = serial execution
			assert.NoError(t, err)
		}

		for i := 0; i < 3; i++ {
			bus.Publish("topic1", i)
		}

		bus.WaitAsync()
		assert.Equal(t, 9, count) // 3 publishes * 3 handlers
	})
}

func TestSubscribeOnce(t *testing.T) {
	t.Run("subscribe once - handler called only once", func(t *testing.T) {
		bus := NewEventBus()
		count := 0
		var mu sync.Mutex

		err := bus.SubscribeOnce("topic1", func(msg string) {
			mu.Lock()
			count++
			mu.Unlock()
		})
		assert.NoError(t, err)

		bus.Publish("topic1", "test1")
		bus.Publish("topic1", "test2")
		bus.Publish("topic1", "test3")

		mu.Lock()
		assert.Equal(t, 1, count)
		mu.Unlock()
	})

	t.Run("subscribe once - handler removed after execution", func(t *testing.T) {
		bus := NewEventBus()

		handler := func(msg string) {}
		bus.SubscribeOnce("topic1", handler)

		bus.Publish("topic1", "test")
		assert.False(t, bus.HasCallback("topic1"))
	})
}

func TestSubscribeOnceAsync(t *testing.T) {
	t.Run("subscribe once async", func(t *testing.T) {
		bus := NewEventBus()
		count := 0
		var mu sync.Mutex

		err := bus.SubscribeOnceAsync("topic1", func(msg string) {
			mu.Lock()
			count++
			mu.Unlock()
		})
		assert.NoError(t, err)

		bus.Publish("topic1", "test1")
		bus.Publish("topic1", "test2")
		bus.WaitAsync()

		mu.Lock()
		assert.Equal(t, 1, count)
		mu.Unlock()
	})
}

func TestUnsubscribe(t *testing.T) {
	t.Run("unsubscribe handler", func(t *testing.T) {
		bus := NewEventBus()
		called := false

		handler := func(msg string) {
			called = true
		}

		bus.Subscribe("topic1", handler)
		assert.True(t, bus.HasCallback("topic1"))

		err := bus.Unsubscribe("topic1", handler)
		assert.NoError(t, err)
		assert.False(t, bus.HasCallback("topic1"))

		bus.Publish("topic1", "test")
		assert.False(t, called)
	})

	t.Run("unsubscribe from non-existent topic", func(t *testing.T) {
		bus := NewEventBus()
		handler := func(msg string) {}

		err := bus.Unsubscribe("non-existent", handler)
		assert.Error(t, err)
	})
}

func TestHasCallback(t *testing.T) {
	t.Run("has callback returns true when subscribed", func(t *testing.T) {
		bus := NewEventBus()
		handler := func(msg string) {}

		bus.Subscribe("topic1", handler)
		assert.True(t, bus.HasCallback("topic1"))
	})

	t.Run("has callback returns false when no subscribers", func(t *testing.T) {
		bus := NewEventBus()
		assert.False(t, bus.HasCallback("topic1"))
	})

	t.Run("has callback returns false after unsubscribe", func(t *testing.T) {
		bus := NewEventBus()
		handler := func(msg string) {}

		bus.Subscribe("topic1", handler)
		bus.Unsubscribe("topic1", handler)

		assert.False(t, bus.HasCallback("topic1"))
	})
}

func TestPublish(t *testing.T) {
	t.Run("publish with no arguments", func(t *testing.T) {
		bus := NewEventBus()
		called := false

		bus.Subscribe("topic1", func() {
			called = true
		})

		bus.Publish("topic1")
		assert.True(t, called)
	})

	t.Run("publish with multiple arguments", func(t *testing.T) {
		bus := NewEventBus()

		bus.Subscribe("topic1", func(s string, i int, b bool) {
			assert.Equal(t, "test", s)
			assert.Equal(t, 42, i)
			assert.Equal(t, true, b)
		})

		bus.Publish("topic1", "test", 42, true)
	})

	t.Run("publish to non-existent topic", func(t *testing.T) {
		bus := NewEventBus()
		assert.NotPanics(t, func() {
			bus.Publish("non-existent")
		})
	})

	t.Run("publish with nil argument", func(t *testing.T) {
		bus := NewEventBus()
		called := false

		bus.Subscribe("topic1", func(s string) {
			called = true
			assert.Equal(t, "", s)
		})

		bus.Publish("topic1", nil)
		assert.True(t, called)
	})
}

func TestWaitAsync(t *testing.T) {
	t.Run("wait for async handlers", func(t *testing.T) {
		bus := NewEventBus()
		count := 0
		var mu sync.Mutex

		bus.SubscribeAsync("topic1", func(msg string) {
			time.Sleep(50 * time.Millisecond)
			mu.Lock()
			count++
			mu.Unlock()
		}, false)

		bus.Publish("topic1", "test")

		// Wait should block until async handler completes
		bus.WaitAsync()

		mu.Lock()
		assert.Equal(t, 1, count)
		mu.Unlock()
	})
}

func TestConcurrentPublish(t *testing.T) {
	t.Run("concurrent publishes", func(t *testing.T) {
		bus := NewEventBus()
		count := 0
		var mu sync.Mutex

		handler := func(msg int) {
			mu.Lock()
			count++
			mu.Unlock()
		}

		bus.Subscribe("topic1", handler)

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				bus.Publish("topic1", idx)
			}(i)
		}

		wg.Wait()
		assert.Equal(t, 10, count)
	})
}

func TestMultipleTopics(t *testing.T) {
	t.Run("different topics are independent", func(t *testing.T) {
		bus := NewEventBus()
		topic1Called := false
		topic2Called := false

		bus.Subscribe("topic1", func() {
			topic1Called = true
		})

		bus.Subscribe("topic2", func() {
			topic2Called = true
		})

		bus.Publish("topic1")
		assert.True(t, topic1Called)
		assert.False(t, topic2Called)

		bus.Publish("topic2")
		assert.True(t, topic2Called)
	})
}

func TestHandlerRemovalDuringPublish(t *testing.T) {
	t.Run("SubscribeOnce removes handler during publish", func(t *testing.T) {
		bus := NewEventBus()
		count := 0
		var mu sync.Mutex

		bus.SubscribeOnce("topic1", func() {
			mu.Lock()
			count++
			mu.Unlock()
		})

		bus.Subscribe("topic1", func() {
			mu.Lock()
			count++
			mu.Unlock()
		})

		bus.Publish("topic1")
		bus.Publish("topic1")

		mu.Lock()
		assert.Equal(t, 3, count) // Once handler + regular handler * 2 publishes
		mu.Unlock()
	})
}

func TestSubscribeAsyncNonTransactional(t *testing.T) {
	t.Run("concurrent async handlers", func(t *testing.T) {
		bus := NewEventBus()
		var mu sync.Mutex
		executed := make([]int, 0)

		for i := 0; i < 5; i++ {
			idx := i
			bus.SubscribeAsync("topic1", func(msg int) {
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				executed = append(executed, idx)
				mu.Unlock()
			}, false) // non-transactional = concurrent
		}

		bus.Publish("topic1", 0)
		bus.WaitAsync()

		mu.Lock()
		assert.Len(t, executed, 5)
		mu.Unlock()
	})
}

func TestEmptyTopic(t *testing.T) {
	t.Run("publish to empty topic", func(t *testing.T) {
		bus := NewEventBus()
		assert.NotPanics(t, func() {
			bus.Publish("")
		})
	})

	t.Run("subscribe to empty topic", func(t *testing.T) {
		bus := NewEventBus()
		called := false

		err := bus.Subscribe("", func() {
			called = true
		})
		assert.NoError(t, err)

		bus.Publish("")
		assert.True(t, called)
	})
}

func TestNilFunction(t *testing.T) {
	t.Run("subscribe with nil function", func(t *testing.T) {
		bus := NewEventBus()
		// Nil function causes panic, so we test with defer/recover
		assert.Panics(t, func() {
			bus.Subscribe("topic1", nil)
		})
	})
}
