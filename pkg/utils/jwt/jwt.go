package jwt

import (
	"time"

	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/errorcode"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type EmailClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateJWT(secret []byte, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ValidateToken[T jwt.Claims](secret []byte, tokenString string) (*T, error) {
	var claims T
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorcode.ErrUnexpectedSigningToken
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(T)
	if !ok || !token.Valid {
		return nil, errorcode.ErrInvalidToken
	}
	return &claims, nil
}

// GenerateAcAndRtTokens concurrently creates access token and refresh token
func GenerateAcAndRtTokens(config *config.Config, userID uuid.UUID) (string, string, error) {
	accessToken, err := CreateJWT([]byte(config.JWT.AccessTokenKey), UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWT.AccessTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		return "", "", err
	}

	refreshToken, err := CreateJWT([]byte(config.JWT.RefreshTokenKey), UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWT.RefreshTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}


