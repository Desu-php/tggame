package str

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateSessionKey() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "",fmt.Errorf("GenerateSessionKey: error %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
