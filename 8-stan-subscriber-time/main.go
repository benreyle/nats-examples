package main

import (
	"encoding/json"
	"flag"
	"log"
	event "nats-examples/0-event"
	"runtime"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	// get the queue name
	queue := flag.String("queue", "groupA", "the queue name to subscribe to")
	option := flag.String("option", "sequence", "how to control the delivery")
	sequence := flag.Uint64("sequence", 0, "the desired start sequence position")
	strDate := flag.String("datetime", time.Now().Format("02/01/2006 15:04:05"), "the desired start time position")
	flag.Parse()

	ctrlOpts := make([]stan.SubscriptionOption, 0)

	// set how STAN will control the deliveries
	switch *option {
	case "sequence":
		ctrlOpts = append(ctrlOpts, stan.StartAtSequence(*sequence))

	case "datetime":
		date, err := time.Parse("02/01/2006 15:04:05", *strDate)
		if err != nil {
			log.Fatal("Failed to parse datetime", err)
			return
		}

		date = date.Add(3 * time.Hour)
		ctrlOpts = append(ctrlOpts, stan.StartAtTime(date))
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
	}, ctrlOpts...)

	log.Printf("Listening on 'event' topic with the queue name '%s'\n", *queue)

	runtime.Goexit()
	subs.Close()
}
