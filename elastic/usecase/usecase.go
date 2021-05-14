package usecase

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"log"
	"musicAPI/elastic"
)

type elasticUseCase struct {
	ElasticRepo elastic.Delivery
}

func New(elasticRepo elastic.Delivery) elastic.UseCase {
	return &elasticUseCase{ElasticRepo: elasticRepo}
}
func (euc *elasticUseCase) WsSending(conn *websocket.Conn) {
	var auth bool
	var validToken bool
	var newWs elastic.SocketSend
	var token string
	for {
		messageType, r, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if auth == false {
			token = bytesToString(r)
			validToken = checkTokenForSocket(token)
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
			res, err := euc.ElasticRepo.FullTextSearch(newWs)
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
func checkTokenForSocket(myToken string) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(myToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})

	if err == nil && token.Valid {
		return true
	}
	return false
}

func bytesToString(data []byte) string {
	return string(data[:])
}
