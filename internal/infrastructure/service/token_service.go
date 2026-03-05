package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func GenerateRefreshToken() (token string, id string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("Failed to generate random bytes for refresh token: %w", err)
	}
	token = base64.StdEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))
	id = hex.EncodeToString(hash[:])
	return token, id, nil
}
