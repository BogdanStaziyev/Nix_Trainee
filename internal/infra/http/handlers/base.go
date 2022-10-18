package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ok(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func success(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Print(err)
	}
}

func created(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Print(err)
	}
}

//func noContent(w http.ResponseWriter) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusNoContent)
//}

func BadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	encodeErrorBody(w, err)
}

//func Forbidden(w http.ResponseWriter, err error) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusForbidden)
//
//	encodeErrorBody(w, err)
//}

func InternalServerError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	encodeErrorBody(w, err)
}

//nolint
//func validationError(w http.ResponseWriter, err error) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusUnprocessableEntity)
//
//	encodeErrorBody(w, err)
//}

//nolint
//func genericError(w http.ResponseWriter, err error) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusBadRequest)
//
//	encodeErrorBody(w, err)
//}

func NotFound(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	body := "Not Found"

	if err != nil {
		body = fmt.Sprint(err)
	}

	e := json.NewEncoder(w).Encode(map[string]interface{}{"error": body})
	if e != nil {
		log.Print(e)
	}
}

func encodeErrorBody(w http.ResponseWriter, err error) {
	e := json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
	if e != nil {
		log.Print(e)
	}
}
