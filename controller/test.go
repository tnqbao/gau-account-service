package controller

import (
	"github.com/gin-gonic/gin"
	"time"
)

func (ctrl *Controller) TestDeployment(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Checked Deployment: " + time.Now().String(),
	})
}

func (ctrl *Controller) CheckHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Checked Health: " + time.Now().String(),
	})
}
