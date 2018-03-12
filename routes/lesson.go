package routes

import (
	"encoding/json"
	"fmt"
	"lavazares/models"
	"log"
	"net/http"
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

func HandleUserCompletedLesson(w http.ResponseWriter, r *http.Request) {
	lessonCompleteReq, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("Could not convert lesson request to bytes: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	completedLesson := models.LessonComplete{}
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
