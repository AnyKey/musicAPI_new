package main

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"github.com/vrischmann/envconfig"
	"log"
	"musicAPI/handlers"
	"musicAPI/repository"
	"musicAPI/user"
	tokenM "musicAPI/user/delivery"
	redisRep "musicAPI/user/repository/redis"
	tokenUseCase "musicAPI/user/usecase"
	"net/http"
	"time"
)

type config struct {
	Database    string `envconfig:"DATABASE"`
	HttpAddress string `envconfig:"HTTP_ADDRESS"`
	RedisPort   string `envconfig:"REDIS_PORT"`
	QueuePort   string `envconfig:"QUEUE_PORT"`
}
type App struct {
	tokenUC user.UseCase
}

func NewApp() *App {
	var sConfig config
	err := envconfig.Init(&sConfig)
	if err != nil {
		log.Fatalln(err)
	}
	db := mustDBConn(sConfig.Database)
	rdb := redis.NewClient(&redis.Options{
		Addr: sConfig.RedisPort,
	})
	tokenRepo := redisRep.New(rdb)
	_ = db
	return &App{
		tokenUC: tokenUseCase.New(
			tokenRepo,
		),
	}
}
func (a *App) Run() {
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

	rdb := redis.NewClient(&redis.Options{
		Addr: sConfig.RedisPort,
	})
	conn := mustDBConn(sConfig.Database)
	repo := repository.Repository{Conn: conn, Redis: rdb}

	router := mux.NewRouter()
	router.Use(tokenM.NewTokenHandler(a.tokenUC).TokenMiddleware)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Handle("/api/artist/{artist}", handlers.ArtistHandler{Repo: repo}).Methods(http.MethodGet, http.MethodOptions) //4
	router.Handle("/api/login/{name}", tokenM.NewAccTokenHandler(a.tokenUC)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/api/refresh", tokenM.NewRefTokenHandler(a.tokenUC)).Methods(http.MethodGet, http.MethodOptions)

	srv := &http.Server{
		Handler:      router,
		Addr:         sConfig.HttpAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Serve http ON", sConfig.HttpAddress)
	log.Fatal(srv.ListenAndServe())

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func main() {
	app := NewApp()
	app.Run()
}
