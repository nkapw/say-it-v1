package helper

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strconv"
	"strings"
)

func GetUserIDFromToken(r *http.Request) (int, error) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("missing Authorization header")
	}

	token := strings.Split(authHeader, " ")
	if len(token) != 2 || token[0] != "Bearer" {
		return 0, fmt.Errorf("invalid token format")
	}

	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}
