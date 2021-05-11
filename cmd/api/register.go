package main

import (
	"database/sql"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"github.com/vrischmann/envconfig"
	"log"
)

type config struct {
	Database    string `envconfig:"DATABASE"`
	HttpAddress string `envconfig:"HTTP_ADDRESS"`
	RedisPort   string `envconfig:"REDIS_PORT"`
	QueuePort   string `envconfig:"QUEUE_PORT"`
}
type register struct {
	dbConn  *sql.DB
	rConn   *redis.Client
	qConn   *amqp.Connection
	esConn  *elasticsearch.Client
	address string
}

func NewReg() *register {
	var sConfig config
	err := envconfig.Init(&sConfig)
	if err != nil {
		panic(err)
	}
	connQueue, err := amqp.Dial(sConfig.QueuePort)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connQueue.Close()
	ch, err := connQueue.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Username: "elastic",
		Password: "changeme",
	})
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: sConfig.RedisPort,
	})
	conn := mustDBConn(sConfig.Database)
	return &register{
		dbConn:  conn,
		rConn:   rdb,
		qConn:   connQueue,
		esConn:  es,
		address: sConfig.HttpAddress,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
func mustDBConn(database string) *sql.DB {
	db, err := sql.Open("postgres", database)
	if err != nil {
		log.Fatalln(err)
	}
	if db.Ping() != nil {
		log.Fatalln(err)
	}
	return db
}
