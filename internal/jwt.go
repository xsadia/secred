package internal

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var notSoSecret = "Karinne"

func CreateToken(id string, exp time.Duration) (token string, err error) {

	claims := jwt.MapClaims{}
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Minute * exp).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = at.SignedString([]byte(notSoSecret))

	return token, err
}

func VerifyToken(token string) (jwt.MapClaims, error) {

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(notSoSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

// func ExtractUser(token string) (repository.User, error) {
// 	claims, err := VerifyToken(token)

// 	if err != nil {
// 		return repository.User{}, err
// 	}

// 	uid := fmt.Sprintf("%v", claims["user_id"])

// 	u := repository.User{Id: uid}

// 	if err = u.GetUserById(db); err != nil {
// 		return repository.User{}, errors.New("invalid token claim")
// 	}

// 	return u, nil
// }
