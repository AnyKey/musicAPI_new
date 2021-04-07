package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"musicAPI/repository"
	"net/http"
)

type ChartHandler struct {
	Repo repository.Repository
}

func (ch ChartHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	var err error
	defer func() {
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			err = WriteJsonToResponse(writer, err.Error())
			if err != nil {
				fmt.Println(writer, err.Error())
			}
		}
	}()
	vars := mux.Vars(req)
	var sortTo = vars["sortto"]
	chart, err := ch.Repo.GetChart(sortTo)
	if chart != nil {
		err = WriteJsonToResponse(writer, chart)
	} else if chart == nil && err == nil {
		//writer.WriteHeader(http.)
		err = WriteJsonToResponse(writer, err.Error())

	}
}
