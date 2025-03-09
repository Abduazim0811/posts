package jwt

import (
	"time"
	"user-service/internal/logger"

	"github.com/golang-jwt/jwt/v5"
)

var (
	SecretKey = "abduazim"      
	ExpiresIn = 24*time.Hour
)

type Claims struct {
	Email string 
	jwt.RegisteredClaims
}

func GenerateToken(email string) (string, error) {
	logger.Logger.Printf("GenerateToken boshlandi: email=%s", email)

	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		logger.Logger.Printf("GenerateToken: Token imzolashda xato: %v", err)
		return "", err
	}

	logger.Logger.Printf("GenerateToken muvaffaqiyatli yakunlandi: email=%s", email)
	return signedToken, nil
}