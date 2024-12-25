package api_user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	provider "github.com/tnqbao/gau_services/api"
	api_user_authed "github.com/tnqbao/gau_services/api/authed/user"
)

func Register(c *gin.Context) {
	var req provider.ClientReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserRequest binding error: " + err.Error()})
		return
	}
	*req.Password = provider.HashPassword(*req.Password)
	if req.Username == nil || req.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and Password must be provided"})
		return
	}

	log.Println("Parsed Request:", req)

	api_user_authed.CreateUser(c, req)
}
