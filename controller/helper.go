package controller

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/models"
	"github.com/tnqbao/gau-account-service/providers"
	"github.com/tnqbao/gau-account-service/repositories"
	"time"
)

func (ctrl *Controller) CreateAuthToken(claims ClaimsResponse) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"permission": claims.UserPermission,
		"fullname":   claims.FullName,
		"exp":        claims.ExpiresAt.Unix(),
		"iat":        time.Now().Unix(),
	})

	return token.SignedString([]byte(ctrl.config.JWT.SecretKey))
}

func (ctrl *Controller) SetAuthCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.config.CORS.GlobalDomain
	c.SetCookie("auth_token", token, timeExpired, "/", globalDomain, false, true)
}

func (ctrl *Controller) SetRefreshCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.config.CORS.GlobalDomain
	c.SetCookie("refresh_token", token, timeExpired, "/", globalDomain, false, true)
}

func isValidLoginRequest(req providers.ClientRequestLogin) bool {
	return req.Password != nil && (req.Username != nil || req.Email != nil || req.Phone != nil)
}

func (ctrl *Controller) AuthenticateUser(req *providers.ClientRequestLogin, c *gin.Context) (*models.User, error) {
	hashedPassword := providers.HashPassword(*req.Password)

	if req.Username != nil {
		return repositories.GetUserByIdentifierAndPassword("username", *req.Username, hashedPassword, c)
	} else if req.Email != nil {
		return repositories.GetUserByIdentifierAndPassword("email", *req.Email, hashedPassword, c)
	} else if req.Phone != nil {
		return repositories.GetUserByIdentifierAndPassword("phone", *req.Phone, hashedPassword, c)
	}
	return nil, fmt.Errorf("missing login identifier")
}

func (ctrl *Controller) GenerateAccessToken(userID string) (string, error) {
	// sử dụng jwt-go
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ctrl.config.JWT.SecretKey))
}

func (ctrl *Controller) GenerateRefreshToken() string {
	return uuid.NewString() + uuid.NewString()
}

func (ctrl *Controller) hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum(nil))
}
