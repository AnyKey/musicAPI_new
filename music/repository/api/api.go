package api

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"musicAPI/model"
	"musicAPI/music"
	"net/http"
	"time"
)

type Repository struct {
	text string
}

func New(text string) *Repository {
	return &Repository{
		text: text,
	}
}

func (a Repository) AlbumInfoReq(album string, artist string) (*model.Root, error) {
	http.DefaultClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, music.BaseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error create new request")
	}

	query := req.URL.Query()
	query.Add("api_key", music.ApiKey)
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
	var val model.Root
	err = json.Unmarshal(bytes, &val)
	if err != nil {
		return nil, errors.Wrap(err, "Error Unmarshal")
	}

	return &val, nil
}

func (a *Repository) TrackSearchReq(track string, artist string) (*model.OwnTrack, error) {
	http.DefaultClient = &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, music.BaseURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error create new request")
	}

	query := req.URL.Query()
	query.Add("api_key", music.ApiKey)
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
	var tracks model.TrackRoot
	err = json.Unmarshal(bytes, &tracks)
	if err != nil {
		return nil, errors.Wrap(err, "Error Unmarshal")
	}
	ts := tracks.Track
	return &ts, nil
}
