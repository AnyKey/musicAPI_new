package delivery

import (
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
}

func (eh *elasticHandler) WsHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.HasPrefix(r.RequestURI, "/ws") {
			conn, err := upgrader.Upgrade(w, r, nil)
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
