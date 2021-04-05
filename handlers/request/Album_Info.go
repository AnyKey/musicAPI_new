package request

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	baseURL = "http://ws.audioscrobbler.com/2.0/"
	apiKey  = "d84296d9388306355db600e324a85b9b"
)

type Track struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type Album struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Tracks struct {
		Tracks []Track `json:"track"`
	} `json:"tracks"`
}

type Root struct {
	Album Album `json:"album"`
}

var Val Root

func AlbumInfoReq(album string, artist string) Root {
	http.DefaultClient = &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	query := req.URL.Query()
	query.Add("api_key", apiKey)
	query.Add("artist", artist)
	query.Add("album", album)
	query.Add("format", "json")
	query.Add("method", "album.getinfo")
	req.URL.RawQuery = query.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &Val)
	if err != nil {
		log.Fatal(err)
	}

	return Val
}
