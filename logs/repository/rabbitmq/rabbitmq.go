package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type Repository struct {
	Chan *amqp.Channel
}

func New(ch *amqp.Channel) *Repository {
	return &Repository{
		Chan: ch,
	}
}
func (repo *Repository) NewQueue() amqp.Queue {
	q, err := repo.Chan.QueueDeclare(
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
	return q
}
func (repo *Repository) PushToChan(body []byte, q amqp.Queue) error {
	err := repo.Chan.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Println(err)
	}
	return nil
}
