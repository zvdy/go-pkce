package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// GenerateCodeVerifier generates a random code verifier.
func GenerateCodeVerifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// GenerateCodeChallenge generates a code challenge from the code verifier.
func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// ValidatePKCE validates the provided code verifier and challenge.
func ValidatePKCE(verifier, challenge string) error {
	expectedChallenge := GenerateCodeChallenge(verifier)
	if expectedChallenge != challenge {
		return errors.New("invalid PKCE parameters")
	}
	return nil
}
