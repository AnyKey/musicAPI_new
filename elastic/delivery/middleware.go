package delivery

import (
	"log"
	"musicAPI/elastic"
	"net/http"
	"strings"
)

type elasticHandler struct {
	usecase elastic.UseCase
}

func NewTrackHandler(usecase elastic.UseCase) *elasticHandler {
	return &elasticHandler{
		usecase: usecase,
	}
}
func (eh elasticHandler) WsHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.HasPrefix(r.RequestURI, "/ws") {
			conn, err := elastic.Upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
			eh.usecase.WsSending(conn)
			//

			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
		return
	})

}
