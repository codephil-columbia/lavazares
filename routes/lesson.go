package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"lavazares/models"
)

func GetNextLessonForStudent(w http.ResponseWriter, r *http.Request) {
	req, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := make(map[string]string)
	err = json.Unmarshal(req, &body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	nextLessonInfo, err := models.NextLessonForStudent(body["uid"])
	fmt.Println(nextLessonInfo)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(nextLessonInfo)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}

func logErrAndReturn(err error, w http.ResponseWriter) {
	log.Printf("%v", err)
	w.WriteHeader(http.StatusBadRequest)
	return
}

func GetCurrentLessonForStudent(w http.ResponseWriter, r *http.Request) {
	req, err := requestToBytes(r.Body)
	if err != nil {
		logErrAndReturn(err, w)
	}

	body := make(map[string]string)
	err = json.Unmarshal(req, &body)
	if err != nil {
		logErrAndReturn(err, w)
	}

	currentLesson, err := models.GetCurrent(body["uid"])
	if err != nil {
		logErrAndReturn(err, w)
	}

	err = json.NewEncoder(w).Encode(currentLesson)
	if err != nil {
		logErrAndReturn(err, w)
	}
	return

}

func GetChapterProgress(w http.ResponseWriter, r *http.Request) {
	req, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body := make(map[string]string)
	err = json.Unmarshal(req, &body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	progress, err := models.GetProgressForCurrentUserLesson(body["uid"])
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(progress)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}

func GetHollisticStats(w http.ResponseWriter, r *http.Request) {
	req, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body := make(map[string]string)
	err = json.Unmarshal(req, &body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stats, err := models.GetOverallWPMAndAccuracy(body["uid"])
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(stats)
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}

func GetCompletedLessonsForUser(w http.ResponseWriter, r *http.Request) {
	req, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body := make(map[string]string)
	err = json.Unmarshal(req, &body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lessonsComplete, err := models.GetCompletedLessonsForUser(body["uid"])
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(lessonsComplete)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}

func GetAllLessonsForAllChapters(w http.ResponseWriter, r *http.Request) {
	allInfo, err := models.GetAllLessonsChapters()
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(allInfo)
	return
}

func GetChapterNames(w http.ResponseWriter, r *http.Request) {
	names, err := models.GetAllChapterNames()
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(names)
	return
}

func GetLessonByID(w http.ResponseWriter, r *http.Request) {
	req, _ := requestToBytes(r.Body)
	body := make(map[string]string)
	json.Unmarshal(req, &body)

	fmt.Println(body["lessonid"])

	l, err := models.GetLesson(body["lessonid"])
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(l)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func UpdateModel(w http.ResponseWriter, r *http.Request) {
	req, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := make(map[string]string)
	err = json.Unmarshal(req, &body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(body)

	err = models.UpdateModel(body["model"], body["field"], body["val"], body["identifier"], body["identifierVal"])
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return

}

func HandleLessonCreate(w http.ResponseWriter, r *http.Request) {
	// lessonRequest, err := requestToBytes(r.Body)
	// if err != nil {
	// 	log.Printf("Could not convert request to bytes: %v", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// l := models.Lesson{}

	// err = json.Unmarshal(lessonRequest, &l)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// _, err = models.NewLesson(lessonRequest)
	// if err != nil {
	// 	log.Printf("Error creating error object: %v", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// return
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

	completedLesson := models.LessonsComplete{}
	err = json.Unmarshal(lessonCompleteReq, &completedLesson)
	if err != nil {
		log.Printf("Could not convert lesson request to struct: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	completedChapter := models.ChaptersComplete{}
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
	// unitCompleteReq, err := requestToBytes(r.Body)
	// if err != nil {
	// 	log.Printf("Error converting request to bytes: %v", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// completedUnit := models.UnitComplete{}
	// json.Unmarshal(unitCompleteReq, &completedUnit)

	// err = models.UserCompletedUnit(completedUnit)
	// if err != nil {
	// 	log.Printf("Error inserting into db: %v", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// return
}

func HandleBulkGet(w http.ResponseWriter, r *http.Request) {
	uidReq, err := requestToBytes(r.Body)
	if err != nil {
		log.Printf("Error getting uid: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(uidReq, &m)
	if err != nil {
		log.Printf("Error unmarshalling uid: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bulkInfo, err := models.AllLessons()
	if err != nil {
		log.Printf("Error getting bulk info for user with id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(bulkInfo)
	if err != nil {
		log.Printf("Error writing json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	req, err := requestToBytes(r.Body)
	defer r.Body.Close()
	if err != nil {
		logErrAndReturn(err, w)
	}

	t := models.LessonsComplete{}
	err = json.Unmarshal(req, &t)
	if err != nil {
		logErrAndReturn(err, w)
	}

	err = models.UserDidFinishLesson(t)
	if err != nil {
		logErrAndReturn(err, w)
	}

	w.WriteHeader(200)
}
