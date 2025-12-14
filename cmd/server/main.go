package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Junior580/go-barber-api/configs"
)

const nonceSize = 16

type endpointPayload struct {
	EncryptedAESKey   string `json:"encrypted_aes_key"`
	EncryptedFlowData string `json:"encrypted_flow_data"`
	InitialVector     string `json:"initial_vector"`
}

type decryptionResult struct {
	DecryptedBody      map[string]any
	AESKeyBytes        []byte
	InitialVectorBytes []byte
}

func main() {
	env, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("❌ error load env file: %v", err)
	}

	pemBytes, err := os.ReadFile(env.PRIVATE_KEY_PATH)
	privateKey := string(pemBytes)
	passphrase := env.PASSPHRASE
	if err != nil {
		log.Fatalf("❌ error load env file: %v", err)
	}

	if privateKey == "" || passphrase == "" {
		log.Fatal("Environment variables 'PRIVATE_KEY' and 'PASSPHRASE' are required.")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		encryptedResponse, err := processRequestHTTP(r, privateKey, passphrase)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(encryptedResponse))
	})

	log.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func processRequestHTTP(r *http.Request, privateKey, passphrase string) (string, error) {
	var payload endpointPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return "", err
	}

	decrypted, err := decryptRequest(
		payload.EncryptedAESKey,
		payload.EncryptedFlowData,
		payload.InitialVector,
		privateKey,
		passphrase,
	)
	if err != nil {
		return "", err
	}

	action, ok := decrypted.DecryptedBody["action"].(string)
	if ok {
		log.Println("Action:", action)
	}

	// screen := "START" // fallback

	// if action == "ping" {
	// 	screen = "PING" // TEM que existir no Flow JSON
	// }

	response := map[string]any{
		"data": map[string]any{
			"status": "active",
		},
	}

	return encryptResponse(
		response,
		decrypted.AESKeyBytes,
		decrypted.InitialVectorBytes,
	)

	// response := map[string]any{
	// 	"screen": "SCREEN_NAME",
	// 	"data":   map[string]string{"some_key": "some_value"},
	// }

	// return encryptResponse(
	// 	response,
	// 	decrypted.AESKeyBytes,
	// 	decrypted.InitialVectorBytes,
	// )
}

func decryptRequest(encryptedAESKey string, encryptedFlowData string, initialVector string, privatePem string, passphrase string) (decryptionResult, error) {
	block, _ := pem.Decode([]byte(privatePem))

	if block == nil || !x509.IsEncryptedPEMBlock(block) {
		return decryptionResult{}, errors.New("invalid PEM format or not encrypted")
	}

	decryptedKey, err := x509.DecryptPEMBlock(block, []byte(passphrase))
	if err != nil {
		fmt.Println("decryptRequest-block: ", err)
		return decryptionResult{}, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(decryptedKey)
	if err != nil {
		return decryptionResult{}, err
	}

	encryptedAESKeyBytes, _ := base64.StdEncoding.DecodeString(encryptedAESKey)
	aesKeyBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedAESKeyBytes, nil)
	if err != nil {
		return decryptionResult{}, err
	}

	initialVectorBytes, _ := base64.StdEncoding.DecodeString(initialVector)
	flowDataBytes, _ := base64.StdEncoding.DecodeString(encryptedFlowData)

	blockCipher, err := aes.NewCipher(aesKeyBytes)
	if err != nil {
		return decryptionResult{}, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(blockCipher, nonceSize)
	if err != nil {
		return decryptionResult{}, err
	}

	decryptedPlaintext, err := gcm.Open(nil, initialVectorBytes, flowDataBytes, nil)
	if err != nil {
		return decryptionResult{}, err
	}

	var decryptedBody map[string]any
	if err := json.Unmarshal(decryptedPlaintext, &decryptedBody); err != nil {
		return decryptionResult{}, err
	}

	return decryptionResult{
		DecryptedBody:      decryptedBody,
		AESKeyBytes:        aesKeyBytes,
		InitialVectorBytes: initialVectorBytes,
	}, nil
}

func encryptResponse(response map[string]any, aesKeyBytes, initialVectorBytes []byte) (string, error) {
	flippedIV := make([]byte, len(initialVectorBytes))
	for i, b := range initialVectorBytes {
		flippedIV[i] = ^b
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	blockCipher, err := aes.NewCipher(aesKeyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCMWithNonceSize(blockCipher, nonceSize)
	if err != nil {
		return "", err
	}

	encryptedData := gcm.Seal(nil, flippedIV, jsonResponse, nil)
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}
