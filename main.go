package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"lavazares/routes"

	"lavazares/models"

	"github.com/rs/cors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	// connStr = "user=codephil dbname=lavazaresdb password=codephil! port=5432 host=lavazares-db1.cnodp99ehkll.us-west-2.rds.amazonaws.com sslmode=disable"
	connStr = "port=1000 host=localhost sslmode=disable user=postgres dbname=postgres"
)

func main() {

	if err := models.InitDB(connStr); err != nil {
		fmt.Println(err)
	}

	router := mux.NewRouter()

	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", routes.HandleLogin).Methods("POST")
	auth.HandleFunc("/signup", routes.HandleSignup).Methods("POST")
	auth.HandleFunc("/usernameValid", routes.CheckUsernameAvailable).Methods("POST")
	auth.HandleFunc("/newPassword", routes.HandleNewPassword).Methods("POST")

	lesson := router.PathPrefix("/lesson").Subrouter()
	lesson.HandleFunc("/create", routes.HandleLessonCreate).Methods("POST")
	lesson.HandleFunc("/finished", routes.HandleUserCompletedLesson).Methods("POST")
	lesson.HandleFunc("/get", routes.GetLessonByID).Methods("POST")
	lesson.HandleFunc("/getNext", routes.GetNextLessonForStudent).Methods("POST")
	lesson.HandleFunc("/getCurrent", routes.GetCurrentLessonForStudent).Methods("POST")
	lesson.HandleFunc("/getCompletedLessons", routes.GetCompletedLessonsForUser).Methods("POST")
	lesson.HandleFunc("/complete", routes.HandleLessonComplete).Methods("POST")

	chapter := router.PathPrefix("/chapter").Subrouter()
	chapter.HandleFunc("/create", routes.HandleChapterCreate).Methods("POST")
	chapter.HandleFunc("/completed", routes.HandleUserCompletedChapter).Methods("POST")
	chapter.HandleFunc("/getAllNames", routes.GetChapterNames).Methods("GET")
	chapter.HandleFunc("/getAllInfo", routes.GetAllLessonsForAllChapters).Methods("GET")
	chapter.HandleFunc("/getChapterProgress", routes.GetChapterProgress).Methods("POST")

	router.HandleFunc("/bulk", routes.HandleBulkGet).Methods("POST")
	router.HandleFunc("/update", routes.UpdateModel).Methods("POST")
	router.HandleFunc("/hollisticStats", routes.GetHollisticStats).Methods("POST")

	// home.Use(routes.AuthMiddleware)
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Println("listening on port 5000")
	log.Println(http.ListenAndServe(":5000", cors.Default().Handler(loggingRouter)))
}
