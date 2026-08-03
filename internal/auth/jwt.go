package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func SignUserJWT(username string, expireAt time.Time) (string, error) {
	signingKey := os.Getenv("JWT_SECRET")

	if signingKey == "" {
		fmt.Print("Set JWT_SECRET as Env variable")
		return "", errors.New("Did Set the JWT_SECRET variable")
	}

	claims := UserClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(signingKey))

	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("error while siging the token")
	}

	return ss, nil
}

func ValidateJWT(tokenString string) (*UserClaims, error) {
	signingKey := os.Getenv("JWT_SECRET")
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(signingKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
