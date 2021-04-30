package delivery

import (
	"log"
	"musicAPI/handlers"
	"musicAPI/user"
	"net/http"
)

type RefTokenHandler struct {
	usecase user.UseCase
}

func NewRefTokenHandler(usecase user.UseCase) *RefTokenHandler {
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
