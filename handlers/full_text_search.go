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
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]interface{}{
				"name": "*" + strings.ToLower(track) + "*",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
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
