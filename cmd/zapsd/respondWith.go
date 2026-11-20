package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, struct {
		Error string `json:"error"`
	}{Error: message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	msg, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshaling json: %v", err)
		return
	}

	w.WriteHeader(code)
	w.Write(msg)
}
