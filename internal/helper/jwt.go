package helper

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignUserJWT(username string, expireAt jwt.NumericDate) (string, error) {
	signingKey := os.Getenv("JWT_SECRET")

	if signingKey == "" {
		log.Fatal("Set JWT_SECRET as Env variable")
		return "", errors.New("Did Set the JWT_SECRET variable")
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, err := token.SignedString(signingKey)

	if err != nil {
		return "", errors.New("error while siging the token")
	}

	return ss, nil
}
