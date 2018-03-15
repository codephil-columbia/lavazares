package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lavazares/routes"

	"github.com/lavazares/models"

	"github.com/rs/cors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	connStr = "user=codephil dbname=lavazaresDB password=password port=5432 host=localhost sslmode=disable"
)

func main() {

	if err := models.InitDB(connStr); err != nil {
		fmt.Println(err)
	}

	if err := models.InitRedisCache(); err != nil {
		fmt.Println(err)
	}

	router := mux.NewRouter()

	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", routes.HandleLogin).Methods("POST")
	auth.HandleFunc("/signup", routes.HandleSignup).Methods("POST")

	lesson := router.PathPrefix("/lesson").Subrouter()
	lesson.HandleFunc("/create", routes.HandleLessonCreate).Methods("POST")
	lesson.HandleFunc("/completed", routes.HandleUserCompletedLesson).Methods("POST")

	chapter := router.PathPrefix("/chapter").Subrouter()
	chapter.HandleFunc("/create", routes.HandleChapterCreate).Methods("POST")
	chapter.HandleFunc("/completed", routes.HandleUserCompletedChapter).Methods("POST")

	unit := router.PathPrefix("/unit").Subrouter()
	unit.HandleFunc("/create", routes.HandleUnitCreate).Methods("POST")
	unit.HandleFunc("/completed", routes.HandleUserCompletedUnit).Methods("POST")

	// home.Use(routes.AuthMiddleware)
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Println("listening on port 8081")
	log.Println(http.ListenAndServe(":8081", cors.Default().Handler(loggingRouter)))
}
