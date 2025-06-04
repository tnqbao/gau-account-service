package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/models"
	"github.com/tnqbao/gau-account-service/providers"
	"github.com/tnqbao/gau-account-service/repositories"
	"net/http"
)

func (ctrl *Controller) GetAccountInfo(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(400, gin.H{"error": "User ID is required"})
		return
	}

	// Convert the user ID to a UUID
	userId, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid User ID format"})
		return
	}

	// Fetch user information from the repository
	userInfo, err := repositories.GetUserById(userId, c)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(404, gin.H{"error": "User not found"})
		} else {
			c.JSON(500, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	UserInfoResponse := providers.UserInfoResponse{
		UserId:      userInfo.UserID,
		FullName:    *userInfo.FullName,
		Email:       *userInfo.Email,
		Phone:       *userInfo.Phone,
		DateOfBirth: userInfo.DateOfBirth,
	}
	c.JSON(200, gin.H{
		"user_info": UserInfoResponse,
	})
}

func (ctrl *Controller) UpdateAccountInfo(c *gin.Context) {
	id := c.Param("user_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

	var req providers.UserInformationUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	user, err := repositories.GetUserById(userID, c)
	if err != nil {
		status := http.StatusInternalServerError
		errMsg := "Internal server error: " + err.Error()
		if err.Error() == "record not found" {
			status = http.StatusNotFound
			errMsg = "User not found"
		}
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	updateData := &models.User{
		UserID:      user.UserID,
		Username:    req.Username,
		FullName:    providers.CheckNullString(req.FullName),
		Email:       providers.CheckNullString(req.Email),
		Phone:       providers.CheckNullString(req.Phone),
		DateOfBirth: req.DateOfBirth,
		Gender:      req.Gender,
		FacebookURL: providers.CheckNullString(req.FacebookURL),
		GithubURL:   providers.CheckNullString(req.GitHubURL),
	}

	updatedUser, err := repositories.UpdateUser(updateData, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User information updated successfully", "user_info": updatedUser})
}

//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/tnqbao/gau-account-service/models"
//	"github.com/tnqbao/gau-account-service/providers"
//	"gorm.io/gorm"
//	"log"
//	"net/http"
//	"strconv"
//)
//
////func CreateUser(c *gin.Context, r providers.ClientReq) {
////	db := c.MustGet("db").(*gorm.DB)
////	var userName *string
////	if r.Username != nil && *r.Username != "" {
////		userName = r.Username
////	} else if r.Email != nil && *r.Email != "" {
////		userName = r.Email
////	}
////
////	err := db.Transaction(func(tx *gorm.DB) error {
////		userAuth := models.{
////			Username:   userName,
////			Password:   r.Password,
////			Permission: "member",
////
////		}
////		if err := tx.Create(&userAuth).Error; err != nil {
////			return err
////		}
////
////		userInfor := models.UserInformation{
////			FullName:    r.Fullname,
////			Email:       providers.CheckNullString(r.Email),
////			Phone:       providers.CheckNullString(r.Phone),
////			DateOfBirth: r.DateOfBirth,
////			UserId:      userAuth.UserId,
////		}
////
////		if err := tx.Create(&userInfor).Error; err != nil {
////			return err
////		}
////		return nil
////	})
////	if err != nil {
////		log.Println("Transaction error:", err)
////		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create user1: " + err.Error()})
////		return
////	}
////	c.JSON(http.StatusOK, gin.H{"message": "User successfully created"})
////}
//
//func GetUserById(c *gin.Context) {
//	db := c.MustGet("db").(*gorm.DB)
//	idStr := c.Param("id")
//	id, err := strconv.ParseUint(idStr, 10, 32)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user1 ID format"})
//		return
//	}
//
//	tokenId, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
//		return
//	}
//
//	tokenIdUint, ok := tokenId.(uint)
//	if !ok {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id format"})
//		return
//	}
//
//	permission, exists := c.Get("permission")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "permission not found"})
//		return
//	}
//
//	if uint64(tokenIdUint) != id && permission != "admin" {
//		c.JSON(http.StatusForbidden, gin.H{"status": http.StatusForbidden, "error": "You don't have permission to access this resource!"})
//		return
//	}
//	var user models.UserAuthentication
//	var userInfo models.UserInformation
//	err = db.Transaction(func(tx *gorm.DB) error {
//		if err := tx.First(&user, "user_id = ?", id).Error; err != nil {
//			return err
//		}
//		if err := tx.First(&userInfo, "user_id = ?", id).Error; err != nil {
//			return err
//		}
//		return nil
//	})
//
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		} else {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		}
//		return
//	}
//	response := providers.ServerResp{
//		Fullame: providers.ToString(userInfo.FullName),
//	}
//	c.JSON(http.StatusOK, response)
//}
//
//func GetMe(c *gin.Context) {
//	db := c.MustGet("db").(*gorm.DB)
//	tokenId, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
//		return
//	}
//
//	var user models.UserAuthentication
//	var userInfo models.UserInformation
//	err := db.Transaction(func(tx *gorm.DB) error {
//		if err := tx.First(&user, "user_id = ?", tokenId).Error; err != nil {
//			return err
//		}
//		if err := tx.First(&userInfo, "user_id = ?", tokenId).Error; err != nil {
//			return err
//		}
//		return nil
//	})
//
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		} else {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		}
//		return
//	}
//	response := providers.ServerResp{
//		UserId:     user.UserId,
//		Fullame:    providers.ToString(userInfo.FullName),
//		Email:      providers.ToString(userInfo.Email),
//		Phone:      providers.ToString(userInfo.Phone),
//		DateBirth:  userInfo.DateOfBirth,
//		Permission: user.Permission,
//	}
//	c.JSON(http.StatusOK, response)
//}
//
//func DeleteUserById(c *gin.Context) {
//	db := c.MustGet("db").(*gorm.DB)
//	tokenId, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
//		return
//	}
//
//	err := db.Transaction(func(tx *gorm.DB) error {
//		if result := tx.Delete(&models.UserAuthentication{}, tokenId); result.Error != nil {
//			return result.Error
//		} else if result.RowsAffected == 0 {
//			return gorm.ErrRecordNotFound
//		}
//		if err := tx.Delete(&models.UserInformation{}, tokenId).Error; err != nil {
//			return err
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		} else {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		}
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
//}
//
//func UpdateUserInformation(c *gin.Context) {
//	db := c.MustGet("db").(*gorm.DB)
//
//	tokenId, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
//		return
//	}
//
//	tokenIdUint, ok := tokenId.(uint)
//	if !ok {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id format"})
//		return
//	}
//
//	userUpate := models.UserInformation{}
//	if err := c.ShouldBindJSON(&userUpate); err != nil {
//		log.Println("UserRequest binding error:", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "UserRequest binding error: " + err.Error()})
//		return
//	}
//
//	var userInfor models.UserInformation
//
//	err := db.Transaction(func(tx *gorm.DB) error {
//		if err := tx.First(&userInfor, "user_id = ?", tokenIdUint).Error; err != nil {
//			return err
//		}
//		db.Model(&userInfor).Updates(userUpate)
//		return nil
//	})
//
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		} else {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		}
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{"message": "Update successful"})
//}
