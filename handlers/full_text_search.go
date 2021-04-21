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
)

func ElasticHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Done!")
	return
	track := mux.Vars(r)["track"]
	var q map[string]interface{}
	var trackList model.TrackSelect

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
			"match": map[string]interface{}{
				"name": track,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Println("Error encoding query: ", err)
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
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
		for _, hit := range q["hits"].(map[string]interface{})["hits"].([]interface{}) {
			trackList = model.TrackSelect{
				hit.(map[string]interface{})["_source"].(map[string]interface{})["name"].(string),
				hit.(map[string]interface{})["_source"].(map[string]interface{})["artist"].(string),
				hit.(map[string]interface{})["_source"].(map[string]interface{})["album"].(string),
			}
			break
			err = WriteJsonToResponse(w, trackList)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				WriteJsonToResponse(w, err.Error())
			}
		}

		WriteJsonToResponse(w, trackList)
	}
}
