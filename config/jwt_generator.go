package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateSecureJWTSecret generates a secure JWT secret key
func GenerateSecureJWTSecret() (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	// Encode to base64
	secret := base64.URLEncoding.EncodeToString(bytes)
	return secret, nil
}

// IsSecureJWTSecret checks if the JWT secret is secure (not the default weak one)
func IsSecureJWTSecret(secretKey string) bool {
	// Check if it's the default weak secret
	if secretKey == "123456" {
		return false
	}

	// Check if the secret is at least 32 characters (256 bits)
	if len(secretKey) < 32 {
		return false
	}

	return true
}