// internal/infrastructure/service/token_generator.go
package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
)


func GenerateRefreshToken() (token string, id string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", appErrors.Internal("Failed to generate refresh token", err)
	}
	token = base64.StdEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))
	id = hex.EncodeToString(hash[:])
	return token, id, nil
}
