package api_user

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", os.Getenv("GLOBAL_DOMAIN"), false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout succesful"})
}
