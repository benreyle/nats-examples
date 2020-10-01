package main

import (
	"encoding/json"
	"log"
	event "nats-examples/0-event"
	"runtime"

	"github.com/nats-io/nats.go"
)

func main() {
	// connect to NATS Server
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS", err)
	}

	log.Println("Connected to NATS at " + nats.DefaultURL)

	conn.Subscribe("event", func(msg *nats.Msg) {
		var ev event.Event

		err := json.Unmarshal(msg.Data, &ev)
		if err != nil {
			log.Println("Failed to unmarshal event", err)
			return
		}

		log.Println("Event received", ev)
	})

	log.Println("Listening on event topic")

	runtime.Goexit()
}
