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

	var err error
	vars := mux.Vars(req)
	Genre, err := gh.Repo.GetGenreTracks(vars["genre"])
	if Genre != nil && err == nil {
		_ = WriteJsonToResponse(writer, Genre)
	} else if Genre == nil && err == nil {
		writer.WriteHeader(http.StatusBadRequest)
		_ = WriteJsonToResponse(writer, err.Error())
	} else {
		writer.WriteHeader(http.StatusInternalServerError)
		_ = WriteJsonToResponse(writer, err.Error())
	}
}
