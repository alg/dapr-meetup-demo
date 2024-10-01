package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/client"
)

func main() {
	http.HandleFunc("POST /order", handlePlaceOrder)

	log.Printf("Listening on 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// Places the order if warehouse has enough items.
func handlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	cl, err := client.NewClient()
	if err != nil {
		log.Fatalf("error initializing client: %v", err)
	}

	content := &client.DataContent{
		ContentType: "application/json",
		Data:        []byte(`{"items": 1}`),
	}

	out, err := cl.InvokeMethodWithContent(r.Context(), "warehouse", "stock:decrease", "post", content)
	if err != nil {
		handleError(w, err)
		return
	}

	var resp struct {
		Success      bool   `json:"success"`
		ErrorMessage string `json:"errorMessage"`
	}

	if err := json.Unmarshal(out, &resp); err != nil {
		http.Error(w, fmt.Sprintf("error parsing warehouse response: %v", err), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, string(out), http.StatusPreconditionFailed)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
