package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Junior580/go-barber-api/internal/whatsapp/handlers"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.VerifyWebhook(w, r)
		case http.MethodPost:
			handlers.ReceiveMessage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/available-days", handlers.AvailableDays)
	http.HandleFunc("/available-hours", handlers.GetAvailableHours)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s 🚀", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
