package helpers

import (
	"backend-golang-api/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(config.GetEnv("JWT_SECRET", "secret_key"))

func GenerateToken(userID uint, email string, role string) (string, error) {

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := JwtClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID uint, email string, role string) (string, time.Time, error) {

	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := JwtClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, expirationTime, err
}

func ParseToken(tokenStr string) (*JwtClaims, error) {

	claims := &JwtClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
