package delivery

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"musicAPI/handlers"
	"musicAPI/music"
	"net/http"
)

type trackHandler struct {
	usecase music.UseCase
}

func newTrackHandler(usecase music.UseCase) *trackHandler {
	return &trackHandler{
		usecase: usecase,
	}
}

func (th trackHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error
	var value bool
	defer func() {

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			err = handlers.WriteJsonToResponse(writer, err.Error())
		}
		if value == false {
			writer.WriteHeader(http.StatusBadRequest)
			err = handlers.WriteJsonToResponse(writer, "Bad request")
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}

	}()
	vars := mux.Vars(req)

	var trackV = vars["track"]
	var artistV = vars["artist"]
	var ctx = context.Background()
	tracks, value, err := th.usecase.TrackReq(ctx, trackV, artistV)
	if err == nil && value == true {
		err = handlers.WriteJsonToResponse(writer, tracks)
	}
	return
}
