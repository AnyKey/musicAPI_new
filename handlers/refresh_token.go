package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"musicAPI/model"
	"musicAPI/repository"
	"net/http"
	"time"
)

type RefreshHandler struct {
	Repo repository.Repository
}

func (rh RefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()
	rToken := r.Header.Get("refresh")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(rToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})
	if err == nil && token.Valid {
		tokenMap := make(map[string]string)
		for key, val := range claims {
			str := fmt.Sprintf("%v", val)
			tokenMap[key] = str
		}
		user := tokenMap["name"]
		res := rh.Repo.GetToken(user)

		log.Println("res - ", res.Refresh, "\n", rToken)
		if res.Refresh == rToken && res.Valid == true {
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
			newToken := model.Tokens{access, refresh, true}
			if err != nil {
				return
			}
			fullToken, err := json.Marshal(newToken)
			rh.Repo.Redis.Set(ctx, "JWT:"+user, fullToken, 1*time.Hour)
			err = WriteJsonToResponse(w, newToken)

		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	err = WriteJsonToResponse(w, err.Error())
}
