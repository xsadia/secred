package config

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(id string) (token string, err error) {

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = id
	atClaims["exp"] = time.Now().Add(time.Hour * 9999).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err = at.SignedString([]byte("Karinne"))

	return token, err
}
