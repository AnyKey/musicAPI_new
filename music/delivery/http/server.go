package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"musicAPI/client"
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
func (mu *musicHandler) album(writer http.ResponseWriter, req *http.Request) {

	var ctx = context.Background()
	vars := mux.Vars(req)
	album, artist := vars["album"], vars["artist"]
	result, err := mu.usecase.AlbumInfoRes(ctx, album, artist)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		err = client.WriteJsonToResponse(writer, err.Error())
		if err != nil {
			log.Println(err.Error())
		}
	}
	if result != nil {
		err = client.WriteJsonToResponse(writer, result)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return
}

func (mu *musicHandler) chart(writer http.ResponseWriter, req *http.Request) {

	var err error
	var ctx = context.Background()
	defer func() {
		if err != nil {
			err = client.WriteJsonToResponse(writer, err.Error())
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}
	}()
	vars := mux.Vars(req)
	var sortTo = vars["sortto"]
	chart, err := mu.usecase.ChartReq(ctx, sortTo)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if chart != nil {
		err = client.WriteJsonToResponse(writer, chart)
		if err != nil {
			log.Println(err)
		}
		return
	}
	if chart == nil && err == nil {
		writer.WriteHeader(http.StatusNoContent)
		err = errors.New("empty response")
		return
	}
	writer.WriteHeader(http.StatusBadRequest)
}

func (mu *musicHandler) artist(writer http.ResponseWriter, req *http.Request) {

	var err error
	defer func() {
		if err != nil {
			err = client.WriteJsonToResponse(writer, err.Error())
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}
	}()
	vars := mux.Vars(req)
	var artistV = vars["artist"]
	var ctx = context.Background()
	tracks, err := mu.usecase.ArtistReq(ctx, artistV)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if tracks != nil {
		err = client.WriteJsonToResponse(writer, tracks)
		if err != nil {
			log.Println(err)
		}
		return
	}
	writer.WriteHeader(http.StatusBadRequest)

}
func (mu *musicHandler) genre(writer http.ResponseWriter, req *http.Request) {

	var err error
	defer func() {
		if err != nil {
			err = client.WriteJsonToResponse(writer, err.Error())
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}
	}()
	vars := mux.Vars(req)
	var genre = vars["genre"]
	var ctx = context.Background()
	tracks, err := mu.usecase.GenreReq(ctx, genre)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if tracks != nil {
		err = client.WriteJsonToResponse(writer, tracks)
		if err != nil {
			log.Println(err)
		}
		return
	}
	writer.WriteHeader(http.StatusBadRequest)

}

func (mu *musicHandler) track(writer http.ResponseWriter, req *http.Request) {

	var err error
	var value bool
	defer func() {

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			err = client.WriteJsonToResponse(writer, err.Error())
		}
		if value == false {
			writer.WriteHeader(http.StatusBadRequest)
			err = client.WriteJsonToResponse(writer, "Bad request")
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}

	}()
	vars := mux.Vars(req)

	var trackV = vars["track"]
	var artistV = vars["artist"]
	var ctx = context.Background()
	tracks, value, err := mu.usecase.TrackReq(ctx, trackV, artistV)
	if err == nil && value == true {
		err = client.WriteJsonToResponse(writer, tracks)
	}
	return
}
