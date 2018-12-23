package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type RequestJSON *map[string]interface{}

var (
	lessonLogger = log.New(os.Stdout, "RouteLessonLogger", log.Lshortfile)
)

// GetLesson
func LessonHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
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
