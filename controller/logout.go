package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (ctrl *Controller) Logout(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		refreshToken, _ = c.Cookie("refresh_token")
	}

	if refreshToken == "" {
		log.Println("No refresh token provided in header or cookie")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token provided"})
		c.Abort()
		return
	}

	hashedToken := ctrl.hashToken(refreshToken)
	deviceID := c.GetHeader("X-Device-ID")

	refreshTokenRecord, err := ctrl.repository.GetRefreshTokenByTokenAndDevice(hashedToken, deviceID)
	if err != nil {
		log.Println("Error fetching refresh token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		c.Abort()
		return
	}

	if refreshTokenRecord != nil {
		rowsAffected, err := ctrl.repository.DeleteRefreshTokenByTokenAndDevice(hashedToken, deviceID)
		if err != nil {
			log.Println("Error deleting refresh token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if rowsAffected > 0 {
			ttl := time.Until(refreshTokenRecord.ExpiresAt)
			if ttl > 0 {
				if err := ctrl.service.Redis.ReleaseAndBlacklistIDWithTTL(
					c.Request.Context(),
					refreshTokenRecord.ID,
					ttl,
				); err != nil {
					log.Println("Failed to blacklist refresh token ID with TTL:", err)
				} else {
					log.Printf("Refresh token ID %d blacklisted for %s\n", refreshTokenRecord.ID, ttl)
				}
			}
		}
	}

	// Clear cookies
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	c.Abort()
}
