package main

import (
	"database/sql"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"musicAPI/client"
	elasticM "musicAPI/elastic/delivery"
	elastic2 "musicAPI/elastic/delivery/elastic"
	elasticUseCase "musicAPI/elastic/usecase"
	graphExpRep "musicAPI/graphexp/repository/postgres"
	graphExpUC "musicAPI/graphexp/usecase"
	logsM "musicAPI/logs/delivery"
	"musicAPI/logs/delivery/rabbitmq"
	logsUseCase "musicAPI/logs/usecase"
	elastic3 "musicAPI/music/delivery/elastic"
	musicH "musicAPI/music/delivery/http"
	dbMusicRep "musicAPI/music/repository/postgres"
	redisMusicRep "musicAPI/music/repository/redis"
	musicUseCase "musicAPI/music/usecase"
	userM "musicAPI/user/delivery"
	userH "musicAPI/user/delivery/http"
	redisUserRep "musicAPI/user/repository/redis"
	userUseCase "musicAPI/user/usecase"
)

type config struct {
	Database    string `envconfig:"DATABASE"`
	HttpAddress string `envconfig:"HTTP_ADDRESS"`
	RedisPort   string `envconfig:"REDIS_PORT"`
	QueuePort   string `envconfig:"QUEUE_PORT"`
	ElasticName string `envconfig:"ELASTIC_NAME"`
	ElasticPass string `envconfig:"ELASTIC_PASS"`
}

func register(router *mux.Router, conf config) {
	queueChan := initRabbitMq(conf.QueuePort)
	elasticClient := initElasticSearch(conf.ElasticName, conf.ElasticPass)
	redisClient := initRedis(conf.RedisPort)
	postgres := initPostgres(conf.Database)

	// graphQl
	graphRegister(router, postgres)

	// user
	userRegister(router, redisClient)

	// music
	musicRegister(router, postgres, redisClient, elasticClient)

	// logs
	logsRegister(router, queueChan)

	// es
	elasticRegister(router, elasticClient)

	// other mdws
	router.Use(mux.CORSMethodMiddleware(router))

	// render
	client.Template(router)
}
func userRegister(router *mux.Router, redis *redis.Client) {
	postgresRep := redisUserRep.New(redis)
	uc := userUseCase.New(postgresRep)
	userH.UserHandlers(router, uc)
	router.Use(userM.NewUserHandler(uc).UserMiddleware)
}
func musicRegister(router *mux.Router, postgres *sql.DB, redis *redis.Client, elastic *elasticsearch.Client) {
	redisRep := redisMusicRep.New(redis)
	postgresRep := dbMusicRep.New(postgres)
	api := musicH.New()
	es := elastic3.New(elastic)
	uc := musicUseCase.New(
		redisRep,
		postgresRep,
		api,
		es,
	)
	musicH.MusicHandlers(router, uc)
}
func logsRegister(router *mux.Router, queue *amqp.Channel) {
	rabbit := rabbitmq.New(queue)
	uc := logsUseCase.New(rabbit)
	router.Use(logsM.NewLogHandler(uc).LogMiddleware)
}
func elasticRegister(router *mux.Router, elastic *elasticsearch.Client) {
	es := elastic2.New(elastic)
	uc := elasticUseCase.New(es)
	router.Use(elasticM.NewTrackHandler(uc).WsHandler)
}
func graphRegister(router *mux.Router, postgres *sql.DB) {
	db := graphExpRep.New(postgres)
	uc := graphExpUC.New(db)
	graphQlHandler(router, uc)
}
