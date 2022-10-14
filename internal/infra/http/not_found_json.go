package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func NotFoundJSON() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content_type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(writer).Encode("Resource Not Found")
		if err != nil {
			fmt.Printf("writting response: %s", err)
		}
	}
}
