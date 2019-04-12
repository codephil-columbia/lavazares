package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// func addLessonRecordHandler(w http.ResponseWriter, r *http.Request) {
// 	data, err := utils.ReadBody(r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	record, err := records.NewLessonRecord(data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	err = tutorialRecordManager.Save(record)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	w.WriteHeader(200)
// 	return
// }

func getNextNonCompletedLesson(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		http.Error(w, "Missing uid", http.StatusBadRequest)
		return
	}

	lesson, err := tutorialRecordManager.GetNextNoncompletedLesson(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(lesson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getNextNonCompletedChapter(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		http.Error(w, "Missing uid", http.StatusBadRequest)
		return
	}

	chapter, err := tutorialRecordManager.GetNextNoncompletedChapter(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(chapter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getLessonRecordsForUserHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		http.Error(w, "Missing uid", http.StatusBadRequest)
		return
	}

	lessonRecords, err := tutorialRecordManager.GetLessonRecords(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(lessonRecords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTutorialHollisticLessonStatsHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		http.Error(w, "Missing uid", http.StatusBadRequest)
		return
	}

	stats, err := tutorialRecordManager.LessonsStats(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTutorialLessonStatsHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		http.Error(w, "Missing uid", http.StatusBadRequest)
		return
	}

	lessonID, ok := mux.Vars(r)["lessonid"]
	if !ok {
		http.Error(w, "Missing lessonid", http.StatusBadRequest)
		return
	}

	stat, err := tutorialRecordManager.LessonStats(lessonID, uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
