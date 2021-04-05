package request

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type TrackSearch struct {
	Name      string `json:"name"`
	Artist    string `json:"artist"`
	Listeners string `json:"listeners"`
}
type ResTrackSearch struct {
	Result struct {
		Trackmatches struct {
			Track []TrackSearch `json:"track"`
		} `json:"trackmatches"`
	} `json:"results"`
}
type OwnTrack struct{
	Name string `json:"name"`
	Album TrackAlbum `json:"album"`
	Listeners string `json:"listeners"`
	Playcount string `json:"playcount"`
	TopTags struct {
		Genre []Tags `json:"tag"`
	} `json:"toptags"`
}
type Tags struct{
	Tag string `json:"name"`
}
type TrackAlbum struct {
	Artist string `json:"artist"`
	Album string `json:"title"`
}
type TrackRoot struct{
	Track OwnTrack `json:"track"`
}
var Tracks TrackRoot

func TrackSearchReq(track string, artist string) OwnTrack {
	http.DefaultClient = &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseURL, nil)
	if err != nil {
		log.Fatal(err) // za4em fatal?
	}

	query := req.URL.Query()
	query.Add("api_key", apiKey)
	query.Add("artist", artist)
	query.Add("track", track)
	query.Add("format", "json")
	query.Add("method", "track.getInfo")
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

	err = json.Unmarshal(bytes, &Tracks)
	if err != nil {
		log.Fatal(err)
	}
	var ts OwnTrack = Tracks.Track
	return ts
}
