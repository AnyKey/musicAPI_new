package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"log"
	"musicAPI/model"
	"strings"
)

func ElasticAdd(tracks []model.TrackSelect) error {
	log.SetFlags(0)

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Username: "elastic",
		Password: "changeme",
	})
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	// Build the request body.
	if tracks == nil {
		return nil
	}
	track := model.TrackSelect{tracks[0].Name, tracks[0].Artist, tracks[0].Album}
	bytes, err := json.Marshal(track)
	if err != nil {
		return errors.Wrap(err, "marshal elastic tracks")

	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:   "test",
		Body:    strings.NewReader(string(bytes)),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err = req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document", res.Status())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	return nil
}
func ElasticGet(tracks []model.TrackSelect) bool {
	var r map[string]interface{}
	log.SetFlags(0)

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Username: "elastic",
		Password: "changeme",
	})
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": tracks[0].Name,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	if r != nil {
		var checkVal interface{}
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			checkVal = hit.(map[string]interface{})["_source"]
			log.Printf(" * %v", hit.(map[string]interface{})["_source"].(map[string]interface{})["artist"])
		}
		if checkVal != nil {
			return true
		}
	}
	return false
}
