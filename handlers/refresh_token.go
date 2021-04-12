package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"musicAPI/model"
	"musicAPI/repository"
	"net/http"
)

type RefreshHandler struct {
	Repo repository.Repository
}

func (rh RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rToken := r.Header.Get("refresh")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(rToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})
	if err != nil {
		log.Println(err)
		return
	}
	if !token.Valid {
		return
		w.Write([]byte("invalid token"))
	}
	user := claims["name"].(string)
	res := rh.Repo.GetTokens(r.Context(), user)

	if res.Refresh != rToken {
		return
		w.Write([]byte("invalid refresh token"))
	}
	tokenA, err := NewTokenA(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = WriteJsonToResponse(w, err.Error())
	}
	tokenR, err := NewTokenR(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = WriteJsonToResponse(w, err.Error())
	}
	access, err := tokenA.SignedString([]byte("key"))
	refresh, err := tokenR.SignedString([]byte("key"))
	newToken := model.Tokens{
		Access:  access,
		Refresh: refresh,
		Valid:   true,
	}
	if err != nil {
		log.Println(err)
		return
	}
	err = rh.Repo.SetTokens(r.Context(), user, newToken)
	if err != nil {
		log.Println(err)
		return
	}

	err = WriteJsonToResponse(w, newToken)
	if err != nil {
		log.Println(err)
	}
	return
}
