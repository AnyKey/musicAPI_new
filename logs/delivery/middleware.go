package delivery

import (
	"log"
	"musicAPI/logs"
	"net/http"
)

type logMiddleHandler struct {
	usecase logs.UseCase
}

func NewLogHandler(usecase logs.UseCase) *logMiddleHandler {
	return &logMiddleHandler{
		usecase: usecase,
	}
}

func (lh *logMiddleHandler) LogMiddleware(next http.Handler) http.Handler {
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
