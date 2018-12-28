package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ChapterHandler corresponds to GET /chapter/{id}
func ChapterHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "id not found in request", http.StatusBadRequest)
		return
	}

	chapter, err := chapterManager.GetChapter(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sendJSON(chapter, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// ChaptersHandler corresponds to GET /chapter
func ChaptersHandler(w http.ResponseWriter, r *http.Request) {
	chapters, err := chapterManager.GetChapters()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sendJSON(chapters, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
