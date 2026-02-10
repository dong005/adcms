package utils

import (
	"adcms/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	TenantID uint   `json:"tenant_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TempClaims struct {
	UserID     uint   `json:"user_id"`
	TenantID   uint   `json:"tenant_id"`
	Username   string `json:"username"`
	RequireOTP bool   `json:"require_otp"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, tenantID uint, username string) (string, error) {
	cfg := config.GlobalConfig.JWT
	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func GenerateTempToken(userID, tenantID uint, username string) (string, error) {
	cfg := config.GlobalConfig.JWT
	claims := TempClaims{
		UserID:     userID,
		TenantID:   tenantID,
		Username:   username,
		RequireOTP: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GlobalConfig.JWT

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ParseTempToken(tokenString string) (*TempClaims, error) {
	cfg := config.GlobalConfig.JWT

	token, err := jwt.ParseWithClaims(tokenString, &TempClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TempClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	return GenerateToken(claims.UserID, claims.TenantID, claims.Username)
}
