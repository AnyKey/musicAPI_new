package usecase

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"log"
	"musicAPI/logs"
	"time"
)

type logsUseCase struct {
	RabbitRepo logs.Repository
}

func New(rabbitRepo logs.Repository) logs.UseCase {
	return &logsUseCase{RabbitRepo: rabbitRepo}
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
	queue := lu.RabbitRepo.NewQueue()
	body := logs.LogBody{
		Name:   user,
		Action: addr,
		Time:   time.Now(),
	}
	bytes, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	err = lu.RabbitRepo.PushToChan(bytes, queue)
	if err != nil {
		return err
	}
	return nil
}
