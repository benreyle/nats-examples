package main

import (
	"encoding/json"
	"flag"
	"log"
	event "nats-examples/0-event"
	"runtime"

	"github.com/nats-io/nats.go"
)

func main() {
	// get the queue name
	queue := flag.String("queue", "groupA", "the queue name to subscribe to")
	flag.Parse()

	// connect to NATS Server
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS", err)
	}

	log.Println("Connected to NATS at " + nats.DefaultURL)

	// subscribe
	conn.QueueSubscribe("event", *queue, func(msg *nats.Msg) {
		var ev event.Event

		err := json.Unmarshal(msg.Data, &ev)
		if err != nil {
			log.Println("Failed to unmarshal event", err)
			return
		}

		log.Println("Event received", ev)
	})

	log.Printf("Listening on 'event' topic with the queue name '%s'\n", *queue)

	runtime.Goexit()
}
