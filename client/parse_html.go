package client

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"html/template"
	"net/http"
)

func Template(router *mux.Router) {
	router.HandleFunc("/index", parseHtml).Methods(http.MethodGet)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		return
	}).Methods(http.MethodGet)
}
func parseHtml(w http.ResponseWriter, r *http.Request) {
	err, w := parsePage(w)
	if err != nil {
		WriteJsonToResponse(w, err.Error())
	}
	return
}

func parsePage(w http.ResponseWriter) (error, http.ResponseWriter) {
	tmpl, err := template.ParseFiles("client/template/index.html")
	if err != nil {
		WriteJsonToResponse(w, err.Error())
		return err, w
	}
	tmpl.Execute(w, nil)
	return nil, w
}

func WriteJsonToResponse(rw http.ResponseWriter, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "error while marshal json")
	}

	_, err = rw.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "error write response")
	}

	return nil
}
