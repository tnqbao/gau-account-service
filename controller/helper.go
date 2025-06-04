package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tnqbao/gau-account-service/models"
	"github.com/tnqbao/gau-account-service/providers"
	"github.com/tnqbao/gau-account-service/repositories"
)

func (ctrl *Controller) CreateAuthToken(claims ClaimsResponse) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"permission": claims.UserPermission,
		"fullname":   claims.FullName,
		"exp":        claims.ExpiresAt.Unix(),
		"iat":        claims.IssuedAt.Unix(),
	})

	secretKey := ctrl.config.JWT.SecretKey
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (ctrl *Controller) SetAuthCookie(c *gin.Context, token string, timeExpired int) {
	globalDomain := ctrl.config.CORS.GlobalDomain
	c.SetCookie("auth_token", token, timeExpired, "/", globalDomain, false, true)
	c.Next()
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
