package logs

import (
	"github.com/streadway/amqp"
	"time"
)

type LogBody struct {
	Name   string    `json:"name"`
	Action string    `json:"action"`
	Time   time.Time `json:"time"`
}
type UseCase interface {
	QueueAppend(myToken string, addr string) error
}
type Repository interface {
	NewQueue() amqp.Queue
	PushToChan(body []byte, q amqp.Queue) error
}
