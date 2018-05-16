package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"lavazares/routes"

	"lavazares/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	connStr = "user=codephil dbname=lavazaresdb password=codephil! port=5432 host=lavazares-db1.cnodp99ehkll.us-west-2.rds.amazonaws.com sslmode=disable"
)

func main() {

	if err := models.InitDB(connStr); err != nil {
		fmt.Println(err)
	}

	router := mux.NewRouter()

	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", routes.HandleLogin).Methods("POST")
	auth.HandleFunc("/signup", routes.HandleSignup).Methods("POST")

	lesson := router.PathPrefix("/lesson").Subrouter()
	lesson.HandleFunc("/create", routes.HandleLessonCreate).Methods("POST")
	lesson.HandleFunc("/completed", routes.HandleUserCompletedLesson).Methods("POST")
	lesson.HandleFunc("/get", routes.GetLessonByID).Methods("POST")

	chapter := router.PathPrefix("/chapter").Subrouter()
	chapter.HandleFunc("/create", routes.HandleChapterCreate).Methods("POST")
	chapter.HandleFunc("/completed", routes.HandleUserCompletedChapter).Methods("POST")

	unit := router.PathPrefix("/unit").Subrouter()
	unit.HandleFunc("/create", routes.HandleUnitCreate).Methods("POST")
	unit.HandleFunc("/completed", routes.HandleUserCompletedUnit).Methods("POST")

	router.HandleFunc("/bulk", routes.HandleBulkGet).Methods("POST")
	router.HandleFunc("/update", routes.UpdateModel).Methods("POST")

	// home.Use(routes.AuthMiddleware)
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Println("listening on port 5000")
	log.Println(http.ListenAndServe(":5000", loggingRouter))
}
