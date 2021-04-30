package delivery

import (
	"musicAPI/handlers"
	"musicAPI/user"
	"net/http"
	"strings"
)

type TokenHandler struct { // type tokenMWR struct {}
	usecase user.UseCase
}

func NewTokenHandler(usecase user.UseCase) *TokenHandler {
	return &TokenHandler{
		usecase: usecase,
	}
}

func (th TokenHandler) TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//		w.Header().Set("Access-Control-Allow-Headers", "token")
		//		if r.Method == http.MethodOptions {
		//			return
		//		}
		if strings.HasPrefix(r.RequestURI, "/api/refresh") || strings.HasPrefix(r.RequestURI, "/api/login") {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get("token")
		if token != "" {
			sd := th.usecase.CheckToken(r.Context(), token)
			if sd == true {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		_ = handlers.WriteJsonToResponse(w, "Unauthorized")
	})

}
