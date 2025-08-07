package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tnqbao/gau-account-service/config"
	"mime/multipart"
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

func (p *UploadServiceProvider) UploadAvatarImage(userID string, imageData []byte, filename string) (string, error) {
	url := fmt.Sprintf("%s/api/v2/upload/image", p.UploadServiceURL)

	// Prepare multipart form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add file_path field
	if err := w.WriteField("file_path", "avatar"); err != nil {
		return "", fmt.Errorf("failed to write file_path field: %w", err)
	}

	// Add file field
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := fw.Write(imageData); err != nil {
		return "", fmt.Errorf("failed to write image data: %w", err)
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPatch, url, &b)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Private-Key", p.PrivateKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload service returned status: %d", resp.Body)
	}
	var response struct {
		FilePath string `json:"file_path"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	return response.FilePath, nil
}
