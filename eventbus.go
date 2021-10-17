/*
Package eventbus is a simple implementation of a pub/sub pattern.

The core structue is the EventBus that holds a map of topics and subscribed channels:

	EventBus:
		subscribers map[string][]DataChannel
		sync        sync.RWMutex

The structure of data sent on the event bus is:

	DataEvent:
		Topic string
		Data  interface{}

Topics are subscribed and publish to, on DataChannels that are a
chan of the DataEvent type.
*/
package eventbus

import (
	"reflect"
	"sync"
)

// DataEvent is a container for data to be published on the EventBus
type DataEvent struct {
	Topic string
	Data  interface{}
}

// DataChannel is a DataEvent channel. This is the type used to subscribe and publish on
type DataChannel chan DataEvent

// EventBus is a structure to hold the list of topics and subscribed channels
type EventBus struct {
	subscribers map[string][]DataChannel
	sync        sync.RWMutex
}

// New returns a reference to a new EventBus
func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]DataChannel),
	}
}

// Subscribe allows systems to subscribe to topics using channels, in order to receive a DataEvent when published
func (bus *EventBus) Subscribe(topic string, ch DataChannel) {
	bus.sync.Lock()
	defer bus.sync.Unlock()

	if t, exists := bus.subscribers[topic]; exists {
		bus.subscribers[topic] = append(t, ch)
	} else {
		bus.subscribers[topic] = append([]DataChannel{}, ch)
	}
}

// Unsubscribe allows channels to be removed from topic subscriptions
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

// Publish allows a DataEvent to be sent to all subscribers of a topic
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
