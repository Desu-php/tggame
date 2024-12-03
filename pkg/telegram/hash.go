package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func ValidateInitData(initData, botToken string) bool {
	// Parse the initData query string
	parsedData, err := url.ParseQuery(initData)
	if err != nil {
		fmt.Println("Error parsing initData:", err)
		return false
	}

	// Extract the hash parameter
	receivedHash := parsedData.Get("hash")
	if receivedHash == "" {
		fmt.Println("Hash not found in initData")
		return false
	}
	parsedData.Del("hash") // Remove the hash field

	// Construct the data_check_string
	var keys []string
	for key := range parsedData {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var dataCheckStringBuilder strings.Builder
	for _, key := range keys {
		dataCheckStringBuilder.WriteString(fmt.Sprintf("%s=%s\n", key, parsedData.Get(key)))
	}
	dataCheckString := dataCheckStringBuilder.String()
	dataCheckString = strings.TrimSuffix(dataCheckString, "\n") // Remove trailing newline

	// Generate the secret_key
	botTokenKey := []byte("WebAppData")
	h := hmac.New(sha256.New, botTokenKey)
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	// Generate the hash of the data_check_string
	h = hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	// Compare the calculated hash with the received hash
	return calculatedHash == receivedHash
}

