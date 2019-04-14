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

// ReadBodyToMap reads JSON bytes from an io.ReaderCloser and returns it as a map[string]string.
// Important to note that this probably only works for non-nested JSON, but haven't tested it.
func ReadBodyToMap(body io.ReadCloser) (map[string]string, error) {
	byt, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return nil, err
	}
	data := make(map[string]string)
	err = json.Unmarshal(byt, &data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func SendJSON(obj interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(obj)
}
