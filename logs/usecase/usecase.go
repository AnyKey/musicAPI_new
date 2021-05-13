package usecase

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"log"
	"musicAPI/logs"
	"time"
)

type logsUseCase struct {
	RabbitDeli logs.Delivery
}

func New(rabbitDeli logs.Delivery) logs.UseCase {
	return &logsUseCase{RabbitDeli: rabbitDeli}
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
	queue := lu.RabbitDeli.NewQueue()
	body := logs.LogBody{
		Name:   user,
		Action: addr,
		Time:   time.Now(),
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	err = lu.RabbitDeli.PushToChan(bytes, queue)
	if err != nil {
		return err
	}
	return nil
}
