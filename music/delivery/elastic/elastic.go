package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"log"
	"musicAPI/music"
	"strings"
)

type Delivery struct {
	Es *elasticsearch.Client
}

func New(es *elasticsearch.Client) *Delivery {
	return &Delivery{
		Es: es,
	}
}

func (d *Delivery) ElasticAdd(tracks []music.TrackSelect) error {
	// Build the request body.
	if tracks == nil {
		return nil
	}
	track := music.TrackSelect{tracks[0].Name, tracks[0].Artist, tracks[0].Album}
	bytes, err := json.Marshal(track)
	if err != nil {
		return errors.Wrap(err, "marshal elastic tracks")

	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:   "tracks",
		Body:    strings.NewReader(string(bytes)),
		Refresh: "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), d.Es)
	if err != nil {
		log.Println("Error getting response: ", err)
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
func (d *Delivery) ElasticGet(tracks []music.TrackSelect) bool {
	var buf bytes.Buffer
	var r map[string]interface{}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": tracks[0].Name,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Println("Error encoding query: ", err)
	}

	// Perform the search request.
	res, err := d.Es.Search(
		d.Es.Search.WithContext(context.Background()),
		d.Es.Search.WithIndex("tracks"),
		d.Es.Search.WithBody(&buf),
		d.Es.Search.WithTrackTotalHits(true),
	)
	if err != nil || res.StatusCode == 404 {
		log.Println("Error getting response: ", err)
		return false
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Println("Error parsing the response body: ", err)
	}
	if r != nil {
		var checkVal interface{}
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			checkVal = hit.(map[string]interface{})["_source"]
			break
		}
		if checkVal != nil {
			return true
		}
	}
	return false
}
