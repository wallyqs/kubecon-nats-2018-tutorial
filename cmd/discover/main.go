package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	log.SetFlags(0)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	t := time.AfterFunc(1*time.Second, func() {
		cancel()
	})

	inbox := nats.NewInbox()
	replies := make([]string, 0)
	sub, err := nc.SubscribeSync(inbox)
	if err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	nc.PublishRequest("_NYFT.discovery", inbox, []byte(""))
	for {
		msg, err := sub.NextMsgWithContext(ctx)
		if err != nil {
			break
		}
		id := string(msg.Data)
		replies = append(replies, id)

		// Extend deadline on each successful response.
		t.Reset(150 * time.Millisecond)
	}
	log.Printf("Found %d components in %.3fms", len(replies), time.Since(startTime).Seconds())

	// Checking status of available components
	for _, componentID := range replies {
		statusSubject := fmt.Sprintf("_NYFT.%s.statz", componentID)
		resp, err := nc.Request(statusSubject, []byte(""), 500*time.Millisecond)
		if err != nil {
			continue
		}

		log.Printf("[%s]: %s", componentID, string(resp.Data))
	}
}
