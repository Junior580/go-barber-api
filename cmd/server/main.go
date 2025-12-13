package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Junior580/go-barber-api/configs"
)

func main() {
	env, err := configs.LoadConfig("../../")
	if err != nil {
		log.Print(err)
	}
	pemBytes, err := os.ReadFile(env.PRIVATE_KEY_PATH)
	pemContent := string(pemBytes)
	if err != nil {
		log.Fatalf("erro ao ler private key: %v", err)
	}
	fmt.Printf("privateKey: %v - passphrase: %v \n", env.PASSPHRASE, pemContent)
}
