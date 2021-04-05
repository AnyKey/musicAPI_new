package response

import (
	"github.com/gorilla/mux"
	"musicAPI/handlers"
	"musicAPI/repository"
	"net/http"
)

type ArtistHandler struct {
	Repo repository.Repository
}

func (ah ArtistHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	//writer.Header().Set("Access-Control-Allow-Origin", "*")
	var err error
// ?
	if err != nil {
		_ = handlers.WriteJsonToResponse(writer, map[string]string{
			"error": err.Error(),
		})
	}
	vars := mux.Vars(req)
	var ArtistV = vars["artist"]
	Tracks, err := ah.Repo.GetArtistTracks(ArtistV)
	if Tracks != nil {
		_ = handlers.WriteJsonToResponse(writer, Tracks)
	//	writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		_ = handlers.WriteJsonToResponse(writer, err.Error())

	}
}
