package delivery

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"musicAPI/handlers"
	"musicAPI/music"
	"net/http"
)

type artistHandler struct {
	usecase music.UseCase
}

func newArtistHandler(usecase music.UseCase) *artistHandler {
	return &artistHandler{
		usecase: usecase,
	}
}

func (ah artistHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error
	defer func() {
		if err != nil {
			err = handlers.WriteJsonToResponse(writer, err.Error())
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}
	}()
	vars := mux.Vars(req)
	var artistV = vars["artist"]
	var ctx = context.Background()
	tracks, err := ah.usecase.ArtistReq(ctx, artistV)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if tracks != nil {
		err = handlers.WriteJsonToResponse(writer, tracks)
		if err != nil {
			log.Println(err)
		}
		return
	}
	writer.WriteHeader(http.StatusBadRequest)

}
