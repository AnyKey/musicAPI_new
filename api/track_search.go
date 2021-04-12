package api

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"musicAPI/model"
	"net/http"
	"time"
)

var Tracks model.TrackRoot

func TrackSearchReq(track string, artist string) (*model.OwnTrack, error) {
	http.DefaultClient = &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error create new request")
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
		return nil, errors.Wrap(err, "Error request")
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error conv bytes")
	}

	err = json.Unmarshal(bytes, &Tracks)
	if err != nil {
		return nil, errors.Wrap(err, "Error Unmarshal")
	}
	ts := Tracks.Track
	return &ts, nil
}
