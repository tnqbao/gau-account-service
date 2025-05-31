package providers

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func createAuthToken(claims ClaimsResponse) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"permission": claims.UserPermission,
		"fullname":   claims.FullName,
		"exp":        claims.ExpiresAt.Unix(),
		"iat":        claims.IssuedAt.Unix(),
	})

	secretKey := os.Getenv("SECRET_KEY")
	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
