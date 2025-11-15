package flow

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (h *FlowHandler) Handle(c *gin.Context) {
	var req struct {
		EncryptedFlowData string `json:"encrypted_flow_data"`
		EncryptedAESKey   string `json:"encrypted_aes_key"`
		InitialVector     string `json:"initial_vector"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// 1. DESCRIPTOGRAFAR AES KEY
	aesKey, err := h.Crypto.DecryptAESKey(req.EncryptedAESKey)
	if err != nil {
		c.JSON(421, gin.H{})
		return
	}

	// 2. DESCRIPTOGRAFAR PAYLOAD
	payload, err := h.Crypto.DecryptPayload(aesKey, req.InitialVector, req.EncryptedFlowData)
	if err != nil {
		c.JSON(421, gin.H{})
		return
	}

	action := payload["action"].(string)
	screen := ""
	if payload["screen"] != nil {
		screen = payload["screen"].(string)
	}

	data := map[string]interface{}{}
	if payload["data"] != nil {
		data = payload["data"].(map[string]interface{})
	}

	flowToken := payload["flow_token"].(string)

	// ----------------------------
	//   BUSINESS LOGIC DO FLOW
	// ----------------------------

	response := map[string]interface{}{}

	// PRIMEIRA TELA (INIT)
	if action == "INIT" {
		response = map[string]interface{}{
			"screen": "SELECT_DAY",
			"data": map[string]interface{}{
				"dias": generateDays(),
			},
		}
	}

	// SEGUNDA TELA
	if action == "data_exchange" && screen == "SELECT_DAY" {
		selected := data["dia"].(string)
		response = map[string]interface{}{
			"screen": "SELECT_HOUR",
			"data": map[string]interface{}{
				"dia":      selected,
				"horarios": generateHours(),
			},
		}
	}

	// FINALIZAÇÃO
	if action == "data_exchange" && screen == "SELECT_HOUR" {
		response = map[string]interface{}{
			"screen": "SUCCESS",
			"data": map[string]interface{}{
				"extension_message_response": map[string]interface{}{
					"params": map[string]interface{}{
						"flow_token": flowToken,
						"dia":        data["dia"],
						"horario":    data["horario"],
					},
				},
			},
		}
	}

	// 3. CRIPTOGRAFAR RESPOSTA
	encrypted, err := h.Crypto.EncryptPayload(aesKey, req.InitialVector, response)
	if err != nil {
		c.JSON(500, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"encrypted_flow_data": encrypted,
	})
}

func generateDays() []map[string]string {
	days := []map[string]string{}

	for i := 0; i < 5; i++ {
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
