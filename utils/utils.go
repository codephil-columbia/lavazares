package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// RequestJSON provides a wrapper for json.Unmarshal
type RequestJSON []byte

func ReadBody(body io.ReadCloser) (RequestJSON, error) {
	json, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return nil, err
	}
	return json, nil
}

func SendJSON(obj interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(obj)
}
