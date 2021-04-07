package handlers

import (
	"github.com/gorilla/mux"
	"musicAPI/repository"
	"net/http"
)

type GenreHandler struct {
	Repo repository.Repository
}

func (gh GenreHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	//writer.Header().Set("Access-Control-Allow-Origin", "*")
	var err error
	vars := mux.Vars(req)
	Genre, err := gh.Repo.GetGenreTracks(vars["genre"])
	if Genre != nil {
		_ = WriteJsonToResponse(writer, Genre)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		_ = WriteJsonToResponse(writer, err.Error())

	}
}
