package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func VerifyWebhook(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, r.URL.Query().Get("hub.challenge"))
}

func ReceiveMessage(w http.ResponseWriter, r *http.Request) {
	var data map[string]any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	log.Printf("Message received: %+v\n", data)
	w.WriteHeader(http.StatusOK)
}
