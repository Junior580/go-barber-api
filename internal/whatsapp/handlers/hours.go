package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Junior580/go-barber-api/internal/whatsapp/types"
)

func GetAvailableHours(w http.ResponseWriter, r *http.Request) {
	day := r.URL.Query().Get("day") // keeps 'dia' only because your flow expects this
	if day == "" {
		http.Error(w, "Missing parameter: day", http.StatusBadRequest)
		return
	}

	// generate hours from 09:00 to 19:00
	hours := []types.Option{}

	for hour := 9; hour <= 19; hour++ {
		h := fmt.Sprintf("%02d:00", hour)

		hours = append(hours, types.Option{
			ID:    h,
			Title: h,
		})
	}

	resp := types.Response{}
	resp.Data.Options = hours

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
