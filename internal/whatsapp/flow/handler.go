package flow

import (
	"encoding/json"
	"net/http"
	"time"
)

type FlowHandler struct {
	Crypto *Crypto
}

func NewFlowHandler() (*FlowHandler, error) {
	crypto, err := NewCrypto()
	if err != nil {
		return nil, err
	}

	return &FlowHandler{
		Crypto: crypto,
	}, nil
}

func (h *FlowHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		EncryptedFlowData string `json:"encrypted_flow_data"`
		EncryptedAESKey   string `json:"encrypted_aes_key"`
		InitialVector     string `json:"initial_vector"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// 1. DESCRYPT AES KEY
	aesKey, err := h.Crypto.DecryptAESKey(req.EncryptedAESKey)
	if err != nil {
		w.WriteHeader(421)
		return
	}

	// 2. DECRYPT PAYLOAD
	payload, err := h.Crypto.DecryptPayload(aesKey, req.InitialVector, req.EncryptedFlowData)
	if err != nil {
		w.WriteHeader(421)
		return
	}

	action := payload["action"].(string)
	screen := ""
	if payload["screen"] != nil {
		screen = payload["screen"].(string)
	}

	data := map[string]any{}
	if payload["data"] != nil {
		data = payload["data"].(map[string]any)
	}

	flowToken := payload["flow_token"].(string)

	response := map[string]any{}

	// ---------------------------
	// BUSINESS LOGIC
	// ---------------------------

	// FIRST SCREEN (INIT)
	if action == "INIT" {
		response = map[string]any{
			"screen": "SELECT_DAY",
			"data": map[string]any{
				"days": generateDays(),
			},
		}
	}

	// SECOND SCREEN
	if action == "data_exchange" && screen == "SELECT_DAY" {
		selected := data["day"].(string)
		response = map[string]any{
			"screen": "SELECT_HOUR",
			"data": map[string]any{
				"day":   selected,
				"hours": generateHours(),
			},
		}
	}

	// FINALIZATION
	if action == "data_exchange" && screen == "SELECT_HOUR" {
		response = map[string]any{
			"screen": "SUCCESS",
			"data": map[string]any{
				"extension_message_response": map[string]any{
					"params": map[string]any{
						"flow_token": flowToken,
						"day":        data["day"],
						"hour":       data["hour"],
					},
				},
			},
		}
	}

	// 3. ENCRYPT RESPONSE
	encrypted, err := h.Crypto.EncryptPayload(aesKey, req.InitialVector, response)
	if err != nil {
		http.Error(w, "internal server error", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{
		"encrypted_flow_data": encrypted,
	}); err != nil {
		http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func generateDays() []map[string]string {
	days := []map[string]string{}

	for i := range 5 {
		d := time.Now().AddDate(0, 0, i)
		days = append(days, map[string]string{
			"id":    d.Format("2006-01-02"),
			"title": d.Format("02/01/2006"),
		})
	}
	return days
}

func generateHours() []map[string]string {
	return []map[string]string{
		{"id": "09:00", "title": "09:00"},
		{"id": "10:00", "title": "10:00"},
		{"id": "11:00", "title": "11:00"},
		{"id": "13:00", "title": "13:00"},
		{"id": "15:00", "title": "15:00"},
	}
}
