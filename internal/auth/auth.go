package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
}

func New(secret string) (*Service, error) {
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is unset or empty")
	}
	return &Service{secret: []byte(secret)}, nil
}

func (s *Service) GenerateJWT(username string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"is-admin": isAdmin,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *Service) ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
}
