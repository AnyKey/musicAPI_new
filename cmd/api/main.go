package main

import (
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {
	conf := newConfig()
	router := mux.NewRouter()

	register(router, conf)

	srv := &http.Server{
		Handler:      router,
		Addr:         conf.HttpAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Serve http ON", conf.HttpAddress)
	log.Fatal(srv.ListenAndServe())
}
