package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/service"
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware(cfg *config.EnvConfig, svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			abortUnauthorized(c, "Authorization token is required")
			return
		}

		token, err := parseToken(tokenStr, cfg)
		if err != nil || !token.Valid {
			abortUnauthorized(c, "Invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			abortUnauthorized(c, "Invalid token claims")
			return
		}

		// Extract JID (JWT ID)
		jid, err := extractJID(claims)
		if err != nil {
			abortUnauthorized(c, err.Error())
			return
		}

		// Check Redis blacklist
		ctx := c.Request.Context()
		revoked, err := svc.Redis.GetBit(ctx, "blacklist_bitmap", jid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			c.Abort()
			return
		}
		if revoked == 1 {
			abortUnauthorized(c, "Token has been revoked")
			return
		}

		// inject claims into context
		if err := injectClaimsToContext(c, claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if token, err := c.Cookie("access_token"); err == nil && token != "" {
		return token
	}
	authHeader := c.GetHeader("Authorization")
	parts := strings.Fields(authHeader)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}

func parseToken(tokenString string, cfg *config.EnvConfig) (*jwt.Token, error) {
	secret := []byte(cfg.JWT.SecretKey)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
}

func extractJID(claims jwt.MapClaims) (int64, error) {
	if val, ok := claims["jti"]; ok {
		return parseJIDValue(val)
	}
	if val, ok := claims["jid"]; ok {
		return parseJIDValue(val)
	}
	return 0, errors.New("Token is missing jti/jid")
}

func parseJIDValue(val interface{}) (int64, error) {
	switch v := val.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, errors.New("Invalid jid format")
	}
}

func injectClaimsToContext(c *gin.Context, claims jwt.MapClaims) error {
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return errors.New("Invalid user_id format")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("Invalid user_id format")
	}
	c.Set("user_id", userID)

	if permission, ok := claims["permission"].(string); ok {
		c.Set("permission", permission)
	} else {
		c.Set("permission", "")
	}
	return nil
}

func abortUnauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
	c.Abort()
}
