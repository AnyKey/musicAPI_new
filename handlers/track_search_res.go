package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"musicAPI/api"
	"musicAPI/repository"
	"net/http"
)

type TrackHandler struct {
	Repo repository.Repository
}

func (th TrackHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

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

	var trackV = vars["track"]
	var artistV = vars["artist"]
	tracks, err := th.Repo.GetTracks(trackV, artistV)
	if tracks != nil {
		err = WriteJsonToResponse(writer, tracks)
		return
	} else if tracks == nil {

		re, err := api.TrackSearchReq(trackV, artistV)
		if err != nil {
			fmt.Println(writer, err.Error())
		}
		go func() {
			err = th.Repo.SetTracks(*re)
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}()
		err = WriteJsonToResponse(writer, re)
	}
}
