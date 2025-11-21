package util

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// TODO: Secret key as well as tls keys/certs stored in secret management service like keyvault

var secretKey = []byte("secret-key-example")

type ContextKey string

const UserIdKey ContextKey = "userId"
const ExpirationKey ContextKey = "exp"

func CreateJWTSessionToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			string(UserIdKey):     userId,
			string(ExpirationKey): time.Now().Add(time.Hour * 24).Unix(),
		},
	)
	sessionToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}
	return sessionToken, nil
}

func VerifyJWTSessionToken(sessionToken string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(sessionToken, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to verify session token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unable to get claims from token: %w", err)
	}
	return &claims, nil
}

func GetUserIdFromClaims(claims *jwt.MapClaims) (string, error) {
	userId := (*claims)[string(UserIdKey)]
	if userId == "" {
		return "", fmt.Errorf("unable to get userid from claims")
	}
	return fmt.Sprint(userId), nil
}
