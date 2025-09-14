package auth

import (
	"errors"
	"fmt"
	"time"
	"net/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	token, err := t.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	t, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !t.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subject UUID: %w", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix){
		return "", errors.New("invalid Authorization header format")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader,prefix))
	if token == ""{
		return "", errors.New("empty bearer token")
	}
	return token, nil
}
