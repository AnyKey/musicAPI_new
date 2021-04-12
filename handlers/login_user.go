package handlers

import (
	"github.com/gorilla/mux"
	"log"
	"musicAPI/model"
	"musicAPI/repository"
	"net/http"
)

type LoginHandler struct {
	Repo repository.Repository
}

func (lh LoginHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	user := vars["name"]
	if user == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenA, err := NewTokenA(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		err = WriteJsonToResponse(writer, err.Error())
		if err != nil {
			log.Println(err)
		}
	}
	tokenR, err := NewTokenR(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		err = WriteJsonToResponse(writer, err.Error())
		if err != nil {
			log.Println(err)
		}
	}
	access, err := tokenA.SignedString([]byte("key"))
	refresh, err := tokenR.SignedString([]byte("key"))
	if err != nil {
		log.Println(err)
	}

	newTokens := model.Tokens{
		Access:  access,
		Refresh: refresh,
		Valid:   true,
	}
	err = lh.Repo.SetTokens(req.Context(), user, newTokens)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		err = WriteJsonToResponse(writer, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = WriteJsonToResponse(writer, newTokens)
	if err != nil {
		log.Println(err)
	}

}
