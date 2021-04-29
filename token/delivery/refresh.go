package delivery

import (
	"log"
	"musicAPI/handlers"
	"musicAPI/token"
	"net/http"
)

type RefTokenHandler struct {
	usecase token.UseCase
}

func NewRefTokenHandler(usecase token.UseCase) *RefTokenHandler {
	return &RefTokenHandler{
		usecase: usecase,
	}
}
func (rh RefTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rToken := r.Header.Get("refresh")
	tokens, err := rh.usecase.RefreshToken(r.Context(), rToken)
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
