package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("POST /stock:renew", handleRenewStock)
	http.HandleFunc("POST /stock:decrease", handleDecreaseStock)
	http.HandleFunc("GET /stock", handleGetStock)

	log.Printf("Listening on 8001")
	log.Fatal(http.ListenAndServe(":8001", nil))
}

// Restocks the warehouse
func handleRenewStock(w http.ResponseWriter, r *http.Request) {
	withWarehouse(r.Context(), w, func(wh *Warehouse) (bool, error) {
		wh.Stock = 5

		return true, nil
	})
}

// Attempts to decrease stock and sends the OutOfStock event when nothing left.
func handleDecreaseStock(w http.ResponseWriter, r *http.Request) {
	// parse request
	var req struct {
		Items int `json:"items"`
	}

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading request: %v", err), http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(reqBytes, &req); err != nil {
		http.Error(w, fmt.Sprintf("error parsing warehouse response: %v", err), http.StatusInternalServerError)
		return
	}

	wh := withWarehouse(r.Context(), w, func(wh *Warehouse) (changed bool, err error) {
		err = wh.Decrease(req.Items)

		if err != nil {
			_, ok := err.(*ErrOutOfStock)
			if ok {
				notifyOutOfStock(r.Context())
			}

			respBytes, _ := json.Marshal(map[string]any{"errorMessage": err.Error()})
			http.Error(w, string(respBytes), http.StatusPreconditionFailed)

			return false, err
		}

		return true, nil
	})

	if wh == nil {
		return
	}

	if wh.IsOutOfStock() {
		notifyOutOfStock(r.Context())
	} else if wh.IsStockLow() {
		notifyLowStock(r.Context(), wh)
	}

	w.WriteHeader(http.StatusCreated)
	respBytes, _ := json.Marshal(map[string]any{"remaining": wh.Stock})
	w.Write(respBytes)

}

// Returns current stock
func handleGetStock(w http.ResponseWriter, r *http.Request) {
	withWarehouse(r.Context(), w, func(wh *Warehouse) (bool, error) {
		respBytes, _ := json.Marshal(map[string]any{"items": wh.Stock})
		w.Write(respBytes)

		return false, nil
	})
}

func withWarehouse(ctx context.Context, w http.ResponseWriter, fn func(*Warehouse) (bool, error)) *Warehouse {
	// fetch warehouse
	wh, etag, err := FetchWarehouse(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching warehouse state: %v", err), http.StatusInternalServerError)
		return nil
	}

	changed, err := fn(wh)

	if err != nil {
		return nil
	}

	if changed {
		err = StoreWarehouse(ctx, wh, etag)
		if err != nil {
			http.Error(w, fmt.Sprintf("error storing warehouse state: %v", err), http.StatusInternalServerError)
			return nil
		}
	}

	return wh
}
