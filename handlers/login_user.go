package handlers

import (
	"context"
	"github.com/gorilla/mux"
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
		uToken, err := NewToken(user)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			err = WriteJsonToResponse(writer, err.Error())
		}
		t, err := uToken.SignedString([]byte(user))
		if err != nil {
			return
		}
		lh.Repo.Redis.Set(ctx, "JWT:"+user, t, 1*time.Hour)
	}

}
