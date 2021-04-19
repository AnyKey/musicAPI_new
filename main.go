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
	"net/http"
	"time"
)

type config struct {
	Database    string `envconfig:"DATABASE"`
	HttpAddress string `envconfig:"HTTP_ADDRESS"`
	RedisPort   string `envconfig:"REDIS_PORT"`
	QueuePort   string `envconfig:"QUEUE_PORT"`
}

func main() {

	var sconfig config
	err := envconfig.Init(&sconfig)
	if err != nil {
		panic(err)
	}
	connQueue, err := amqp.Dial(sconfig.QueuePort)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connQueue.Close()
	ch, err := connQueue.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: sconfig.RedisPort,
	})
	conn := MustDBConn(sconfig.Database)
	repo := repository.Repository{Conn: conn, Redis: rdb}

	router := mux.NewRouter()
	router.Use(handlers.TokenHandler{Repo: repo, Chann: ch}.AuthUser)

	router.Handle("/api/track/{artist}/{track}", handlers.TrackHandler{Repo: repo}).Methods(http.MethodGet) //1
	router.Handle("/api/album/{artist}/{album}", handlers.AlbumHandler{Repo: repo}).Methods(http.MethodGet) //2
	router.Handle("/api/genre/{genre}", handlers.GenreHandler{Repo: repo}).Methods(http.MethodGet)          //3
	router.Handle("/api/artist/{artist}", handlers.ArtistHandler{Repo: repo}).Methods(http.MethodGet)       //4
	router.Handle("/api/chart/{sortto}", handlers.ChartHandler{Repo: repo}).Methods(http.MethodGet)         //5
	router.Handle("/api/login/{name}", handlers.LoginHandler{Repo: repo}).Methods(http.MethodGet)
	router.Handle("/api/refresh", handlers.RefreshHandler{Repo: repo}).Methods(http.MethodGet)
	//router.HandleFunc("/api/elastic", handlers.ElasticHandler()).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      router,
		Addr:         sconfig.HttpAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Serve http ON", sconfig.HttpAddress)
	log.Fatal(srv.ListenAndServe())

}

func MustDBConn(database string) *sql.DB {
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
