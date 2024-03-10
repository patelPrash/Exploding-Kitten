package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTToken(email string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	secretKey := "your_actual_secret_key_here"
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
