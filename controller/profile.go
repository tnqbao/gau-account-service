package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/schemas"
	"github.com/tnqbao/gau-account-service/utils"
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

	updateData := &schemas.User{
		UserID:      user.UserID,
		Username:    utils.Coalesce(req.Username, user.Username),
		FullName:    utils.Coalesce(req.FullName, user.FullName),
		Email:       utils.Coalesce(req.Email, user.Email),
		Phone:       utils.Coalesce(req.Phone, user.Phone),
		DateOfBirth: utils.Coalesce(req.DateOfBirth, user.DateOfBirth),
		Gender:      utils.Coalesce(req.Gender, user.Gender),
		FacebookURL: utils.Coalesce(req.FacebookURL, user.FacebookURL),
		GithubURL:   utils.Coalesce(req.GitHubURL, user.GithubURL),
		ImageURL:    user.ImageURL,
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
