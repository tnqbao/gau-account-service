package provider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tnqbao/gau-account-service/shared/provider/dto"
)

func GetUserInfoFromGoogle(token string) (*dto.GoogleUserInfo, error) {

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call google api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api returned status: %d", resp.StatusCode)
	}

	var gResp dto.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return nil, fmt.Errorf("failed to decode google response: %w", err)
	}

	return &gResp, nil

}
