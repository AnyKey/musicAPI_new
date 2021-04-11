package handlers

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"musicAPI/repository"
	"net/http"
	"strings"
	"time"
)

type TokenHandler struct {
	Repo repository.Repository
}

func (th TokenHandler) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/api/refresh") || strings.HasPrefix(r.RequestURI, "/api/login") {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("token")
		if token != "" {
			sd := th.CheckToken(r.Context(), token)
			if sd == true {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
		_ = WriteJsonToResponse(w, "Unauthorized")
	})

}

func NewTokenA(user string) (*jwt.Token, error) {
	return jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		jwt.MapClaims{
			"name": user,
			"exp": time.Now().Add(time.Hour * 1).Unix(),
			"root": true,
		}), nil
}
func NewTokenR(user string) (*jwt.Token, error) {
	return jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		jwt.MapClaims{
			"name": user,
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		}), nil
}

func (th TokenHandler) CheckToken(ctx context.Context, myToken string) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(myToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})

	if err == nil && token.Valid {
		tokenMap := make(map[string]string)
		for key, val := range claims {
			str := fmt.Sprintf("%v", val)
			tokenMap[key] = str
		}

		user := tokenMap["name"]
		res := th.Repo.GetTokens(ctx, user)
		if res.Access == myToken && res.Valid == true {
			return true
		}
	}
	return false
}
