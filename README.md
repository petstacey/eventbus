# eventbus

A simple pub/sub implementation in Go. The EventBus uses channels to subscribe to topics and publishes back to the channels when data is sent for a topic.

## Installaton

```bash
go get github.com/pscompsci/eventbus
```

## Usage

```Go
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

# prints "some example message"
```

## Contributing
Pull requests to add functionality are welcome. For major changes, please open an issue first, to discuss the changes in advance.

Please make sure to update tests as appropriate

## License
[MIT](https://choosealicense.com/licenses/mit/)

