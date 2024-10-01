package main

import (
	"context"
	"log"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
)

func main() {
	s, err := daprd.NewService(":8002")
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}

	// Add out-of-stock handler
	sub := &common.Subscription{PubsubName: "pubsub", Topic: "out-of-stock"}
	if err := s.AddTopicEventHandler(sub, handleOutOfStock); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	// Add low-stock handler
	sub = &common.Subscription{PubsubName: "pubsub", Topic: "low-stock"}
	if err := s.AddTopicEventHandler(sub, handleLowStock); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	// Start the service
	if err := s.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleLowStock(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	// send email
	in := &dapr.InvokeBindingRequest{
		Name:      "email-low-stock",
		Operation: "create",
		Data:      []byte("Hey, we are getting really low. Resupply!"),
		Metadata:  map[string]string{"emailTo": "Supplier <supplier@example.com>"},
	}

	client, err := dapr.NewClient()
	if err != nil {
		return true, err
	}

	_, err = client.InvokeBinding(ctx, in)
	if err != nil {
		return true, err
	}

	return false, nil
}
func handleOutOfStock(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	// send email
	in := &dapr.InvokeBindingRequest{
		Name:      "email-out-of-stock",
		Operation: "create",
		Data:      []byte("OMG, we are out of stock!"),
		Metadata:  map[string]string{"emailTo": "Supplier <supplier@example.com>"},
	}

	client, err := dapr.NewClient()
	if err != nil {
		return true, err
	}

	_, err = client.InvokeBinding(ctx, in)
	if err != nil {
		return true, err
	}

	return false, nil
}
