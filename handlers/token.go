package handlers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"musicAPI/repository"
	"net/http"
	"time"
)

type TokenHandler struct {
	Repo repository.Repository
}

func (th TokenHandler) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if r.RequestURI == "/api/login/"+vars["name"] {
			next.ServeHTTP(w, r)
			return
		}
		if r.RequestURI == "/api/refresh/" {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get("token")
		if token != "" {
			sd := CheckToken(th, token)
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
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user
	claims["root"] = true
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	return token, nil
}
func NewTokenR(user string) (*jwt.Token, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	return token, nil
}

func CheckToken(tt TokenHandler, myToken string) bool {
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
		res := tt.Repo.GetToken(user)
		if res.Access == myToken && res.Valid == true {
			return true
		}
	}
	return false
}
