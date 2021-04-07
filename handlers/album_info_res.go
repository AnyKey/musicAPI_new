package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"musicAPI/api"
	"net/http"
)

func AlbumInfoRes(writer http.ResponseWriter, req *http.Request) {

	var err error
	vars := mux.Vars(req)
	re, err := api.AlbumInfoReq(vars["album"], vars["artist"])
	fmt.Println(re, err)
	err = WriteJsonToResponse(writer, re)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_ = WriteJsonToResponse(writer, err.Error())
	}
}
