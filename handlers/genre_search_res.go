package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"musicAPI/repository"
	"net/http"
	"time"
)

type GenreHandler struct {
	Repo repository.Repository
}

func (gh GenreHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error
	var ctx = context.Background()
	vars := mux.Vars(req)
	Genre := gh.Repo.GetGenreRedis(vars["genre"])
	if Genre != nil {
		err = WriteJsonToResponse(writer, Genre)
		return
	}

	Genre, err = gh.Repo.GetGenreTracks(vars["genre"])
	if Genre != nil && err == nil {
		bytes, err := json.Marshal(Genre)
		if err == nil {
			gh.Repo.Redis.Set(ctx, "Genre:"+vars["genre"], bytes, 5*time.Minute)
		}
		err = WriteJsonToResponse(writer, Genre)
		return
	} else if Genre == nil && err == nil {
		writer.WriteHeader(http.StatusBadRequest)
		_ = WriteJsonToResponse(writer, err.Error())
		return
	}
	writer.WriteHeader(http.StatusInternalServerError)
	_ = WriteJsonToResponse(writer, err.Error())
}
