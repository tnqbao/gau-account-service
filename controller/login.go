package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnqbao/gau-account-service/providers"
	"log"
	"net/http"
	"time"
)

func (ctrl *Controller) LoginWithIdentifierAndPassword(c *gin.Context) {
	var req providers.ClientRequestLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	if !isValidLoginRequest(req) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email/Username/Phone and Password are required"})
		return
	}

	user, err := ctrl.AuthenticateUser(&req, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var expirationTime time.Time
	if req.KeepLogin != nil && *req.KeepLogin == "true" {
		expirationTime = time.Now().Add(7 * 24 * time.Hour)
	} else {
		expirationTime = time.Now().Add(15 * time.Minute)
	}

	claims := &ClaimsResponse{
		UserID:         user.UserId,
		FullName:       *user.FullName,
		UserPermission: user.Permission,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token, err := ctrl.CreateAuthToken(*claims)
	if err != nil {
		log.Println("Error creating auth token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create auth token: " + err.Error()})
		return
	}

	var timeExpired int
	if req.KeepLogin != nil && *req.KeepLogin == "true" {
		timeExpired = 3600 * 24 * 30
	} else {
		timeExpired = 0
	}

	ctrl.SetAuthCookie(c, token, timeExpired)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
