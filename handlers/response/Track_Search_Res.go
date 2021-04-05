package response

import (
	"github.com/gorilla/mux"
	"musicAPI/handlers"
	"musicAPI/handlers/request"
	"musicAPI/repository"
	"net/http"
)

type TrackHandler struct {
	Repo repository.Repository
}

func (th TrackHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	// writer.Header().Set("Access-Control-Allow-Origin", "*")
	var err error

	// 4to eto?
	if err != nil {
		_ = handlers.WriteJsonToResponse(writer, map[string]string{
			"error": err.Error(),
		})
	}


	vars := mux.Vars(req)
	// non exp
	var TrackV string = vars["track"]
	var ArtistV string = vars["artist"]
	Tracks, err := th.Repo.GetTracks(TrackV, ArtistV)
	if Tracks != nil {
		_ = handlers.WriteJsonToResponse(writer, Tracks)
	}else if Tracks == nil {

		re := request.TrackSearchReq(TrackV, ArtistV)
		go th.Repo.SetTracks(re)

		_ = handlers.WriteJsonToResponse(writer, re)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		_ = handlers.WriteJsonToResponse(writer, err.Error())
	}
}