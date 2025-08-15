package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/entity"
	"github.com/tnqbao/gau-account-service/utils"
	"gorm.io/gorm"
	"io"
	"net/http"
)

func (ctrl *Controller) GetAccountInfo(c *gin.Context) {
	userId := c.MustGet("user_id")
	if userId == nil {
		utils.JSON400(c, "User ID is required")
		return
	}

	var uuidUserId uuid.UUID
	switch v := userId.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			utils.JSON400(c, "Invalid User ID format")
			return
		}
		uuidUserId = parsed
	case uuid.UUID:
		uuidUserId = v
	default:
		utils.JSON400(c, "Invalid User ID type")
		return
	}

	userInfo, err := ctrl.Repository.GetUserById(uuidUserId)
	if err != nil {
		if err.Error() == "record not found" {
			utils.JSON404(c, "User not found")
		} else {
			utils.JSON500(c, "Internal server error")
		}
		return
	}

	UserInfoResponse := UserInfoResponse{
		UserId:          userInfo.UserID,
		FullName:        ctrl.CheckNullString(userInfo.FullName),
		Email:           ctrl.CheckNullString(userInfo.Email),
		Phone:           ctrl.CheckNullString(userInfo.Phone),
		GithubUrl:       ctrl.CheckNullString(userInfo.GithubURL),
		FacebookUrl:     ctrl.CheckNullString(userInfo.FacebookURL),
		AvatarURL:       ctrl.CheckNullString(userInfo.AvatarURL),
		IsEmailVerified: userInfo.IsEmailVerified,
		IsPhoneVerified: userInfo.IsPhoneVerified,
		DateOfBirth:     userInfo.DateOfBirth,
	}
	utils.JSON200(c, gin.H{
		"user_info": UserInfoResponse,
	})
}

func (ctrl *Controller) UpdateAccountInfo(c *gin.Context) {
	userIdRaw := c.MustGet("user_id")
	if userIdRaw == nil {
		utils.JSON400(c, "User ID is required")
		return
	}

	var userID uuid.UUID
	switch v := userIdRaw.(type) {
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			utils.JSON400(c, "Invalid User ID format")
			return
		}
		userID = id
	case uuid.UUID:
		userID = v
	default:
		utils.JSON400(c, "Invalid User ID type")
		return
	}

	user, err := ctrl.Repository.GetUserById(userID)
	if err != nil {
		if err.Error() == "record not found" {
			utils.JSON404(c, "User not found")
		} else {
			utils.JSON500(c, "Internal server error")
		}
		return
	}

	var req UserInformationUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	updateData := &entity.User{
		UserID:      user.UserID,
		Username:    utils.Coalesce(req.Username, user.Username),
		FullName:    utils.Coalesce(req.FullName, user.FullName),
		Email:       utils.Coalesce(req.Email, user.Email),
		Phone:       utils.Coalesce(req.Phone, user.Phone),
		DateOfBirth: utils.Coalesce(req.DateOfBirth, user.DateOfBirth),
		Gender:      utils.Coalesce(req.Gender, user.Gender),
		FacebookURL: utils.Coalesce(req.FacebookURL, user.FacebookURL),
		GithubURL:   utils.Coalesce(req.GitHubURL, user.GithubURL),
		AvatarURL:   user.AvatarURL,
		Permission:  user.Permission,
	}

	// Gọi hàm cập nhật DB
	updatedUser, err := ctrl.Repository.UpdateUser(updateData)
	if err != nil {
		utils.JSON500(c, "Internal server error")
		return
	}

	utils.JSON200(c, gin.H{
		"message":   "User information updated successfully",
		"user_info": updatedUser,
	})
}

func (ctrl *Controller) UpdateAvatarImage(c *gin.Context) {
	userIdRaw := c.MustGet("user_id")
	if userIdRaw == nil {
		utils.JSON400(c, "User ID is required")
		return
	}

	var userID uuid.UUID
	switch v := userIdRaw.(type) {
	case string:
		id, err := uuid.Parse(v)
		if err != nil {
			utils.JSON400(c, "Invalid User ID format")
			return
		}
		userID = id
	case uuid.UUID:
		userID = v
	default:
		utils.JSON400(c, "Invalid User ID type")
		return
	}

	// Change from "avatar_image" to "file" to match the upload service expectation
	file, err := c.FormFile("file")
	if err != nil {
		utils.JSON400(c, "File not found or invalid")
		return
	}

	if file.Size > 50*1024*1024 {
		utils.JSON400(c, "File size exceeds the limit of 50MB")
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		utils.JSON500(c, "Failed to open uploaded file: "+err.Error())
		return
	}
	defer openedFile.Close()

	fileBytes, err := io.ReadAll(openedFile)
	if err != nil {
		utils.JSON500(c, "Failed to read uploaded file: "+err.Error())
		return
	}

	// Get content type from the uploaded file and detect if needed
	contentType := file.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		contentType = http.DetectContentType(fileBytes)
	}

	// Validate image type using helper function
	if !ctrl.ValidateImageContentType(contentType) {
		utils.JSON400(c, "Invalid file type. Only JPEG, JPG, PNG, GIF, WEBP, SVG, and ICO are allowed")
		return
	}

	// Use GORM's Transaction method for avatar upload and database update
	var fullImageURL string
	err = ctrl.ExecuteInTransaction(func(tx *gorm.DB) error {
		// Upload avatar image
		imageURL, err := ctrl.UploadAvatarFromFile(userID, fileBytes, contentType)
		if err != nil {
			return fmt.Errorf("failed to upload image: %w", err)
		}

		// Add CDN URL prefix
		fullImageURL = fmt.Sprintf("%s/images/%s", ctrl.Config.EnvConfig.ExternalService.CDNServiceURL, imageURL)

		// Get user and update avatar URL
		user, err := ctrl.Repository.GetUserById(userID)
		if err != nil {
			if err.Error() == "record not found" {
				return fmt.Errorf("user not found")
			}
			return fmt.Errorf("failed to get user: %w", err)
		}

		user.AvatarURL = &fullImageURL

		// Update user within transaction using repository method
		if _, err := ctrl.Repository.UpdateUserWithTransaction(tx, user); err != nil {
			return fmt.Errorf("failed to update user information: %w", err)
		}

		return nil
	})

	if err != nil {
		utils.JSON500(c, err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message":    "Avatar image updated successfully",
		"avatar_url": fullImageURL,
	})
}
