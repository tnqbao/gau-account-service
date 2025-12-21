package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/tnqbao/gau-account-service/shared/config"
)

type UploadServiceProvider struct {
	UploadServiceURL string `json:"upload_service_url"`
	CDNServiceURL    string `json:"cdn_service_url"`
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
		CDNServiceURL:    config.ExternalService.CDNServiceURL,
		PrivateKey:       config.PrivateKey,
	}
}

// UploadResponse represents the response from upload service
type UploadResponse struct {
	Bucket      string `json:"bucket"`
	ContentType string `json:"content_type"`
	Duplicated  bool   `json:"duplicated"`
	FileHash    string `json:"file_hash"`
	FilePath    string `json:"file_path"`
	Message     string `json:"message"`
	Size        int64  `json:"size"`
	Status      int    `json:"status"`
}

func (p *UploadServiceProvider) UploadAvatarImage(imageData []byte, filename string, contentType string) (string, error) {
	// Sử dụng API mới POST /api/v2/upload/file
	url := fmt.Sprintf("%s/api/v2/upload/file", p.UploadServiceURL)

	// Prepare multipart form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add bucket field
	if err := w.WriteField("bucket", "images"); err != nil {
		return "", fmt.Errorf("failed to write bucket field: %w", err)
	}

	Add path field
	if err := w.WriteField("path", "avatar"); err != nil {
		return "", fmt.Errorf("failed to write path field: %w", err)
	}

	// Add file field with proper content type
	h := make(map[string][]string)
	h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename)}
	h["Content-Type"] = []string{contentType}

	fw, err := w.CreatePart(h)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := fw.Write(imageData); err != nil {
		return "", fmt.Errorf("failed to write image data: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &b)
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
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload service returned %d: %s", resp.StatusCode, string(raw))
	}

	var response UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	cdnURL := fmt.Sprintf("%s/%s/%s", p.CDNServiceURL, response.Bucket, response.FilePath)

	return cdnURL, nil
}
