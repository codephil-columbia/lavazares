package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// RequestJSON provides a wrapper for json.Unmarshal
type RequestJSON *map[string]interface{}

func sendJSON(obj interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(obj)
}

func requestToJSON(reader io.ReadCloser) (RequestJSON, error) {
	m := make(map[string]interface{})

	body, err := ioutil.ReadAll(reader)
	defer reader.Close()

	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	return &m, err
}
