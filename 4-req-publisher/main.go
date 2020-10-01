package main

import (
	"encoding/json"
	"log"
	event "nats-examples/0-event"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// connect to NATS Server
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS", err)
	}

	log.Println("Connected to NATS at " + nats.DefaultURL)

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

		log.Println("Requesting event", string(b))

		msg, err := conn.Request("event", b, 5*time.Second)
		if err != nil {
			log.Println("Failed to request event", err)
			return
		}

		log.Println("Response received", string(msg.Data))
		i++
	}
}
