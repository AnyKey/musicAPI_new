package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/vrischmann/envconfig"
	"log"
	"musicAPI/handlers"
	"musicAPI/repository"
	"net/http"
	"time"
)

func main() {
	var sconfig config
	err := envconfig.Init(&sconfig)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: sconfig.RedisPort,
	})
	conn := MustDBConn(sconfig.Database)
	repo := repository.Repository{Conn: conn, Redis: rdb}

	router := mux.NewRouter()
	router.Handle("/api/track/{artist}/{track}", handlers.TrackHandler{Repo: repo}).Methods(http.MethodGet) //1
	router.Handle("/api/album/{artist}/{album}", handlers.AlbumHandler{Repo: repo}).Methods(http.MethodGet) //2
	router.Handle("/api/genre/{genre}", handlers.GenreHandler{Repo: repo}).Methods(http.MethodGet)          //3
	router.Handle("/api/artist/{artist}", handlers.ArtistHandler{Repo: repo}).Methods(http.MethodGet)       //4
	router.Handle("/api/chart/{sortto}", handlers.ChartHandler{Repo: repo}).Methods(http.MethodGet)         //5

	srv := &http.Server{
		Handler:      router,
		Addr:         sconfig.HttpAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Serve http ON", sconfig.HttpAddress)
	log.Fatal(srv.ListenAndServe())

}
