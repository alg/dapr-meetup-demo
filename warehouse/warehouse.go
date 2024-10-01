package main

import (
	"context"
	"encoding/json"

	dapr "github.com/dapr/go-sdk/client"
)

const LowThreshold = 3
const StateStoreName = "statestore"
const StateKey = "warehouse"

type ErrOutOfStock struct{}

func (e *ErrOutOfStock) Error() string {
	return "Warehouse is out of stock"
}

type Warehouse struct {
	Stock int `json:"stock"`
}

// Decrease attempts to use some items from stock.
func (wh *Warehouse) Decrease(items int) error {
	if wh.Stock < items {
		return &ErrOutOfStock{}
	}

	wh.Stock -= items

	return nil
}

func (wh *Warehouse) IsStockLow() bool {
	return wh.Stock < LowThreshold
}

func (wh *Warehouse) IsOutOfStock() bool {
	return wh.Stock == 0
}

// FetchWarehouse fetches the state of the warehouse
func FetchWarehouse(ctx context.Context) (wh *Warehouse, etag string, err error) {
	client, err := dapr.NewClient()
	if err != nil {
		return nil, "", err
	}

	item, err := client.GetState(ctx, StateStoreName, StateKey, map[string]string{"content-type": "application/json"})
	if err != nil {
		return nil, "", err
	}

	if len(item.Value) == 0 {
		return &Warehouse{Stock: 5}, "", nil
	}

	err = json.Unmarshal(item.Value, &wh)
	if err != nil {
		return nil, "", err
	}

	return wh, item.Etag, nil
}

func StoreWarehouse(ctx context.Context, wh *Warehouse, etag string) error {
	client, err := dapr.NewClient()
	if err != nil {
		return err
	}

	dataBytes, err := json.Marshal(wh)
	if err != nil {
		return err
	}

	err = client.SaveState(
		ctx,
		StateStoreName,
		StateKey,
		dataBytes,
		map[string]string{"content-type": "application/json"},
	)
	if err != nil {
		return err
	}

	return nil
}
