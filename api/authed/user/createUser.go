package api_user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	provider "github.com/tnqbao/gau_services/api"
	"github.com/tnqbao/gau_services/models"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context, r provider.ClientReq) {
	db := c.MustGet("db").(*gorm.DB)
	err := db.Transaction(func(tx *gorm.DB) error {
		userAuth := models.UserAuthentication{
			Username:   r.Username,
			Password:   r.Password,
			Permission: "member",
		}
		if err := tx.Create(&userAuth).Error; err != nil {
			return err
		}

		userInfor := models.UserInformation{
			FullName:    r.Fullname,
			Email:       provider.CheckNullString(r.Email),
			Phone:       provider.CheckNullString(r.Phone),
			DateOfBirth: provider.FormatStringToDate(r.DateOfBirth),
			UserId:      userAuth.UserId,
		}

		if err := tx.Create(&userInfor).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("Transaction error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create user: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User successfully created"})
}
