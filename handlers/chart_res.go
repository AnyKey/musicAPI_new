package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"musicAPI/repository"
	"net/http"
	"time"
)

type ChartHandler struct {
	Repo repository.Repository
}

func (ch ChartHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error
	var ctx = context.Background()
	defer func() {
		if err != nil {
			err = WriteJsonToResponse(writer, err.Error())
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}
	}()
	vars := mux.Vars(req)
	var sortTo = vars["sortto"]
	chart := ch.Repo.GetChartRedis(sortTo)
	if chart != nil {
		err = WriteJsonToResponse(writer, chart)
		if err != nil {
			log.Println(err)
		}
		return
	}
	chart, err = ch.Repo.GetChart(sortTo)
	if chart != nil {
		bytes, err := json.Marshal(chart)
		if err == nil {
			ch.Repo.Redis.Set(ctx, "SortTo:"+sortTo, bytes, 5*time.Minute)
		}
		err = WriteJsonToResponse(writer, chart)
		if err != nil {
			log.Println(err)
		}
		return
	} else if chart == nil && err == nil {
		writer.WriteHeader(http.StatusNoContent)
		err = WriteJsonToResponse(writer, errors.New("empty response"))
		if err != nil {
			log.Println(err)
		}
		return
	}
	writer.WriteHeader(http.StatusBadRequest)

}
