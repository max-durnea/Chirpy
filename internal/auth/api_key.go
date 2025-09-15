package auth

import(
	"strings"
	"net/http"
	"errors"
)

func GetAPIKey(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}
	const prefix = "ApiKey "
	if !strings.HasPrefix(authHeader, prefix){
		return "", errors.New("invalid Authorization header format")
	}

	key := strings.TrimSpace(strings.TrimPrefix(authHeader,prefix))
	if key == ""{
		return "", errors.New("empty bearer token")
	}
	return key, nil
}