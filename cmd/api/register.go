package main

import (
	"database/sql"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"musicAPI/client"
	elasticMdw "musicAPI/elastic/delivery"
	elasticEs "musicAPI/elastic/delivery/elastic"
	elasticUseCase "musicAPI/elastic/usecase"
	graphExpRep "musicAPI/graphexp/repository/postgres"
	graphExpUC "musicAPI/graphexp/usecase"
	logsMdw "musicAPI/logs/delivery"
	"musicAPI/logs/delivery/rabbitmq"
	logsUseCase "musicAPI/logs/usecase"
	elastic3 "musicAPI/music/delivery/elastic"
	musicHttp "musicAPI/music/delivery/http"
	dbMusicRep "musicAPI/music/repository/postgres"
	redisMusicRep "musicAPI/music/repository/redis"
	musicUseCase "musicAPI/music/usecase"
	userMdw "musicAPI/user/delivery"
	userHttp "musicAPI/user/delivery/http"
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
	userHttp.UserHandlers(router, uc)
	router.Use(userMdw.NewUserHandler(uc).UserMiddleware)
}
func musicRegister(router *mux.Router, postgres *sql.DB, redis *redis.Client, elastic *elasticsearch.Client) {
	redisRep := redisMusicRep.New(redis)
	postgresRep := dbMusicRep.New(postgres)
	api := musicHttp.New()
	es := elastic3.New(elastic)
	uc := musicUseCase.New(
		redisRep,
		postgresRep,
		api,
		es,
	)
	musicHttp.MusicHandlers(router, uc)
}
func logsRegister(router *mux.Router, queue *amqp.Channel) {
	rabbit := rabbitmq.New(queue)
	uc := logsUseCase.New(rabbit)
	router.Use(logsMdw.NewLogHandler(uc).LogMiddleware)
}
func elasticRegister(router *mux.Router, elastic *elasticsearch.Client) {
	es := elasticEs.New(elastic)
	uc := elasticUseCase.New(es)
	router.Use(elasticMdw.NewTrackHandler(uc).WsHandler)
}
func graphRegister(router *mux.Router, postgres *sql.DB) {
	db := graphExpRep.New(postgres)
	uc := graphExpUC.New(db)
	graphQlHandler(router, uc)
}
