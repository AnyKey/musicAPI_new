package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type Delivery struct {
	Chan *amqp.Channel
}

func New(ch *amqp.Channel) *Delivery {
	return &Delivery{
		Chan: ch,
	}
}
func (repo *Delivery) NewQueue() amqp.Queue {
	queue, err := repo.Chan.QueueDeclare(
		"main_queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Println(err)
	}
	return queue
}
func (repo *Delivery) PushToChan(body []byte, queue amqp.Queue) error {
	err := repo.Chan.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Println(err)
	}
	return nil
}
