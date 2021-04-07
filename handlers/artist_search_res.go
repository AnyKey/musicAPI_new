package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"musicAPI/repository"
	"net/http"
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
	tracks, err := ah.Repo.GetArtistTracks(artistV)
	if tracks != nil {
		err = WriteJsonToResponse(writer, tracks)
	} else if tracks == nil && err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		err = WriteJsonToResponse(writer, err.Error())

	}
}
