package eventbus

import (
	"reflect"
	"sync"
)

type DataEvent struct {
	Topic string
	Data  interface{}
}

type DataChannel chan DataEvent

type EventBus struct {
	subscribers map[string][]DataChannel
	sync        sync.RWMutex
}

func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]DataChannel),
	}
}

func (bus *EventBus) Subscribe(topic string, ch DataChannel) {
	bus.sync.Lock()
	defer bus.sync.Unlock()

	if t, exists := bus.subscribers[topic]; exists {
		bus.subscribers[topic] = append(t, ch)
	} else {
		bus.subscribers[topic] = append([]DataChannel{}, ch)
	}
}

func (bus *EventBus) Unsubscribe(topic string, ch DataChannel) {
	bus.sync.Lock()
	defer bus.sync.Unlock()

	if _, exists := bus.subscribers[topic]; !exists {
		return
	}

	if len(bus.subscribers[topic]) <= 1 {
		bus.subscribers[topic] = []DataChannel{}
		return
	}

	for i, subscriber := range bus.subscribers[topic] {
		if reflect.DeepEqual(subscriber, ch) {
			if i == 0 {
				bus.subscribers[topic] = bus.subscribers[topic][1:]
				return
			}
			if i == len(bus.subscribers[topic])-1 {
				bus.subscribers[topic] = bus.subscribers[topic][:len(bus.subscribers[topic])-1]
				return
			}
			bus.subscribers[topic] = append(bus.subscribers[topic][:i], bus.subscribers[topic][i+1:]...)
		}
	}
}

func (bus *EventBus) Publish(topic string, data interface{}) {
	bus.sync.RLock()
	if s, exists := bus.subscribers[topic]; exists {
		subscribers := append([]DataChannel{}, s...)
		go func(data DataEvent, channels []DataChannel) {
			for _, ch := range channels {
				ch <- data
			}
		}(DataEvent{Topic: topic, Data: data}, subscribers)
	}
}
