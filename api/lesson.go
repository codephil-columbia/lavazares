package api

import (
	"lavazares/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// LessonHandler corresponds to GET /lesson/{id}
func LessonHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, "id not found in request", http.StatusBadRequest)
		return
	}

	lesson, err := lessonManager.GetLesson(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.SendJSON(lesson, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// LessonsHandler corresponds to GET /lesson
func LessonsHandler(w http.ResponseWriter, r *http.Request) {
	lessons, err := lessonManager.GetLessons()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = utils.SendJSON(lessons, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
