package usecase

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"log"
	"musicAPI/logs"
	"time"
)

type logsUseCase struct {
	RabbitDelivery logs.Delivery
}

func New(rabbitDelivery logs.Delivery) logs.UseCase {
	return &logsUseCase{RabbitDelivery: rabbitDelivery}
}

func (lu *logsUseCase) QueueAppend(myToken string, addr string) error {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(myToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})
	if err != nil {
		return err
	}
	user := (claims["name"]).(string)
	queue := lu.RabbitDelivery.NewQueue()
	body := logs.LogBody{
		Name:   user,
		Action: addr,
		Time:   time.Now(),
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	err = lu.RabbitDelivery.PushToChan(bytes, queue)
	if err != nil {
		return err
	}
	return nil
}
