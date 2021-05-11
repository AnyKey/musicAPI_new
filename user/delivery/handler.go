package delivery

import (
	"github.com/gorilla/mux"
	"musicAPI/user"
	"net/http"
)

func UserHandlers(router *mux.Router, tokenUC user.UseCase) {
	router.Handle("/api/login/{name}", NewAccTokenHandler(tokenUC)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/api/refresh", NewRefTokenHandler(tokenUC)).Methods(http.MethodGet, http.MethodOptions)
}
