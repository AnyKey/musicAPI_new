package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"musicAPI/api"
	"musicAPI/model"
	"musicAPI/repository"
	"net/http"
	"time"
)

type TrackHandler struct {
	Repo repository.Repository
}

func (th TrackHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error

	defer func() {
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			err = WriteJsonToResponse(writer, err.Error())
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}
	}()
	vars := mux.Vars(req)

	var trackV = vars["track"]
	var artistV = vars["artist"]
	var ctx = context.Background()
	tracks := th.Repo.GetTracksRedis(trackV, artistV)
	if tracks != nil {
		err = WriteJsonToResponse(writer, tracks)
	}
	if tracks == nil {
		tracks, err = th.Repo.GetTracks(trackV, artistV)
		if tracks != nil {
			bytes, err := json.Marshal(tracks)
			if err == nil {
				th.Repo.Redis.Set(ctx, "Track:"+trackV+"_Artist:"+artistV, bytes, 5*time.Minute)
			}
			err = WriteJsonToResponse(writer, tracks)

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
			result := structConv(re)
			err = WriteJsonToResponse(writer, result)
		}
	}

}

func structConv(trackList *model.OwnTrack) model.TrackSelect {
	return model.TrackSelect{trackList.Name, trackList.Album.Artist, trackList.Album.Album}
}
