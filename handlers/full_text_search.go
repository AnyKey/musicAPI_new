package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	"log"
	"musicAPI/model"
	"net/http"
	"strings"
)

func ElasticHandler(w http.ResponseWriter, r *http.Request) {
	track := mux.Vars(r)["track"]

	var nameValue, artistValue, albumValue string

	if boolConv(mux.Vars(r)["name"]) == true {
		nameValue = "*" + strings.ToLower(track) + "*"
	} else {
		nameValue = ""
	}
	if boolConv(mux.Vars(r)["artist"]) == true {
		artistValue = "*" + strings.ToLower(track) + "*"
	} else {
		artistValue = ""
	}
	if boolConv(mux.Vars(r)["album"]) == true {
		albumValue = "*" + strings.ToLower(track) + "*"
	} else {
		albumValue = ""
	}
	var q map[string]interface{}
	var trackList []model.TrackSelect

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Username: "elastic",
		Password: "changeme",
	})
	res, err := es.Info()
	if err != nil {
		log.Println("Error getting response: ", err)
	}
	defer res.Body.Close()
	var buf bytes.Buffer

	var appString1 = queryString{
		Fields: []string{
			"name",
		},
		Query: nameValue,
	}
	var appString2 = queryString{
		Fields: []string{
			"artist",
		},
		Query: artistValue,
	}
	var appString3 = queryString{
		Fields: []string{
			"album",
		},
		Query: albumValue,
	}
	var newList = [3]should{
		{appString1},
		{appString2},
		{appString3},
	}

	var newQuery queryReq

	for i, _ := range newList {
		newQuery.Query.Bool.Filter.Bool.Should = append(newQuery.Query.Bool.Filter.Bool.Should, newList[i])
	}
	if err := json.NewEncoder(&buf).Encode(newQuery); err != nil {
		log.Println("Error encoding query: ", err)
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("tracks"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Println("Error getting response: ", err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&q); err != nil {
		log.Println("Error parsing the response body: ", err)
	}
	if q != nil {
		i := 0

		for _, hit := range q["hits"].(map[string]interface{})["hits"].([]interface{}) {

			trackList = append(trackList, model.TrackSelect{
				Name:   hit.(map[string]interface{})["_source"].(map[string]interface{})["name"].(string),
				Artist: hit.(map[string]interface{})["_source"].(map[string]interface{})["artist"].(string),
				Album:  hit.(map[string]interface{})["_source"].(map[string]interface{})["album"].(string),
			})
			i++
		}

	}
	//err, w = ParsePage(w, trackList)
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	WriteJsonToResponse(w, trackList)
	return
}

func boolConv(id string) bool {
	switch id {
	case "true":
		return true
	case "false":
		return false
	default:
		return false
	}
}

type queryString struct {
	Fields []string `json:"fields"`
	Query  string   `json:"query"`
}
type should struct {
	QueryString queryString `json:"query_string"`
}
type queryReq struct {
	Query struct {
		Bool struct {
			Filter struct {
				Bool struct {
					Should []should `json:"should"`
				} `json:"bool"`
			} `json:"filter"`
		} `json:"bool"`
	} `json:"query"`
}
