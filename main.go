package main

import (
	"database/sql"
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

type config struct {
	Database    string `envconfig:"DATABASE"`
	HttpAddress string `envconfig:"HTTP_ADDRESS"`
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

func main() {
	var sconfig config
	err := envconfig.Init(&sconfig)
	if err != nil {
		panic(err)
	}
	conn := MustDBConn(sconfig.Database)
	repo := repository.Repository{Conn: conn}

	router := mux.NewRouter()
	router.Handle("/api/track/{artist}/{track}", handlers.TrackHandler{Repo: repo}).Methods(http.MethodGet) //1
	router.HandleFunc("/api/album/{artist}/{album}", handlers.AlbumInfoRes).Methods(http.MethodGet)         //2
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
