package main

import (
	"context"
	pb "coop/proto"
	"log"
	"sync"
)

var (
	subscriptions      = make(map[string]int)
	subscriptionsMutex sync.Mutex
)

func addSubscription(name string) {
	log.Println("Add subsription req:", name)
	subscriptionsMutex.Lock()
	if count, ok := subscriptions[name]; ok {
		log.Println("Incrementing subsription:", name)
		subscriptions[name] = count + 1
	} else {
		log.Println("Adding subsription:", name)
		subscriptions[name] = 1
	}
	subscriptionsMutex.Unlock()
}

func getSubscription(name string) string {
	sub := ""
	subscriptionsMutex.Lock()
	if count, ok := subscriptions[name]; ok && count > 0 {
		sub = name
	}
	subscriptionsMutex.Unlock()
	return sub
}

func maybeRemoveSubscription(name string) {
	subscriptionsMutex.Lock()
	if count, ok := subscriptions[name]; ok && (count-1) == 0 {
		log.Println("Removing subsription:", name)
		delete(subscriptions, name)
	}
	subscriptionsMutex.Unlock()
}

// subscribe
func sub(c pb.EventClient, tag string) error {
	defer maybeRemoveSubscription(tag)

	addSubscription(tag)

	ctx := context.Background()
	defer ctx.Done()

	stream, err := c.Sub(ctx, &pb.SubscribeRequest{
		Tag: tag,
	})
	if err != nil {
		log.Println("Err with stream", err)
		return err
	}

	for {
		event, err := stream.Recv()

		if err != nil {
			continue
		}
		eventMessage <- EventMessage{
			event: event.Event,
			data:  event.Data,
			tag:   event.Tag,
		}
	}

}
