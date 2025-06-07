package controller

import "github.com/gin-gonic/gin"

func (ctrl *Controller) LoginWithFacebook(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Login with Facebook is not implemented yet",
	})
}

func (ctrl *Controller) LoginWithGoogle(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Login with Google is not implemented yet",
	})
}
