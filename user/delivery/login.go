package delivery

import (
	"github.com/gorilla/mux"
	"log"
	"musicAPI/handlers"
	"net/http"
)

func (th UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["name"]
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokens, err := th.usecase.NewToken(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = handlers.WriteJsonToResponse(w, err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}
	err = handlers.WriteJsonToResponse(w, tokens)
	if err != nil {
		log.Println(err)
	}

}
