package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mozillazg/go-unidecode"
	"github.com/tnqbao/gau-account-service/shared/entity"
	"gorm.io/gorm"
)

func (ctrl *Controller) SetAccessCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.Config.EnvConfig.CORS.GlobalDomain
	c.SetCookie("access_token", token, timeExpired, "/", globalDomain, false, true)
}

func (ctrl *Controller) SetRefreshCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.Config.EnvConfig.CORS.GlobalDomain
	c.SetCookie("refresh_token", token, timeExpired, "/", globalDomain, false, true)
}

func isValidLoginRequest(req ClientRequestBasicLogin) bool {
	return req.Password != nil && (req.Username != nil || req.Email != nil || req.Phone != nil)
}

func (ctrl *Controller) HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (ctrl *Controller) AuthenticateUser(req *ClientRequestBasicLogin, c *gin.Context) (*entity.User, error) {
	hashedPassword := ctrl.HashPassword(*req.Password)

	if req.Username != nil {
		return ctrl.Repository.GetUserByIdentifierAndPassword("username", *req.Username, hashedPassword)
	} else if req.Email != nil {
		return ctrl.Repository.GetUserByIdentifierAndPassword("email", *req.Email, hashedPassword)
	} else if req.Phone != nil {
		return ctrl.Repository.GetUserByIdentifierAndPassword("phone", *req.Phone, hashedPassword)
	}
	return nil, fmt.Errorf("missing login identifier")
}

func (ctrl *Controller) GenerateToken() string {
	return uuid.NewString() + uuid.NewString()
}

func (ctrl *Controller) hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (ctrl *Controller) CheckNullString(str *string) string {
	if str == nil || *str == "" {
		return ""
	}
	return *str
}

func (ctrl *Controller) IsValidEmail(email string) bool {
	// Simple regex for email validation
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	at := 0
	for i, char := range email {
		if char == '@' {
			at++
			if at > 1 || i == 0 || i == len(email)-1 {
				return false
			}
		} else if char == '.' && (i == 0 || i == len(email)-1 || email[i-1] == '@') {
			return false
		}
	}
	return at == 1
}

func (ctrl *Controller) IsValidPhone(phone string) bool {
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func handleTokenError(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Refresh token not found"})
		return
	}

	switch err.Error() {
	case "record not found":
		c.JSON(http.StatusNotFound, gin.H{"error": "Refresh token not found or revoked"})
	case "refresh token expired":
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}

// DownloadImageFromURL downloads an image from the given URL and returns the image data and content type
func (ctrl *Controller) DownloadImageFromURL(imageURL string) ([]byte, string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download image: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download image: HTTP %d", resp.StatusCode)
	}

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(fileBytes)
	}

	return fileBytes, contentType, nil
}

// GetFileExtensionFromContentType returns the appropriate file extension for the given content type
func (ctrl *Controller) GetFileExtensionFromContentType(contentType string, fallbackURL string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/webp":
		return "webp"
	case "image/svg+xml":
		return "svg"
	case "image/x-icon", "image/vnd.microsoft.icon":
		return "ico"
	default:
		// Try to get extension from URL or default to jpg
		if fallbackURL != "" {
			if ext := filepath.Ext(fallbackURL); ext != "" {
				return strings.TrimPrefix(ext, ".")
			}
		}
		return "jpg"
	}
}

// GenerateAvatarHash generates a unique hash from username and timestamp
func (ctrl *Controller) GenerateAvatarHash(username string) string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%s_%d", username, timestamp)
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))[:16] // Use first 16 characters for shorter hash
}

// UploadAvatarFromURL downloads an image from URL and uploads it to the upload service
func (ctrl *Controller) UploadAvatarFromURL(username string, imageURL string) (string, error) {
	// Download image
	fileBytes, contentType, err := ctrl.DownloadImageFromURL(imageURL)
	if err != nil {
		return "", err
	}

	// Get file extension
	extension := ctrl.GetFileExtensionFromContentType(contentType, imageURL)

	// Generate filename: userId.{extension}
	filename := fmt.Sprintf("%s.%s", ctrl.GenerateAvatarHash(username), extension)

	// Upload to service
	uploadedURL, err := ctrl.Provider.UploadServiceProvider.UploadAvatarImage(username, fileBytes, filename, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	return uploadedURL, nil
}

// UploadAvatarFromFile processes an uploaded file and uploads it to the upload service
func (ctrl *Controller) UploadAvatarFromFile(username string, fileBytes []byte, contentType string) (string, error) {
	// Get file extension
	extension := ctrl.GetFileExtensionFromContentType(contentType, "")

	// Generate filename: userId.{extension}
	filename := fmt.Sprintf("%s.%s", ctrl.GenerateAvatarHash(username), extension)

	// Upload to service
	uploadedURL, err := ctrl.Provider.UploadServiceProvider.UploadAvatarImage(username, fileBytes, filename, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	return uploadedURL, nil
}

// ValidateImageContentType validates if the content type is allowed for images
func (ctrl *Controller) ValidateImageContentType(contentType string) bool {
	allowedTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
		"image/x-icon",
		"image/vnd.microsoft.icon",
	}

	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}

// ExecuteInTransaction executes a function within a database transaction using GORM's Transaction method
func (ctrl *Controller) ExecuteInTransaction(fn func(tx *gorm.DB) error) error {
	return ctrl.Repository.Db.Transaction(fn)
}

// RemoveVietnameseDiacritics removes Vietnamese diacritics from text using unidecode library
func (ctrl *Controller) RemoveVietnameseDiacritics(text string) string {
	return unidecode.Unidecode(text)
}

// GenerateUsernameFromFullName generates username from fullname: no space, lowercase + diacritic removal
func (ctrl *Controller) GenerateUsernameFromFullName(fullName string) string {
	if fullName == "" {
		return ""
	}

	// Remove Vietnamese diacritics, spaces and convert to lowercase
	normalizedName := ctrl.RemoveVietnameseDiacritics(fullName)
	username := strings.ToLower(strings.ReplaceAll(normalizedName, " ", ""))

	return username
}

// GenerateUsernameFromFullNameWithTransaction generates username with count check within transaction
func (ctrl *Controller) GenerateUsernameFromFullNameWithTransaction(tx *gorm.DB, fullName string) (string, error) {
	if fullName == "" {
		return "", fmt.Errorf("fullname cannot be empty")
	}

	// Remove Vietnamese diacritics, spaces and convert to lowercase
	normalizedName := ctrl.RemoveVietnameseDiacritics(fullName)
	baseUsername := strings.ToLower(strings.ReplaceAll(normalizedName, " ", ""))

	// Get all usernames starting with baseUsername within transaction
	usernames, err := ctrl.Repository.GetUsernamesStartingWithTransaction(tx, baseUsername)
	if err != nil {
		return "", fmt.Errorf("failed to get usernames: %w", err)
	}

	// If no existing usernames, return the base username
	if len(usernames) == 0 {
		return baseUsername, nil
	}

	// Find the maximum suffix number
	maxSuffix := -1
	baseLen := len(baseUsername)

	for _, username := range usernames {
		if username == baseUsername {
			// Exact match, set maxSuffix to 0 if not set
			if maxSuffix < 0 {
				maxSuffix = 0
			}
		} else if len(username) > baseLen && strings.HasPrefix(username, baseUsername) {
			// Extract suffix and check if it's a number
			suffix := username[baseLen:]
			var num int
			if _, err := fmt.Sscanf(suffix, "%d", &num); err == nil {
				if num > maxSuffix {
					maxSuffix = num
				}
			}
		}
	}

	// Generate new username with next suffix
	if maxSuffix < 0 {
		// No matching usernames found (shouldn't happen, but just in case)
		return baseUsername, nil
	}

	return fmt.Sprintf("%s%d", baseUsername, maxSuffix+1), nil
}
