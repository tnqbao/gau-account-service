package auth

import (
	"github.com/tnqbao/gau_user_service/api/user"
	"github.com/tnqbao/gau_user_service/providers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req providers.ClientReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("UserRequest binding error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserRequest binding error: " + err.Error()})
		return
	}
	*req.Password = providers.HashPassword(*req.Password)
	if req.Username == nil || req.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and Password must be provided"})
		return
	}

	log.Println("Parsed Request:", req)

	user.api_user.CreateUser(c, req)
}
