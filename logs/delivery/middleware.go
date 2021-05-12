package delivery

import (
	"log"
	"musicAPI/logs"
	"net/http"
)

type LogMiddleHandler struct {
	usecase logs.UseCase
}

func NewLogHandler(usecase logs.UseCase) *LogMiddleHandler {
	return &LogMiddleHandler{
		usecase: usecase,
	}
}

func (lh LogMiddleHandler) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		addr := r.RequestURI
		err := lh.usecase.QueueAppend(token, addr)
		if err != nil {
			log.Println(err)
		}
		next.ServeHTTP(w, r)
		return
	})

}
