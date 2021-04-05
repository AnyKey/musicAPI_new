package response

import (
	"github.com/gorilla/mux"
	"musicAPI/handlers"
	"musicAPI/handlers/request"
	"net/http"
)

func AlbumInfoRes(writer http.ResponseWriter, req *http.Request) {

	//writer.Header().Set("Access-Control-Allow-Origin", "*")
	var err error
 // ?
	if err != nil {
		_ = handlers.WriteJsonToResponse(writer, map[string]string{
			"error": err.Error(),
		})
	}

	vars := mux.Vars(req)
	re := request.AlbumInfoReq(vars["album"], vars["artist"])
	_ = handlers.WriteJsonToResponse(writer, re)

}
