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
			log.Println("Failed to unmarshal request event", err)
			return
		}

		resp := event.Event{
			ID:  ev.ID,
			Msg: "Hello, Publisher",
		}

		b, err := json.Marshal(resp)
		if err != nil {
			log.Println("Failed to marshal response event", err)
			return
		}

		err = msg.Respond(b)
		if err != nil {
			log.Println("Failed to respond event", err)
			return
		}

		log.Println("Replied", ev)
	})

	log.Println("Listening on event topic")

	runtime.Goexit()
}
