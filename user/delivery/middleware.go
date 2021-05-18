package delivery

import (
	"encoding/json"
	"github.com/pkg/errors"
	"musicAPI/user"
	"net/http"
	"strings"
)

type userMiddleHandler struct {
	usecase user.UseCase
}

func NewUserHandler(usecase user.UseCase) *userMiddleHandler {
	return &userMiddleHandler{
		usecase: usecase,
	}
}

func (th *userMiddleHandler) UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "token")
		if r.Method == http.MethodOptions {
			return
		}
		if strings.HasPrefix(r.RequestURI, "/api/refresh") || strings.HasPrefix(r.RequestURI, "/api/login") || strings.HasPrefix(r.RequestURI, "/ws") || strings.HasPrefix(r.RequestURI, "/index") || strings.HasPrefix(r.RequestURI, "/") {
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
		_ = writeJsonToResponse(w, "Unauthorized")
	})

}

func writeJsonToResponse(rw http.ResponseWriter, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "error while marshal json")
	}

	_, err = rw.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "error write response")
	}

	return nil
}
