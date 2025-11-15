package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Junior580/go-barber-api/internal/whatsapp/types"
)

func AvailableDays(w http.ResponseWriter, r *http.Request) {
	days := []types.Option{}
	daysOfWeek := []string{"Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"}

	for i := range 5 {
		day := time.Now().AddDate(0, 0, i)
		id := day.Format("2006-01-02") // ISO date
		weekday := daysOfWeek[day.Weekday()]

		// title in English style date but weekday in Portuguese
		title := fmt.Sprintf("%s (%s)", day.Format("2006-01-02"), weekday)

		days = append(days, types.Option{
			ID:    id,
			Title: title,
		})
	}

	resp := types.Response{}
	resp.Data.Options = days

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
