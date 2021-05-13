package main

import (
	"database/sql"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"github.com/vrischmann/envconfig"
	"log"
)

func newConfig() config {
	var sConfig config
	err := envconfig.Init(&sConfig)
	if err != nil {
		log.Fatalln(err)
	}
	return sConfig
}

func initPostgres(database string) *sql.DB {
	db, err := sql.Open("postgres", database)
	if err != nil {
		log.Fatalln(err)
	}
	if db.Ping() != nil {
		log.Fatalln(err)
	}
	return db
}

func initRedis(redisPort string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisPort,
	})
	return rdb
}

func initRabbitMq(queuePort string) *amqp.Channel {
	connQueue, err := amqp.Dial(queuePort)
	failOnError(err, "Failed to connect to RabbitMQ")
	//	defer connQueue.Close()
	ch, err := connQueue.Channel()
	failOnError(err, "Failed to open a channel")
	//	defer ch.Close()
	return ch
}

func initElasticSearch(name, pass string) *elasticsearch.Client {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Username: name,
		Password: pass,
	})
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	return es
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
