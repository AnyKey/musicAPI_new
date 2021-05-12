package delivery

import (
	"github.com/gorilla/mux"
	"log"
	"musicAPI/handlers"
	"musicAPI/user"
	"net/http"
)

type UserHandler struct {
	usecase user.UseCase
}

type RefTokenHandler struct {
	usecase user.UseCase
}

func NewRefTokenHandler(usecase user.UseCase) *RefTokenHandler {
	return &RefTokenHandler{
		usecase: usecase,
	}
}

func UserHandlers(router *mux.Router, tokenUC user.UseCase) {
	UserH := UserHandler{
		usecase: tokenUC,
	}

	//router.Handle("/api/login/{name}", NewAccTokenHandler(tokenUC)).Methods(http.MethodGet, http.MethodOptions)
	//router.Handle("/api/refresh", NewRefTokenHandler(tokenUC)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/login/{name}", UserH.login)
	router.HandleFunc("/api/refresh", UserH.refresh)
}
func (u *UserHandler) refresh(w http.ResponseWriter, r *http.Request) {
	rToken := r.Header.Get("refresh")
	tokens, err := u.usecase.RefreshToken(r.Context(), rToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = handlers.WriteJsonToResponse(w, err.Error())
	}
	err = handlers.WriteJsonToResponse(w, tokens)
	if err != nil {
		log.Println(err)
	}
	return
}
func (u *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["name"]
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokens, err := u.usecase.NewToken(r.Context(), user)
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
