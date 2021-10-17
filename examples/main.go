package main

import (
	"fmt"

	"github.com/pscompsci/eventbus"
)

func main() {
	bus := eventbus.New()
	var ch eventbus.DataChannel

	bus.Subscribe("myTopic", ch)

	go func() {
		for {
			data := <-ch
			if data.Topic == "myTopic" {
				fmt.Println(data.Data)
				break
			}
		}
	}()

	bus.Publish("myTopic", "some example message")
}
