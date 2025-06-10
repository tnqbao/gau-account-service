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

	if refreshToken != "" {
		hashedToken := ctrl.hashToken(refreshToken)
		deviceID := c.GetHeader("X-Device-ID")

		// Tìm refresh token trong DB
		refreshTokenRecord, err := ctrl.repository.GetRefreshTokenByTokenAndDevice(hashedToken, deviceID)
		if err != nil {
			log.Println("Error fetching refresh token:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		if refreshTokenRecord == nil {
			log.Printf("No matching refresh token found: hash=%s, deviceID=%s\n", hashedToken, deviceID)
		} else {
			// Xóa trong DB
			rowsAffected, err := ctrl.repository.DeleteRefreshTokenByTokenAndDevice(hashedToken, deviceID)
			if err != nil {
				log.Println("Error deleting refresh token:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}

			if rowsAffected > 0 {
				// Thêm vào blacklist + set TTL bằng thời gian còn lại của token
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
	} else {
		log.Println("No refresh token provided in header or cookie")
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
