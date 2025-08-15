package provider

import (
	"encoding/json"
	"fmt"
	"github.com/tnqbao/gau-account-service/entity"
	"github.com/tnqbao/gau-account-service/provider/dto"
	"net/http"
	"strings"
)

// GetUserInfoFromGoogle gets user info from Google OAuth API
func GetUserInfoFromGoogle(token string) (*entity.User, error) {
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

	// Generate username from fullname: remove spaces, uppercase
	var username *string
	if gResp.Name != "" {
		generatedUsername := generateUsernameFromFullName(gResp.Name)
		username = &generatedUsername
	}

	user := &entity.User{
		Email:           &gResp.Email,
		FullName:        &gResp.Name,
		AvatarURL:       &gResp.Picture,
		Username:        username,
		IsEmailVerified: gResp.EmailVerified,
	}

	return user, nil
}

// generateUsernameFromFullName generates base username from fullname: no space, uppercase
func generateUsernameFromFullName(fullName string) string {
	if fullName == "" {
		return ""
	}

	// Remove spaces and convert to uppercase
	return strings.ToUpper(strings.ReplaceAll(fullName, " ", ""))
}
