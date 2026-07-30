package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignUserJWT(username string, expireAt time.Time) (string, error) {
	signingKey := os.Getenv("JWT_SECRET")

	if signingKey == "" {
		log.Fatal("Set JWT_SECRET as Env variable")
		return "", errors.New("Did Set the JWT_SECRET variable")
	}

	claims := jwt.MapClaims{
		"exp":      expireAt.Unix(),
		"username": username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(signingKey))

	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("error while siging the token")
	}

	return ss, nil
}
