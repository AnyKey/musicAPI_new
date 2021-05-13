package http

import (
	"github.com/gorilla/mux"
	"log"
	"musicAPI/client"
	"musicAPI/user"
	"net/http"
)

type UserHandler struct {
	usecase user.UseCase
}

func UserHandlers(router *mux.Router, tokenUC user.UseCase) {
	UserH := UserHandler{
		usecase: tokenUC,
	}

	router.HandleFunc("/api/login/{name}", UserH.login).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/refresh", UserH.refresh).Methods(http.MethodGet, http.MethodOptions)
}
func (uh *UserHandler) refresh(w http.ResponseWriter, r *http.Request) {
	rToken := r.Header.Get("refresh")
	tokens, err := uh.usecase.RefreshToken(r.Context(), rToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = client.WriteJsonToResponse(w, err.Error())
	}
	err = client.WriteJsonToResponse(w, tokens)
	if err != nil {
		log.Println(err)
	}
	return
}
func (uh *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["name"]
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokens, err := uh.usecase.NewToken(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = client.WriteJsonToResponse(w, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = client.WriteJsonToResponse(w, tokens)
	if err != nil {
		log.Println(err)
	}
}
