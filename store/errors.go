package main

import (
	"fmt"
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func handleError(w http.ResponseWriter, err error) {
	st := status.Convert(err)

	fmt.Printf("Code: %s\n", st.Code().String())
	fmt.Printf("Message: %s\n", st.Message())

	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.ErrorInfo:
			fmt.Printf("ErrorInfo:\n- Domain: %s\n- Reason: %s\n- Metadata: %v\n", t.GetDomain(), t.GetReason(), t.GetMetadata())
			metadata := t.GetMetadata()
			http.Error(w, metadata["http.error_message"], http.StatusUnprocessableEntity)

		default:
			fmt.Printf("Unhandled error detail type: %v\n", t)
			http.Error(w, fmt.Sprintf("error contacting warehouse: %v", err), http.StatusInternalServerError)
		}
	}
}
