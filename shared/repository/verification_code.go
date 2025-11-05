package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func (r *Repository) GenerateVerificationToken(ctx context.Context, userID string, email string) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Store token in Redis with 24 hour expiration
	key := fmt.Sprintf("email_verification:%s", token)
	value := fmt.Sprintf("%s:%s", userID, email)

	if err := r.cacheDb.Set(ctx, key, value, 24*time.Hour).Err(); err != nil {
		return "", fmt.Errorf("failed to store token: %w", err)
	}

	return token, nil
}

// ValidateVerificationToken validates and retrieves user info from token
func (r *Repository) ValidateVerificationToken(ctx context.Context, token string) (userID string, email string, err error) {
	key := fmt.Sprintf("email_verification:%s", token)

	value, err := r.cacheDb.Get(ctx, key).Result()
	if err != nil {
		return "", "", fmt.Errorf("invalid or expired token: %w", err)
	}

	// Parse value (format: "userID:email")
	var parsedUserID, parsedEmail string
	if _, err := fmt.Sscanf(value, "%s:%s", &parsedUserID, &parsedEmail); err != nil {
		return "", "", fmt.Errorf("invalid token format: %w", err)
	}

	// Delete token after successful validation (one-time use)
	r.cacheDb.Del(ctx, key)

	return parsedUserID, parsedEmail, nil
}

func (r *Repository) GetImage(ctx context.Context, key string) ([]byte, string, error) {
	data, err := r.cacheDb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, "", err
	}
	ct, err := r.cacheDb.Get(ctx, key+":content-type").Result()
	if err != nil {
		ct = "application/octet-stream"
	}
	return data, ct, nil
}
