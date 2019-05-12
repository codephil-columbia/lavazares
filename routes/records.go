package routes

import (
	"encoding/json"
	"lavazares/records"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func saveTutorialRecord(w http.ResponseWriter, r *http.Request) {
	typ, ok := mux.Vars(r)["type"]
	if !ok {
		http.Error(w, errMissingPathVar{"type"}.Error(), http.StatusBadRequest)
		return
	}

	// We either save a chapter or a lesson record given the type field
	var err error
	if typ == "chapter" {
		chapterRecord := records.ChapterRecord{}
		json.NewDecoder(r.Body).Decode(&chapterRecord)
		err = records.Save(tutorialRecordManager, chapterRecord)
	} else {
		lessonRecord := records.LessonRecord{}
		json.NewDecoder(r.Body).Decode(&lessonRecord)
		err = records.Save(tutorialRecordManager, lessonRecord)
	}

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

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
