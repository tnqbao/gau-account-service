package utils

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-account-service/config"
	"github.com/tnqbao/gau-account-service/repository"
	"strconv"
	"strings"
)

func ExtractToken(c *gin.Context) string {
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

func ParseToken(tokenString string, config *config.EnvConfig) (*jwt.Token, error) {
	secret := []byte(config.JWT.SecretKey)
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
}

func ExtractJID(claims jwt.MapClaims) (int64, error) {
	if val, ok := claims["jti"]; ok {
		return ParseJIDValue(val)
	}
	if val, ok := claims["jid"]; ok {
		return ParseJIDValue(val)
	}
	return 0, errors.New("Token is missing jti/jid")
}

func ParseJIDValue(val interface{}) (int64, error) {
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

func InjectClaimsToContext(c *gin.Context, claims jwt.MapClaims) error {
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

func ValidateToken(ctx context.Context, tokenStr string, config *config.EnvConfig, repo *repository.Repository) (jwt.MapClaims, error) {
	// Parse token
	token, err := ParseToken(tokenStr, config)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check blacklist
	jid, err := ExtractJID(claims)
	if err != nil {
		return nil, err
	}
	revoked, err := repo.GetBit(ctx, "blacklist_bitmap", jid)
	if err != nil {
		return nil, errors.New("redis error")
	}
	if revoked == 1 {
		return nil, errors.New("token has been revoked")
	}

	return claims, nil
}
