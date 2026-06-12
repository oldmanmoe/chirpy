package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(header http.Header) (string, error) {
	rawHeader := header.Get("Authorization")

	if rawHeader == "" {
		return "", fmt.Errorf("Unable to get header string")
	}

	polkaApiKey := strings.ReplaceAll(rawHeader, "ApiKey ", "")
	if polkaApiKey == "" {
		return "", fmt.Errorf("Unable to get api key")
	}

	return polkaApiKey, nil
}