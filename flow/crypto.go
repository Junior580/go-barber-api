package flow

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
)

type Crypto struct {
	PrivateKey *rsa.PrivateKey
}

func NewCrypto() (*Crypto, error) {
	keyData, err := os.ReadFile("private.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &Crypto{PrivateKey: privateKey}, nil
}

func (c *Crypto) DecryptAESKey(encKeyBase64 string) ([]byte, error) {
	encKey, err := base64.StdEncoding.DecodeString(encKeyBase64)
	if err != nil {
		return nil, err
	}

	aesKey, err := rsa.DecryptPKCS1v15(nil, c.PrivateKey, encKey)
	if err != nil {
		return nil, err
	}

	return aesKey, nil
}

func (c *Crypto) DecryptPayload(aesKey []byte, ivBase64 string, encDataBase64 string) (map[string]interface{}, error) {
	iv, _ := base64.StdEncoding.DecodeString(ivBase64)
	encryptedData, _ := base64.StdEncoding.DecodeString(encDataBase64)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encryptedData))
	mode.CryptBlocks(decrypted, encryptedData)

	// remove padding
	for len(decrypted) > 0 && decrypted[len(decrypted)-1] == 0 {
		decrypted = decrypted[:len(decrypted)-1]
	}

	var result map[string]interface{}
	err = json.Unmarshal(decrypted, &result)
	return result, err
}

func (c *Crypto) EncryptPayload(aesKey []byte, ivBase64 string, payload interface{}) (string, error) {
	iv, _ := base64.StdEncoding.DecodeString(ivBase64)

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	padded := jsonBytes
	for len(padded)%aes.BlockSize != 0 {
		padded = append(padded, 0)
	}

	encrypted := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, padded)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}
