package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tnqbao/gau-account-service/config"
	"net/http"
)

type UploadServiceProvider struct {
	UploadServiceURL string `json:"upload_service_url"`
	PrivateKey       string `json:"private_key,omitempty"`
}

func NewUploadServiceProvider(config *config.EnvConfig) *UploadServiceProvider {
	if config.ExternalService.UploadServiceURL == "" {
		panic("Upload service URL is not configured")
	}

	if config.PrivateKey == "" {
		panic("Private key is not configured")
	}

	return &UploadServiceProvider{
		UploadServiceURL: config.ExternalService.UploadServiceURL,
		PrivateKey:       config.PrivateKey,
	}
}

func (p *UploadServiceProvider) UploadAvatarImage(userID string, imageData []byte) (string, error) {
	url := fmt.Sprintf("%s/api/v2/upload/image", p.UploadServiceURL)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(imageData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Private-Key", p.PrivateKey)
	req.Header.Set("X-User-ID", userID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload service returned status: %d", resp.StatusCode)
	}
	var response struct {
		filePath string `json:"file_path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return response.filePath, nil
}
