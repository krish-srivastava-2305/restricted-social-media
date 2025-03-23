package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret")

func GenerateToken(payload string) (string, error) {
	claims := jwt.MapClaims{
		"email": payload,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("Error generating token: %v", err.Error())
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, fmt.Errorf("Error validating token %v", err.Error())
	}

	if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, fmt.Errorf("Malformed Token")
	}

	if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, fmt.Errorf("Invalid Signature")
	}

	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, fmt.Errorf("Token Expired")
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is invalid")
	}

	return token, nil
}
