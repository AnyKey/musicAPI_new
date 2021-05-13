package api

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"musicAPI/music"
	"net/http"
	"time"
)

type Delivery struct {
}

func New() *Delivery {
	return &Delivery{}
}

const (
	baseURL = "http://ws.audioscrobbler.com/2.0/"
	apiKey  = "d84296d9388306355db600e324a85b9b"
)

func (*Delivery) AlbumInfoReq(album string, artist string) (*music.Root, error) {
	http.DefaultClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error create new request")
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
		return nil, errors.Wrap(err, "Error request")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error conv bytes")
	}
	var val music.Root
	err = json.Unmarshal(bytes, &val)
	if err != nil {
		return nil, errors.Wrap(err, "Error Unmarshal")
	}

	return &val, nil
}

func (*Delivery) TrackSearchReq(track string, artist string) (*music.OwnTrack, error) {
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
	var tracks music.TrackRoot
	err = json.Unmarshal(bytes, &tracks)
	if err != nil {
		return nil, errors.Wrap(err, "Error Unmarshal")
	}
	ts := tracks.Track
	return &ts, nil
}
