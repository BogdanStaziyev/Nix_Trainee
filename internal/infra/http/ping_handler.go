package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func PingHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode("OK")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}
