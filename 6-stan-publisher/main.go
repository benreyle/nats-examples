package main

import (
	"encoding/json"
	"log"
	event "nats-examples/0-event"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	opts := []stan.Option{
		stan.NatsURL(stan.DefaultNatsURL),
	}

	conn, err := stan.Connect("test-cluster", "publisher", opts...)
	if err != nil {
		log.Fatal("Failed to connect to STAN", err)
	}

	log.Println("Connected to STAN at " + stan.DefaultNatsURL)

	i := 0

	// repeat every 2 seconds
	for range time.Tick(2 * time.Second) {
		ev := event.Event{
			ID:  i,
			Msg: "Hello, Subscriber!",
		}

		// marshal the event to JSON
		b, err := json.Marshal(ev)
		if err != nil {
			log.Println("Failed to marshal data", err)
			return
		}

		err = conn.Publish("event", b)
		if err != nil {
			log.Println("Failed to publish event", err)
			return
		}

		log.Println("Event published", string(b))
		i++
	}
}
