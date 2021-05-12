package delivery

import (
	"musicAPI/handlers"
	"musicAPI/user"
	"net/http"
	"strings"
)

type UserMiddleHandler struct {
	usecase user.UseCase
}

func NewUserHandler(usecase user.UseCase) *UserMiddleHandler {
	return &UserMiddleHandler{
		usecase: usecase,
	}
}

func (th UserMiddleHandler) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "token")
		if r.Method == http.MethodOptions {
			return
		}
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
