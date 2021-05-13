package http

import (
	"github.com/gorilla/mux"
	"log"
	"musicAPI/helper"
	"musicAPI/user"
	"net/http"
)

type userHandler struct {
	usecase user.UseCase
}

func UserHandlers(router *mux.Router, tokenUC user.UseCase) {
	UserH := userHandler{
		usecase: tokenUC,
	}

	router.HandleFunc("/api/login/{name}", UserH.login).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/refresh", UserH.refresh).Methods(http.MethodGet, http.MethodOptions)
}
func (uh *userHandler) refresh(w http.ResponseWriter, r *http.Request) {
	rToken := r.Header.Get("refresh")
	tokens, err := uh.usecase.RefreshToken(r.Context(), rToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = helper.WriteJsonToResponse(w, err.Error())
	}
	err = helper.WriteJsonToResponse(w, tokens)
	if err != nil {
		log.Println(err)
	}
	return
}
func (uh *userHandler) login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["name"]
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokens, err := uh.usecase.NewToken(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = helper.WriteJsonToResponse(w, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = helper.WriteJsonToResponse(w, tokens)
	if err != nil {
		log.Println(err)
	}
}
