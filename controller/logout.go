package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ctrl *Controller) Logout(c *gin.Context) {
	c.MustGet("auth_token")
	// add token to blacklist

	c.SetCookie("auth_token", "", -1, "/", ctrl.config.JWT.SecretKey, false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
