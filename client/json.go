package client

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func WriteJsonToResponse(rw http.ResponseWriter, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "error while marshal json")
	}

	_, err = rw.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "error write response")
	}

	return nil
}
