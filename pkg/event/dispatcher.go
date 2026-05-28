package event

import (
	"sync"
)

// Event đại diện cho gói tin sự kiện Pub/Sub
type Event struct {
	Topic string
	Data  interface{}
}

// EventDispatcher là Broker Pub/Sub nội bộ chạy bằng Goroutine và Channel
type EventDispatcher struct {
	mu          sync.RWMutex
	subscribers map[string][]chan interface{}
}

// NewEventDispatcher khởi tạo một Dispatcher mới
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		subscribers: make(map[string][]chan interface{}),
	}
}

// Subscribe đăng ký kênh channel để nhận sự kiện từ một Topic
func (ed *EventDispatcher) Subscribe(topic string, ch chan interface{}) {
	ed.mu.Lock()
	defer ed.mu.Unlock()
	ed.subscribers[topic] = append(ed.subscribers[topic], ch)
}

// Unsubscribe hủy đăng ký nhận sự kiện
func (ed *EventDispatcher) Unsubscribe(topic string, ch chan interface{}) {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	subs, ok := ed.subscribers[topic]
	if !ok {
		return
	}

	for i, sub := range subs {
		if sub == ch {
			ed.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

// Publish phát tán sự kiện bất đồng bộ qua Goroutines
func (ed *EventDispatcher) Publish(topic string, data interface{}) {
	ed.mu.RLock()
	defer ed.mu.RUnlock()

	subs, ok := ed.subscribers[topic]
	if !ok {
		return
	}

	for _, sub := range subs {
		go func(ch chan interface{}, d interface{}) {
			select {
			case ch <- d:
			default:
				// Bỏ qua nếu channel đầy để tránh tắc nghẽn
			}
		}(sub, data)
	}
}
