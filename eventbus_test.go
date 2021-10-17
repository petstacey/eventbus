package eventbus

import "testing"

func TestDataEvent(t *testing.T) {
	event := DataEvent{
		Topic: "test",
		Data:  "test",
	}
	if event.Data != "test" {
		t.Errorf("failed to create DataChannel. got: %v, want 'test'", event.Data)
	}
}

func TestNewEventBus(t *testing.T) {
	bus := New()
	if len(bus.subscribers) != 0 {
		t.Errorf("new eventbus creation error. got: %d subscribers, want: 0", len(bus.subscribers))
	}
}

func TestSubscribers(t *testing.T) {
	bus := New()

	ch1 := make(chan DataEvent)
	ch2 := make(chan DataEvent)
	ch3 := make(chan DataEvent)

	type test struct {
		topic   string
		channel DataChannel
		want    int
	}

	tests := []test{
		{topic: "test", channel: ch1, want: 1},
		{topic: "test", channel: ch2, want: 2},
		{topic: "test", channel: ch3, want: 3},
		{topic: "test2", channel: ch1, want: 1},
	}

	for _, tc := range tests {
		bus.Subscribe(tc.topic, tc.channel)
		got := len(bus.subscribers[tc.topic])
		if got != tc.want {
			t.Errorf("incorrect number of subscribers to topic. got: %d, want: %d", got, tc.want)
		}
	}
}

func TestSubscribedTopics(t *testing.T) {
	bus := New()

	ch1 := make(chan DataEvent)

	type test struct {
		topic   string
		channel DataChannel
		want    int
	}

	tests := []test{
		{topic: "test", channel: ch1, want: 1},
		{topic: "test2", channel: ch1, want: 2},
		{topic: "tes3", channel: ch1, want: 3},
		{topic: "test4", channel: ch1, want: 4},
	}

	for _, tc := range tests {
		bus.Subscribe(tc.topic, tc.channel)
		if len(bus.subscribers) != tc.want {
			t.Errorf("incorrect number of topics. got: %d, want: %d", len(bus.subscribers), tc.want)
		}
	}
}

func TestUnsubscribe(t *testing.T) {

	ch1 := make(chan DataEvent)
	ch2 := make(chan DataEvent)
	ch3 := make(chan DataEvent)

	type test struct {
		topic       string
		channels    []DataChannel
		unsubscribe DataChannel
		want        int
	}

	tests := []test{
		{topic: "test", channels: []DataChannel{ch1}, unsubscribe: ch1, want: 0},
		{topic: "test", channels: []DataChannel{ch1, ch2}, unsubscribe: ch1, want: 1},
		{topic: "test", channels: []DataChannel{ch1, ch2, ch3}, unsubscribe: ch1, want: 2},
		{topic: "test", channels: []DataChannel{ch1, ch2, ch3}, unsubscribe: ch2, want: 2},
		{topic: "test", channels: []DataChannel{ch1, ch2, ch3}, unsubscribe: ch3, want: 2},
	}

	for _, tc := range tests {
		bus := New()
		for _, ch := range tc.channels {
			bus.Subscribe(tc.topic, ch)
		}
		bus.Unsubscribe(tc.topic, tc.unsubscribe)
		if len(bus.subscribers[tc.topic]) != tc.want {
			t.Errorf("incorrect unsubscribe. got: %d, want: %d", len(bus.subscribers[tc.topic]), tc.want)
		}
	}
}

func TestPublish(t *testing.T) {
	bus := New()

	ch1 := make(chan DataEvent)

	bus.Subscribe("test", ch1)

	go func() {
		data := <-ch1
		if data.Data != "test message" {
			t.Errorf("error")
		}
	}()

	bus.Publish("test2", "wrong topic")
	bus.Publish("test", "test message")
}
