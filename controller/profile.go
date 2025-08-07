package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/entity"
	"github.com/tnqbao/gau-account-service/utils"
	"io"
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

	file, err := c.FormFile("avatar_image")
	if err != nil {
		utils.JSON400(c, "File not found or invalid")
		return
	}

	if file.Size > 50*1024*1024 {
		utils.JSON400(c, "File size exceeds the limit of 50MB")
		return
	}

	// Validate image type: allow jpeg, jpg, png, gif, webp
	contentType := file.Header.Get("Content-Type")
	var ext string
	switch contentType {
	case "image/jpeg":
		ext = "jpg"
	case "image/jpg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/gif":
		ext = "gif"
	case "image/webp":
		ext = "webp"
	default:
		utils.JSON400(c, "Invalid file type. Only JPEG, JPG, PNG, GIF, and WEBP are allowed")
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

	filename := userID.String() + "." + ext
	imageURL, err := ctrl.Provider.UploadServiceProvider.UploadAvatarImage(userID.String(), fileBytes, filename)
	if err != nil {
		utils.JSON500(c, "Failed to upload image: "+err.Error())
		return
	}

	user, err := ctrl.Repository.GetUserById(userID)
	if err != nil {
		if err.Error() == "record not found" {
			utils.JSON404(c, "User not found")
			return
		}
		utils.JSON500(c, "Internal server error")
		return
	}

	user.AvatarURL = &imageURL

	if _, err := ctrl.Repository.UpdateUser(user); err != nil {
		utils.JSON500(c, "Failed to update user information: "+err.Error())
		return
	}

	utils.JSON200(c, gin.H{
		"message":    "Avatar image updated successfully",
		"avatar_url": ctrl.Config.EnvConfig.ExternalService.CDNServiceURL + "/" + imageURL,
	})
}
