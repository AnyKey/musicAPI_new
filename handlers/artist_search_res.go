package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"musicAPI/repository"
	"net/http"
	"time"
)

type ArtistHandler struct {
	Repo repository.Repository
}

func (ah ArtistHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error
	defer func() {
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = WriteJsonToResponse(writer, err.Error())
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}
	}()
	vars := mux.Vars(req)
	var artistV = vars["artist"]
	var ctx = context.Background()
	tracks := ah.Repo.GetArtistRedis(artistV)
	if tracks != nil {
		err = WriteJsonToResponse(writer, tracks)
	}
	if tracks == nil {
		tracks, err = ah.Repo.GetArtistTracks(artistV)
		if tracks != nil {
			bytes, err := json.Marshal(tracks)
			if err == nil {
				ah.Repo.Redis.Set(ctx, "Artist:"+artistV, bytes, 5*time.Minute)
			}
			err = WriteJsonToResponse(writer, tracks)
		} else if tracks == nil && err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			err = WriteJsonToResponse(writer, err.Error())
		}

	}

}
