package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/shared/config"
	"github.com/tnqbao/gau-account-service/shared/provider/dto"
)

type AuthorizationServiceProvider struct {
	AuthorizationServiceURL string `json:"authorization_service_url"`
	PrivateKey              string `json:"private_key,omitempty"`
}

func NewAuthorizationServiceProvider(config *config.EnvConfig) *AuthorizationServiceProvider {
	url := config.ExternalService.AuthorizationServiceURL
	if url == "" {
		panic("Authorization service URL is not configured")
	}

	privateKey := config.PrivateKey
	if privateKey == "" {
		panic("Private key is not configured")
	}

	return &AuthorizationServiceProvider{
		AuthorizationServiceURL: url,
		PrivateKey:              privateKey,
	}
}

func (p *AuthorizationServiceProvider) CreateNewToken(userID uuid.UUID, permission string, deviceID string) (string, string, time.Time, error) {
	if deviceID == "" {
		return "", "", time.Time{}, fmt.Errorf("device ID is required")
	}

	url := fmt.Sprintf("%s/api/v2/authorization/token", p.AuthorizationServiceURL)

	body, err := json.Marshal(dto.CreateTokenRequest{
		UserID:     userID,
		Permission: permission,
	})

	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Device-ID", deviceID)
	req.Header.Set("Private-Key", p.PrivateKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", "", time.Time{}, fmt.Errorf("authorization service returned %d: %s", resp.StatusCode, string(raw))
	}

	var tokenResp dto.CreateTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to decode response: %w", err)
	}

	expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return tokenResp.AccessToken, tokenResp.RefreshToken, expiry, nil
}

func (p *AuthorizationServiceProvider) RenewAccessToken(refreshToken, deviceID, oldAccessToken string) (string, time.Time, error) {
	url := fmt.Sprintf("%s/api/v2/authorization/token/renew", p.AuthorizationServiceURL)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Device-ID", deviceID)
	req.Header.Set("X-Refresh-Token", refreshToken)
	req.Header.Set("X-Old-Access-Token", oldAccessToken)
	req.Header.Set("Private-Key", p.PrivateKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", time.Time{}, fmt.Errorf("authorization service returned %d: %s", resp.StatusCode, string(raw))
	}

	var response dto.RenewTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to decode response: %w", err)
	}

	expiry := time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	return response.AccessToken, expiry, nil
}

func (p *AuthorizationServiceProvider) CheckAccessToken(token string) error {
	url := fmt.Sprintf("%s/api/v2/authorization/token/validate?token=%s", p.AuthorizationServiceURL, token)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Private-Key", p.PrivateKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("invalid token: %s", string(raw))
	}

	return nil
}

func (p *AuthorizationServiceProvider) RevokeToken(refreshToken, deviceID string) error {
	url := fmt.Sprintf("%s/api/v2/authorization/token", p.AuthorizationServiceURL)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Refresh-Token", refreshToken)
	req.Header.Set("X-Device-ID", deviceID)
	req.Header.Set("Private-Key", p.PrivateKey)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authorization service returned %d: %s", resp.StatusCode, string(raw))
	}

	return nil
}
