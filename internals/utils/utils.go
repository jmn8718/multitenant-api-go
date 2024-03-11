package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

func GenerateApiKey(apiKeyLength int8, prefix string) (string, error) {
	b := make([]byte, apiKeyLength+4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	key := strings.Replace(strings.Replace(base64.URLEncoding.EncodeToString(b), "-", "", -1), "_", "", -1)[:apiKeyLength]
	return fmt.Sprintf("%s_%s", prefix, key), nil
}

func ToValidDateString(valueDate time.Time) string {
	if valueDate.IsZero() {
		return ""
	}
	return valueDate.UTC().Format(time.RFC3339)
}
