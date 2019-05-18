package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"lavazares/records"
	"log"
	"net/http"
)

type (
	errServerError struct {
		err string
	}
)

func (err errServerError) Error() string {
	return fmt.Sprintln(err.err)
}

func getCurrentLesson(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		http.Error(w, errMissingPathVar{MissingVar: "uid"}.Error(), http.StatusBadRequest)
		return
	}

	currentLesson, err := records.CurrentLesson(lessonRecordManager, lessonManager, uid)
	if err != nil {
		log.Println(err)
		http.Error(w, errServerError{err: "Error getting current lesson"}.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(currentLesson)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error sending json", http.StatusInternalServerError)
		return
	}
}

func getCurrentChapter(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		errVal := errMissingPathVar{MissingVar: "uid"}
		log.Println(errVal)
		http.Error(w, errVal.Error(), http.StatusBadRequest)
		return
	}

	currentChapter, err := records.CurrentChapter(chapterRecordManager, chapterManager, uid)
	if err != nil {
		log.Println(err)
		http.Error(w, errServerError{"Error getting current chapter"}.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(currentChapter)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error sending json", http.StatusInternalServerError)
		return
	}
}

func saveTutorialLessonRecordHandler(w http.ResponseWriter, r *http.Request) {
	lessonRecord := records.LessonRecord{}
	err := json.NewDecoder(r.Body).Decode(&lessonRecord)
	if err != nil {
		log.Println(err)
		http.Error(w, errInvalidBody.Error(), http.StatusBadRequest)
		return
	}

	err = records.SaveLessonRecord(lessonRecordManager, chapterRecordManager, lessonRecord)
	if err != nil {
		log.Println(err)
		http.Error(w, errServerError{err: "Problem saving record"}.Error(), http.StatusInternalServerError)
		return
	}
}

func saveTutorialChapterRecordHandler(w http.ResponseWriter, r *http.Request) {
	chapterRecord := records.ChapterRecord{}
	err := json.NewDecoder(r.Body).Decode(&chapterRecord)
	if err != nil {
		log.Println(err)
		http.Error(w, errInvalidBody.Error(), http.StatusBadRequest)
		return
	}

	err = records.SaveChapterRecord(chapterRecordManager, chapterRecord)
	if err != nil {
		log.Println(err)
		http.Error(w, errServerError{err: "Problem saving record"}.Error(), http.StatusInternalServerError)
		return
	}
}

func getCurrentLessonHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		http.Error(w, errMissingPathVar{MissingVar: "Missing UID"}.Error(), http.StatusBadRequest)
		return
	}

	lesson, err := records.CurrentLesson(lessonRecordManager, lessonManager, uid)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting current lesson", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(lesson)
	if err != nil {
		log.Println(err)
		http.Error(w, errServerError{err: "Problem sending lesson record"}.Error(), http.StatusInternalServerError)
		return
	}
}

func getLessonRecordsForUserHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		http.Error(w, "Missing uid", http.StatusBadRequest)
		return
	}

	var recs []records.LessonRecord
	err := records.QueryLessonRecords(lessonRecordManager, &recs, records.QueryLessonsRecordParams{UID: uid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(recs)
	err = json.NewEncoder(w).Encode(recs)
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

	stats, err := records.QueryAvgTutorialStats(lessonRecordManager, uid)
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

//
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

	stats := records.LessonStats{}
	err := records.QueryLessonStats(lessonRecordManager, stats, records.QueryLessonStatsParams{LessonID: lessonID, UID: uid})
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

func getChapterProgressPercentage(w http.ResponseWriter, r *http.Request) {
	uid, ok := mux.Vars(r)["uid"]
	if !ok {
		http.Error(w, errMissingPathVar{MissingVar: "UID"}.Error(), http.StatusBadRequest)
		return
	}

	chapterProgress, err := records.QueryChapterProgress(chapterRecordManager, records.QueryChaptersRecordParams{UID: uid})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(chapterProgress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
