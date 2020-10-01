package main

import (
	"encoding/json"
	"log"
	"nats-examples/10-services/helpers"

	"runtime"

	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"gopkg.in/resty.v1"
)

var conn stan.Conn

func main() {
	conn = helpers.ConnectToSTAN("notify")

	createSubs, err := subscribeCreatePosition()
	if err != nil {
		log.Fatal("Failed to subscribe on CreatePosition subject", err)
	}

	runtime.Goexit()
	createSubs.Close()
}

func subscribeCreatePosition() (stan.Subscription, error) {
	return conn.QueueSubscribe("create-position", "notifications", func(msg *stan.Msg) {
		var id uuid.UUID

		err := json.Unmarshal(msg.Data, &id)
		if err != nil {
			log.Println("Failed to unmarshal position", err)
			return
		}

		// get position in api service
		resp, err := resty.R().SetQueryParam("id", id.String()).Get("http://localhost:8000/position")
		if err != nil {
			log.Println("Failed to get position", err, id)
			return
		}

		var position helpers.Position

		err = json.Unmarshal(resp.Body(), &position)
		if err != nil {
			log.Println("Failed to unmarshal position", err, id)
			return
		}

		// check if deleted
		if position.Deleted {
			log.Println("Position deleted", id)
			msg.Ack()
			return
		}

		msg.Ack()

		log.Printf("Send email to %s with id %s", position.Email, position.ID.String())
	}, stan.DurableName("durable-notifications"), stan.SetManualAckMode())
}
