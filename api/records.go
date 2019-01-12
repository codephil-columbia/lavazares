package api

import (
	"lavazares/records"
	"lavazares/utils"
	"net/http"
)

func addLessonRecordHandler(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ReadBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := records.NewLessonRecord(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = tutorialRecordManager.Save(record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// func lessonRecordHandler(w http.ResponseWriter, r *http.Request) {
// 	uid, ok := mux.Vars(r)["uid"]
// 	if !ok {
// 		http.Error(w, "Missing uid", http.StatusBadRequest)
// 		return
// 	}

// 	// recordID
// }
