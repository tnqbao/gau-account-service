package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/schemas"
	"net/http"
	"time"
)

func (ctrl *Controller) CreateAccessToken(claims ClaimsToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jid":        claims.JID,
		"user_id":    claims.UserID,
		"permission": claims.UserPermission,
		"fullname":   claims.FullName,
		"exp":        claims.ExpiresAt.Unix(),
		"iat":        time.Now().Unix(),
	})

	return token.SignedString([]byte(ctrl.config.EnvConfig.JWT.SecretKey))
}

func (ctrl *Controller) SetAccessCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.config.EnvConfig.CORS.GlobalDomain
	c.SetCookie("access_token", token, timeExpired, "/", globalDomain, false, true)
}

func (ctrl *Controller) SetRefreshCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.config.EnvConfig.CORS.GlobalDomain
	c.SetCookie("refresh_token", token, timeExpired, "/", globalDomain, false, true)
}

func isValidLoginRequest(req ClientRequestBasicLogin) bool {
	return req.Password != nil && (req.Username != nil || req.Email != nil || req.Phone != nil)
}

func (ctrl *Controller) HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (ctrl *Controller) AuthenticateUser(req *ClientRequestBasicLogin, c *gin.Context) (*schemas.User, error) {
	hashedPassword := ctrl.HashPassword(*req.Password)

	if req.Username != nil {
		return ctrl.repository.GetUserByIdentifierAndPassword("username", *req.Username, hashedPassword)
	} else if req.Email != nil {
		return ctrl.repository.GetUserByIdentifierAndPassword("email", *req.Email, hashedPassword)
	} else if req.Phone != nil {
		return ctrl.repository.GetUserByIdentifierAndPassword("phone", *req.Phone, hashedPassword)
	}
	return nil, fmt.Errorf("missing login identifier")
}

func (ctrl *Controller) GenerateToken() string {
	return uuid.NewString() + uuid.NewString()
}

func (ctrl *Controller) hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (ctrl *Controller) CheckNullString(str *string) string {
	if str == nil || *str == "" {
		return ""
	}
	return *str
}

func (ctrl *Controller) IsValidEmail(email string) bool {
	// Simple regex for email validation
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	at := 0
	for i, char := range email {
		if char == '@' {
			at++
			if at > 1 || i == 0 || i == len(email)-1 {
				return false
			}
		} else if char == '.' && (i == 0 || i == len(email)-1 || email[i-1] == '@') {
			return false
		}
	}
	return at == 1
}

func (ctrl *Controller) IsValidPhone(phone string) bool {
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func handleTokenError(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Refresh token not found"})
		return
	}

	switch err.Error() {
	case "record not found":
		c.JSON(http.StatusNotFound, gin.H{"error": "Refresh token not found or revoked"})
	case "refresh token expired":
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
