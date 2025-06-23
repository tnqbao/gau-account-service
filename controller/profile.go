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
	id := c.Param("user_id")
	if id == "" {
		utils.JSON400(c, "User ID is required")
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		utils.JSON400(c, "Invalid User ID format")
		return
	}

	var req UserInformationUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSON400(c, "Invalid request format: "+err.Error())
		return
	}

	user, err := ctrl.Repository.GetUserById(userID)
	if err != nil {
		if err.Error() == "record not found" {
			utils.JSON404(c, "User not found")
		} else {
			utils.JSON500(c, "Internal server error: "+err.Error())
		}
		return
	}

	updateData := &schemas.User{
		UserID:      user.UserID,
		Username:    req.Username,
		FullName:    req.FullName,
		Email:       req.Email,
		Phone:       req.Phone,
		DateOfBirth: req.DateOfBirth,
		Gender:      req.Gender,
		FacebookURL: req.FacebookURL,
		GithubURL:   req.GitHubURL,
	}

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
