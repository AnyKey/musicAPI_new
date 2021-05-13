package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"musicAPI/helper"
	"musicAPI/music"
	"net/http"
)

type musicHandler struct {
	usecase music.UseCase
}

func MusicHandlers(router *mux.Router, musicUC music.UseCase) {

	musicH := musicHandler{musicUC}
	router.HandleFunc("/api/album/{artist}/{album}", musicH.album).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/chart/{sortto}", musicH.chart).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/artist/{artist}", musicH.artist).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/genre/{genre}", musicH.genre).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/track/{artist}/{track}", musicH.track).Methods(http.MethodGet, http.MethodOptions)
}
func (mu *musicHandler) album(w http.ResponseWriter, r *http.Request) {

	var ctx = context.Background()
	vars := mux.Vars(r)
	album, artist := vars["album"], vars["artist"]
	result, err := mu.usecase.AlbumInfoRes(ctx, album, artist)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = helper.WriteJsonToResponse(w, err.Error())
		if err != nil {
			log.Println(err.Error())
		}
	}
	if result != nil {
		err = helper.WriteJsonToResponse(w, result)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return
}

func (mu *musicHandler) chart(w http.ResponseWriter, r *http.Request) {

	var err error
	var ctx = context.Background()
	defer func() {
		if err != nil {
			err = helper.WriteJsonToResponse(w, err.Error())
			if err != nil {
				fmt.Println(w, err.Error())
			}
		}
	}()
	vars := mux.Vars(r)
	var sortTo = vars["sortto"]
	chart, err := mu.usecase.ChartReq(ctx, sortTo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if chart != nil {
		err = helper.WriteJsonToResponse(w, chart)
		if err != nil {
			log.Println(err)
		}
		return
	}
	if chart == nil && err == nil {
		w.WriteHeader(http.StatusNoContent)
		err = errors.New("empty response")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (mu *musicHandler) artist(w http.ResponseWriter, r *http.Request) {

	var err error
	defer func() {
		if err != nil {
			err = helper.WriteJsonToResponse(w, err.Error())
			if err != nil {
				fmt.Println(w, err.Error())
			}
		}
	}()
	vars := mux.Vars(r)
	var artistV = vars["artist"]
	var ctx = context.Background()
	tracks, err := mu.usecase.ArtistReq(ctx, artistV)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if tracks != nil {
		err = helper.WriteJsonToResponse(w, tracks)
		if err != nil {
			log.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}
func (mu *musicHandler) genre(w http.ResponseWriter, r *http.Request) {

	var err error
	defer func() {
		if err != nil {
			err = helper.WriteJsonToResponse(w, err.Error())
			if err != nil {
				fmt.Println(w, err.Error())
			}
		}
	}()
	vars := mux.Vars(r)
	var genre = vars["genre"]
	var ctx = context.Background()
	tracks, err := mu.usecase.GenreReq(ctx, genre)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if tracks != nil {
		err = helper.WriteJsonToResponse(w, tracks)
		if err != nil {
			log.Println(err)
		}
		return
	}
	w.WriteHeader(http.StatusBadRequest)

}

func (mu *musicHandler) track(w http.ResponseWriter, r *http.Request) {

	var err error
	var value bool
	defer func() {

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = helper.WriteJsonToResponse(w, err.Error())
		}
		if value == false {
			w.WriteHeader(http.StatusBadRequest)
			err = helper.WriteJsonToResponse(w, "Bad request")
			if err != nil {
				fmt.Println(w, err.Error())
			}
		}

	}()
	vars := mux.Vars(r)

	var trackV = vars["track"]
	var artistV = vars["artist"]
	var ctx = context.Background()
	tracks, value, err := mu.usecase.TrackReq(ctx, trackV, artistV)
	if err == nil && value == true {
		err = helper.WriteJsonToResponse(w, tracks)
	}
	return
}
