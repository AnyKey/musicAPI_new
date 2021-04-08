package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func NewToken(user string) (*jwt.Token, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	return token, nil
}
