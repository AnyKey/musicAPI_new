package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"musicAPI/api"
	"musicAPI/repository"
	"net/http"
	"time"
)

type AlbumHandler struct {
	Repo repository.Repository
}

func (ah AlbumHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error
	var ctx = context.Background()
	vars := mux.Vars(req)
	album, artist := vars["album"], vars["artist"]
	result := ah.Repo.GetAlbumRedis(album, artist)
	if result != nil {
		err = WriteJsonToResponse(writer, result)
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	re, err := api.AlbumInfoReq(album, artist)
	bytes, err := json.Marshal(re)
	if err == nil {
		ah.Repo.Redis.Set(ctx, "Album:"+album+"_Artist:"+artist, bytes, 5*time.Minute)
	}
	err = WriteJsonToResponse(writer, re)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		err = WriteJsonToResponse(writer, err.Error())
		if err != nil {
			log.Println(err.Error())
		}
	}

}
