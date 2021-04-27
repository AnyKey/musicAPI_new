package handlers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
}

type socketSend struct {
	Track       string `json:"track"`
	NameCheck   bool   `json:"nameCheck"`
	ArtistCheck bool   `json:"artistCheck"`
	AlbumCheck  bool   `json:"albumCheck"`
}

func (th TokenHandler) WsHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.HasPrefix(r.RequestURI, "/ws") {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
			th.wsSending(conn, r)
			//

			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
		return
	})

}

func (th TokenHandler) wsSending(conn *websocket.Conn, req *http.Request) {
	var auth bool
	var validToken bool
	var newWs socketSend
	var token string
	for {
		messageType, r, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if auth == false {
			token = bytesToString(r)
			validToken = th.CheckTokenForSocket(token, req)
		}
		if validToken != true && auth != true {
			return
		}
		if validToken == true && auth != true {
			auth = true
			err = conn.WriteMessage(messageType, []byte("Token Valid"))
			if err != nil {
				return
			}
		}
		if auth == true && validToken == false {
			err = json.Unmarshal(r, &newWs)
			if err != nil {
				log.Println(err)
				return
			}
			res, err := fullTextSearch(newWs)
			if err != nil {
				log.Println(err)
				return
			}
			bytes, _ := json.Marshal(res)
			err = conn.WriteMessage(messageType, bytes)
			if err != nil {
				return
			}
		}
		if auth == true {
			validToken = false
		}
	}
}
func (th TokenHandler) CheckTokenForSocket(myToken string, r *http.Request) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(myToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})

	if err == nil && token.Valid {
		user := (claims["name"]).(string)
		res := th.Repo.GetTokens(r.Context(), user)
		if res.Access == myToken && res.Valid == true {
			if err != nil {
				log.Println(err)
			}
			return true
		}
	}
	return false
}

func bytesToString(data []byte) string {
	return string(data[:])
}
