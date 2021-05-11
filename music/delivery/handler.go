package delivery

import (
	"github.com/gorilla/mux"
	"musicAPI/music"
	"net/http"
)

func MusicHandlers(router *mux.Router, musicUC music.UseCase) {
	router.Handle("/api/album/{artist}/{album}", newAlbumInfoHandler(musicUC)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/api/chart/{sortto}", newChartHandler(musicUC)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/api/artist/{artist}", newArtistHandler(musicUC)).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/api/genre/{genre}", newGenreHandler(musicUC)).Methods(http.MethodGet, http.MethodOptions)
}
