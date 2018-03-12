package main

import (
	"fmt"
	"lavazares/models"
	"lavazares/routes"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	lesson := router.PathPrefix("/learn").Subrouter()

	auth.HandleFunc("/login", routes.HandleLogin).Methods("POST")
	auth.HandleFunc("/signup", routes.HandleSignup).Methods("POST")

	lesson.HandleFunc("/create", routes.HandleLessonCreate).Methods("POST")
	lesson.HandleFunc("/completed", routes.HandleUserCompletedLesson).Methods("POST")
	// home.Use(routes.AuthMiddleware)
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Println("listening on port 8081")
	log.Println(http.ListenAndServe(":8081", cors.Default().Handler(loggingRouter)))
}
