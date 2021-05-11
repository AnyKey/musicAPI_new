package delivery

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"musicAPI/handlers"
	"musicAPI/music"
	"net/http"
)

type chartHandler struct {
	usecase music.UseCase
}

func newChartHandler(usecase music.UseCase) *chartHandler {
	return &chartHandler{
		usecase: usecase,
	}
}

func (ch chartHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error
	var ctx = context.Background()
	defer func() {
		if err != nil {
			err = handlers.WriteJsonToResponse(writer, err.Error())
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}
	}()
	vars := mux.Vars(req)
	var sortTo = vars["sortto"]
	chart, err := ch.usecase.ChartReq(ctx, sortTo)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if chart != nil {
		err = handlers.WriteJsonToResponse(writer, chart)
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
