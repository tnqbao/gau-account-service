package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (ctrl *Controller) Logout(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}

	if refreshToken != "" {
		hashedToken := ctrl.hashToken(refreshToken)
		deviceID := c.GetHeader("X-Device-ID")

		// Gọi hàm Delete, giờ trả về cả RowsAffected
		rowsAffected, err := ctrl.repository.DeleteRefreshTokenByTokenAndDevice(hashedToken, deviceID)
		if err != nil {
			log.Println("Error deleting refresh token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		if rowsAffected == 0 {
			log.Printf("No refresh token found for hash: %s and deviceID: %s\n", hashedToken, deviceID)
		} else {
			log.Printf("Refresh token deleted: hash=%s, deviceID=%s\n", hashedToken, deviceID)
		}
	} else {
		log.Println("No refresh token provided in header or cookie")
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
