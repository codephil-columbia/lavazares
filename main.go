package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"lavazares/models"
	"lavazares/routes"

	"github.com/rs/cors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	// connStr = "user=codephil dbname=lavazaresdb password=codephil! port=5432 host=lavazares-db1.cnodp99ehkll.us-west-2.rds.amazonaws.com sslmode=disable"
	connStr = "port=5432 host=localhost sslmode=disable user=postgres dbname=postgres"
)

func main() {

	if err := models.InitDB(connStr); err != nil {
		fmt.Println(err)
		return
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

type Config struct {
	Database struct {
		User     string `json:"user"`
		DBName   string `json:"dbName"`
		Password string `json:"password"`
		Port     string `json:"port"`
		Host     string `json:"host"`
		SSLMode  string `json:"sslmode"`
	} `json:"dbCredentials"`
}

func loadFromFile(path string) (*Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	return &config, err
}

func (c *Config) getDBConnString() string {
	return fmt.Sprintf(
		"user=%s dbname=%s password=%s port=%s host=%s sslmode=%s",
		c.Database.User,
		c.Database.DBName,
		c.Database.Password,
		c.Database.Port,
		c.Database.Host,
		c.Database.SSLMode,
	)
}

func Init(isLocalEnv bool, configPath string) error {
	var connStr string
	if isLocalEnv {
		fmt.Printf("Running local env with psql connStr %s\n", connStr)
	} else {
		config, err := loadFromFile(configPath)
		if err != nil {
			return err
		}
		connStr = config.getDBConnString()
		fmt.Printf("Running prod env with psql connStr %s\n", connStr)
	}

	err := models.InitDB(connStr)
	return err
}
