package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"musicAPI/model"
	"musicAPI/repository"
	"net/http"
	"time"
)

type LoginHandler struct {
	Repo repository.Repository
}

func (lh LoginHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	var ctx = context.Background()
	user := vars["name"]
	if user != "" {
		tokenA, err := NewTokenA(user)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			err = WriteJsonToResponse(writer, err.Error())
		}
		tokenR, err := NewTokenR(user)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			err = WriteJsonToResponse(writer, err.Error())
		}
		access, err := tokenA.SignedString([]byte("key"))
		refresh, err := tokenR.SignedString([]byte("key"))
		newToken := model.Tokens{access, refresh, true}

		if err != nil {
			return
		}
		fullToken, err := json.Marshal(newToken)
		lh.Repo.Redis.Set(ctx, "JWT:"+user, fullToken, 1*time.Hour)
		err = WriteJsonToResponse(writer, newToken)
	}

}
