package helpers

import (
	"log"

	"github.com/nats-io/stan.go"
)

func ConnectToSTAN(client string) stan.Conn {
	// connect to NATS Server
	opts := []stan.Option{
		stan.NatsURL(stan.DefaultNatsURL),
	}

	conn, err := stan.Connect("test-cluster", client, opts...)
	if err != nil {
		log.Fatal("Failed to connect to STAN", err)
	}

	log.Println("Connected to STAN at " + stan.DefaultNatsURL)

	return conn
}
