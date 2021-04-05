package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"musicAPI/handlers/response"
	"musicAPI/repository"
	"net/http"
	"time"
)

func MustDBConn() *sql.DB {
	db, err := sql.Open("postgres", "dbname=musicdb user=postgres password=123 port=5432 sslmode=disable" )
	if err != nil {
		log.Fatalln(err)
	}
	if db.Ping() != nil {
		log.Fatalln(err)
	}
	return db
}

func main() {
	conn := MustDBConn()
	repo := repository.Repository{Conn: conn}

	router := mux.NewRouter()
	// Запрос к API lost для получения информации по альбому и исполнителю
	router.HandleFunc("/api/album/{artist}/{album}", response.AlbumInfoRes).Methods(http.MethodGet) //2
	// Запрос к API lost/DB для получения спика треков по исполнителю и треку
	router.Handle("/api/track/{artist}/{track}", response.TrackHandler{Repo: repo}).Methods(http.MethodGet) //1
	// Запрос к DB для получения списка треков по исполнителю
	router.Handle("/api/artist/{artist}", response.ArtistHandler{Repo: repo}).Methods(http.MethodGet) //4


	srv := &http.Server{
		Handler:      router,
		Addr:         ":8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
