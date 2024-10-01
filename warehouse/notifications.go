package main

import (
	"context"
	"encoding/json"

	dapr "github.com/dapr/go-sdk/client"
)

const Pubsub = "pubsub"
const TopicLowStock = "low-stock"
const TopicOutOfStock = "out-of-stock"

func notifyLowStock(ctx context.Context, wh *Warehouse) error {
	data, _ := json.Marshal(map[string]any{"remaning": wh.Stock})
	return notify(ctx, TopicLowStock, data)
}

func notifyOutOfStock(ctx context.Context) error {
	return notify(ctx, TopicOutOfStock, []byte{})
}

func notify(ctx context.Context, topic string, data []byte) error {
	client, err := dapr.NewClient()
	if err != nil {
		return err
	}

	err = client.PublishEvent(ctx, Pubsub, topic, data)

	return err
}
