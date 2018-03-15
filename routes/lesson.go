package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lavazares/models"
)

func HandleLessonCreate(w http.ResponseWriter, r *http.Request) {
	lessonRequest, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("Could not convert request to bytes: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l := models.Lesson{}

	err = json.Unmarshal(lessonRequest, &l)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = models.NewLesson(lessonRequest)
	if err != nil {
		log.Printf("Error creating error object: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func HandleChapterCreate(w http.ResponseWriter, r *http.Request) {
	chapterRequest, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("Could not convert request to bytes: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chapter, err := models.NewChapter(chapterRequest)
	fmt.Println(chapter)
	if err != nil {
		log.Printf("Error inserting chapter: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func HandleUnitCreate(w http.ResponseWriter, r *http.Request) {
	unitRequest, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("Error reading unit request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	unit, err := models.NewUnit(unitRequest)
	if err != nil {
		log.Printf("Error creating unit model: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(unit)
	w.WriteHeader(http.StatusOK)
	return
}

func HandleUserCompletedLesson(w http.ResponseWriter, r *http.Request) {
	lessonCompleteReq, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("Could not convert lesson request to bytes: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	completedLesson := models.LessonsComplete{}
	json.Unmarshal(lessonCompleteReq, &completedLesson)

	err = models.UserCompletedLesson(completedLesson)
	if err != nil {
		log.Printf("Could not add completed lesson: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func HandleUserCompletedChapter(w http.ResponseWriter, r *http.Request) {
	chapterCompleteReq, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("Error converting request to bytes: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	completedChapter := models.ChapterComplete{}
	json.Unmarshal(chapterCompleteReq, &completedChapter)

	err = models.UserCompletedChapter(completedChapter)
	if err != nil {
		log.Printf("Error inserting into db: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func HandleUserCompletedUnit(w http.ResponseWriter, r *http.Request) {
	unitCompleteReq, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("Error converting request to bytes: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	completedUnit := models.UnitComplete{}
	json.Unmarshal(unitCompleteReq, &completedUnit)

	err = models.UserCompletedUnit(completedUnit)
	if err != nil {
		log.Printf("Error inserting into db: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
