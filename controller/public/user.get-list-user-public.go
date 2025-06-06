package public

import (
	"github.com/gin-gonic/gin"
	provider "github.com/tnqbao/gau_user_service/providers"
	"gorm.io/gorm"
	"net/http"
)

func GetListUserPublicByIDs(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var req provider.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	idArr := req.IDs
	if len(idArr) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user IDs provided"})
		return
	}

	var userPublics []provider.UserPublic
	query := `
		SELECT user_id, full_name AS fullname
		FROM user_informations
		WHERE user_id = ANY($1)
	`
	err := db.Raw(query, idArr).Scan(&userPublics).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(userPublics) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userMap := make(map[uint]provider.UserPublic)
	for _, user := range userPublics {
		userMap[user.UserId] = user
	}

	var response []provider.UserPublic
	for _, id := range idArr {
		if user, exists := userMap[id]; exists {
			response = append(response, user)
		}
	}

	c.JSON(http.StatusOK, response)
}
