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

type genreHandler struct {
	usecase music.UseCase
}

func newGenreHandler(usecase music.UseCase) *genreHandler {
	return &genreHandler{
		usecase: usecase,
	}
}

func (gh genreHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

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
	var genre = vars["genre"]
	var ctx = context.Background()
	tracks, err := gh.usecase.GenreReq(ctx, genre)
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
