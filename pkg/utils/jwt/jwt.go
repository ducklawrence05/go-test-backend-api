package jwt

import (
	"time"

	"github.com/ducklawrence05/go-test-backend-api/config"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/errorcode"
	"github.com/ducklawrence05/go-test-backend-api/internal/constants/jwtpurpose"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type CustomClaims struct {
	Purpose jwtpurpose.JWTPurpose `json:"purpose"`
	jwt.RegisteredClaims
}

func createJWT(secret []byte, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ValidateToken(secret []byte, tokenString string, purpose jwtpurpose.JWTPurpose) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorcode.ErrUnexpectedSigningToken
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errorcode.ErrInvalidToken
	}

	if claims.Purpose != purpose {
		return nil, errorcode.ErrInvalidJWTPurpose
	}

	return claims, nil
}

// GenerateAcAndRtTokens concurrently creates access token and refresh token
func GenerateAcAndRtTokens(config *config.Config, userID uuid.UUID) (string, string, error) {
	accessToken, err := createJWT([]byte(config.JWT.AccessTokenKey), CustomClaims{
		Purpose: jwtpurpose.JWTAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWT.AccessTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		return "", "", err
	}

	refreshToken, err := createJWT([]byte(config.JWT.RefreshTokenKey), CustomClaims{
		Purpose: jwtpurpose.JWTRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.JWT.RefreshTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func GenerateEmailToken(
	config *config.Config, email string,
	secret []byte, expiresIn time.Duration, purpose jwtpurpose.JWTPurpose,
) (string, error) {
	emailVerifyToken, err := createJWT(
		secret,
		CustomClaims{
			Purpose: purpose,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   email,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	)
	if err != nil {
		return "", err
	}
	return emailVerifyToken, nil
}
