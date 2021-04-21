package handlers

import (
	"html/template"
	"net/http"
)

func ParseHtml(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("handlers/template/index.html", "handlers/template/elastic.js", "handlers/template/jquery-3.5.0.min.js")
	if err != nil {
		WriteJsonToResponse(w, err.Error())
		return
	}

	tmpl.Execute(w, "Название песни")
	return
}
