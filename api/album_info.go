package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"musicAPI/model"
	"net/http"
	"time"
)

const (
	baseURL = "http://ws.audioscrobbler.com/2.0/"
	apiKey  = "d84296d9388306355db600e324a85b9b"
)

var Val model.Root

func AlbumInfoReq(album string, artist string) (*model.Root, error) {
	http.DefaultClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseURL, nil)
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
	fmt.Println(resp.StatusCode)

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error conv bytes")
	}

	err = json.Unmarshal(bytes, &Val)
	if err != nil {
		return nil, errors.Wrap(err, "Error Unmarshal")
	}

	return &Val, nil
}
