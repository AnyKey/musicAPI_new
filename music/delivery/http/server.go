package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
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
	router.HandleFunc("/api/like/{artist}/{track}", musicH.like).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/likelist/{artist}/{track}", musicH.likeList).Methods(http.MethodGet, http.MethodOptions)
}
func (mu *musicHandler) album(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	vars := mux.Vars(r)
	album, artist := vars["album"], vars["artist"]
	result, err := mu.usecase.AlbumInfoRes(ctx, album, artist)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = WriteJsonToResponse(w, err.Error())
		if err != nil {
			log.Println(err.Error())
		}
	}
	if result != nil {
		err = WriteJsonToResponse(w, result)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return
}

func (mu *musicHandler) chart(w http.ResponseWriter, r *http.Request) {

	var err error
	ctx := context.Background()
	defer func() {
		if err != nil {
			err = WriteJsonToResponse(w, err.Error())
			if err != nil {
				fmt.Println(w, err.Error())
			}
		}
	}()
	vars := mux.Vars(r)
	sortTo := vars["sortto"]
	chart, err := mu.usecase.ChartReq(ctx, sortTo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if chart != nil {
		err = WriteJsonToResponse(w, chart)
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
			err = WriteJsonToResponse(w, err.Error())
			if err != nil {
				fmt.Println(w, err.Error())
			}
		}
	}()
	vars := mux.Vars(r)
	artist := vars["artist"]
	ctx := context.Background()
	tracks, err := mu.usecase.ArtistReq(ctx, artist)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if tracks != nil {
		err = WriteJsonToResponse(w, tracks)
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
			err = WriteJsonToResponse(w, err.Error())
			if err != nil {
				fmt.Println(w, err.Error())
			}
		}
	}()
	vars := mux.Vars(r)
	genre := vars["genre"]
	ctx := context.Background()
	tracks, err := mu.usecase.GenreReq(ctx, genre)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if tracks != nil {
		err = WriteJsonToResponse(w, tracks)
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
			err = WriteJsonToResponse(w, err.Error())
		}
		if value == false {
			w.WriteHeader(http.StatusBadRequest)
			err = WriteJsonToResponse(w, "Bad request")
			if err != nil {
				fmt.Println(w, err.Error())
			}
		}

	}()
	vars := mux.Vars(r)

	track := vars["track"]
	artist := vars["artist"]
	ctx := context.Background()
	tracks, value, err := mu.usecase.TrackReq(ctx, track, artist)
	if err == nil && value == true {
		err = WriteJsonToResponse(w, tracks)
	}
	return
}
func (mu *musicHandler) like(w http.ResponseWriter, r *http.Request) {

	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = WriteJsonToResponse(w, err.Error())
		}
	}()
	vars := mux.Vars(r)

	track := vars["track"]
	artist := vars["artist"]
	token := r.Header.Get("token")
	message, err := mu.usecase.SetLike(track, artist, token)
	if err == nil {
		err = WriteJsonToResponse(w, message)
	}
	return
}

func (mu *musicHandler) likeList(w http.ResponseWriter, r *http.Request) {

	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = WriteJsonToResponse(w, err.Error())
		}
	}()
	vars := mux.Vars(r)

	track := vars["track"]
	artist := vars["artist"]
	token := r.Header.Get("token")
	message, err := mu.usecase.GetLike(track, artist, token)
	if err == nil {
		err = WriteJsonToResponse(w, message)
	}
	return
}

func WriteJsonToResponse(rw http.ResponseWriter, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "error while marshal json")
	}

	_, err = rw.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "error write response")
	}

	return nil
}
