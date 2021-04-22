package handlers

import (
	"html/template"
	"musicAPI/model"
	"net/http"
)

func ParseHtml(w http.ResponseWriter, r *http.Request) {
	err, w := ParsePage(w, nil)
	if err != nil {
		WriteJsonToResponse(w, err.Error())
	}
	return
}

type newList struct {
	Title bool
	Line  []string
}

func ParsePage(w http.ResponseWriter, trackList []model.TrackSelect) (error, http.ResponseWriter) {
	tmpl, err := template.ParseFiles("handlers/template/index.html")
	if err != nil {
		WriteJsonToResponse(w, err.Error())
		return err, w
	}
	var list newList
	if trackList != nil {
		for i, _ := range trackList {
			list.Line = append(list.Line, trackList[i].Name+" | "+trackList[i].Artist+" | "+trackList[i].Album)
		}
		list.Title = true
	} else {
		list.Line = append(list.Line, "Empty")
	}
	tmpl.Execute(w, list)
	return nil, w
}
