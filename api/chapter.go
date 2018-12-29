package api

import (
	"lavazares/utils"
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

	err = utils.SendJSON(chapter, w)
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

	err = utils.SendJSON(chapters, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
