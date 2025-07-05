package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

func EncodeSecrets(secrets map[string]string, encryption_key string) ([]byte, error) {
	secretb, err := json.Marshal(secrets)
	if err != nil {
		return nil, fmt.Errorf("error marshaling secvrets: %v", err)
	}

	ekey, err := hex.DecodeString(encryption_key)
	if err != nil {
		return nil, fmt.Errorf("error decoding encryption key: %v", err)
	}

	block, err := aes.NewCipher(ekey)
	if err != nil {
		return nil, fmt.Errorf("error generating aes cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating new gcm: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {

	}
	ciphertext := gcm.Seal(nonce, nonce, secretb, nil)
	return ciphertext, nil
}

func DecodeSecrets(ciphertext []byte, encryption_key string) (map[string]string, error) {
	ekey, err := hex.DecodeString(encryption_key)
	if err != nil {
		return nil, fmt.Errorf("error decoding encryption key: %v", err)
	}

	block, err := aes.NewCipher(ekey)
	if err != nil {
		return nil, fmt.Errorf("error generating aes cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating new gcm: %v", err)
	}

	decryptedData, err := gcm.Open(nil, ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():], nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %v", err)
	}

	var decryptedDataDict map[string]string
	err = json.Unmarshal(decryptedData, &decryptedData)
	if err != nil {
		return nil, fmt.Errorf("error converting decrypted data to dictionary: %v", err)
	}
	return decryptedDataDict, nil
}
