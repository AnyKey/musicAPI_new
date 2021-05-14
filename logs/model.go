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
	QueueAppend(string, string) error
}
type Delivery interface {
	NewQueue() amqp.Queue
	PushToChan([]byte, amqp.Queue) error
}
