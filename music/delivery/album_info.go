package delivery

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"musicAPI/handlers"
	"musicAPI/music"
	"net/http"
)

type albumInfoHandler struct {
	usecase music.UseCase
}

func newAlbumInfoHandler(usecase music.UseCase) *albumInfoHandler {
	return &albumInfoHandler{
		usecase: usecase,
	}
}

func (aih albumInfoHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var ctx = context.Background()
	vars := mux.Vars(req)
	album, artist := vars["album"], vars["artist"]
	result, err := aih.usecase.AlbumInfoRes(ctx, album, artist)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		err = handlers.WriteJsonToResponse(writer, err.Error())
		if err != nil {
			log.Println(err.Error())
		}
	}
	if result != nil {
		err = handlers.WriteJsonToResponse(writer, result)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return
}
