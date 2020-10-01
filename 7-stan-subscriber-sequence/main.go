package main

import (
	"encoding/json"
	"flag"
	"log"
	event "nats-examples/0-event"
	"runtime"

	"github.com/nats-io/stan.go"
)

func main() {
	// get the queue name
	queue := flag.String("queue", "groupA", "the queue name to subscribe to")
	option := flag.String("option", "sequence", "how to control the delivery")
	sequence := flag.Uint64("sequence", 0, "the desired start sequence")
	flag.Parse()

	ctlrOpts := make([]stan.SubscriptionOption, 0)

	// set how STAN will control the deliveries
	switch *option {
	case "sequence":
		ctlrOpts = append(ctlrOpts, stan.StartAtSequence(*sequence))
	}

	// connect to NATS Server
	opts := []stan.Option{
		stan.NatsURL(stan.DefaultNatsURL),
	}

	conn, err := stan.Connect("test-cluster", "subscriber", opts...)
	if err != nil {
		log.Fatal("Failed to connect to STAN", err)
	}

	log.Println("Connected to STAN at " + stan.DefaultNatsURL)

	subs, _ := conn.QueueSubscribe("event", *queue, func(msg *stan.Msg) {
		var ev event.Event

		err := json.Unmarshal(msg.Data, &ev)
		if err != nil {
			log.Println("Failed to unmarshal event", err)
			return
		}

		log.Println("Event received", ev)
	}, ctlrOpts...)

	log.Printf("Listening on 'event' topic with the queue name '%s'\n", *queue)

	runtime.Goexit()
	subs.Close()
}
