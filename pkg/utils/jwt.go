package utils

import (
	"fmt"
	"time"

	"github.com/ducklawrence05/go-test-backend-api/pkg/setting"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateJWT(secret []byte, userID uuid.UUID, expireInSeconds int) (string, error) {
	expiration := time.Now().Add(time.Second * time.Duration(expireInSeconds))

	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ValidateToken(secret []byte, tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// GenerateAcAndRtTokens concurrently creates access token and refresh token
func GenerateAcAndRtTokens(config *setting.Config, userID uuid.UUID) (string, string, error) {
	accessToken, err := CreateJWT([]byte(config.JWT.AccessTokenKey), userID, config.JWT.AccessTokenExpiresIn)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := CreateJWT([]byte(config.JWT.RefreshTokenKey), userID, config.JWT.RefreshTokenExpiresIn)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
