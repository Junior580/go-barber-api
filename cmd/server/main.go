package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func verifyWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.URL.Query().Get("hub.challenge"))
}

func receiveMessage(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}
	log.Printf("Mensagem recebida: %+v\n", data)
	w.WriteHeader(http.StatusOK)
}

type Option struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Response struct {
	Data struct {
		Options []Option `json:"options"`
	} `json:"data"`
}

func availableDays(w http.ResponseWriter, r *http.Request) {
	days := []Option{}
	daysOfWeek := []string{"Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"}

	for i := 0; i < 5; i++ {
		day := time.Now().AddDate(0, 0, i)
		id := day.Format("2006-01-02") // ISO date
		weekday := daysOfWeek[day.Weekday()]

		// title in English style date but weekday in Portuguese
		title := fmt.Sprintf("%s (%s)", day.Format("2006-01-02"), weekday)

		days = append(days, Option{
			ID:    id,
			Title: title,
		})
	}

	resp := Response{}
	resp.Data.Options = days

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			verifyWebhook(w, r)
		} else if r.Method == http.MethodPost {
			receiveMessage(w, r)
		} else {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/dias-disponiveis", availableDays)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor iniciado na porta %s 🚀", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
