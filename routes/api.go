package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"lavazares/auth"
	"lavazares/content"
	"lavazares/records"
	"log"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
)

var (
	lessonManager         *content.LessonManager
	chapterManager        *content.ChapterManager
	userManager           *auth.UserManager
	tutorialRecordManager *records.TutorialRecordManager

	app *API
)

const (
	secrets      = "./secrets.json"
	localConnStr = "port=5432 host=localhost sslmode=disable user=postgres dbname=postgres"
)

type (
	errMissingPathVar struct {
		MissingVar string
	}
)

func (err errMissingPathVar) Error() string {
	return fmt.Sprintf("Missing path var %v", err.MissingVar)
}

var (
	errMissingUID = errors.New("Missing UID")
	errBadJSON    = errors.New("Could not read json")
)

type productionCredentials struct {
	ProductionDB string `json:"productionDB"`
}

// Run initializes the App
func Run(isLocal bool) *API {
	if isLocal {
		log.Println("Running local db!")
		return initAPI(localConnStr)
	}

	log.Println("Running with production db! :)")
	productionCredentials, err := getProductionCredentials(secrets)
	if err != nil {
		return nil
	}

	return initAPI(productionCredentials.ProductionDB)
}

func getProductionCredentials(path string) (*productionCredentials, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	productionCredentials := productionCredentials{}
	err = json.Unmarshal(data, &productionCredentials)
	if err != nil {
		return nil, err
	}

	return &productionCredentials, nil
}

// API holds information about the currently running app
type API struct {
	BaseRouter *mux.Router
}

func initAPI(connStr string) *API {
	app := API{}

	db, err := initDB(connStr)
	if err != nil {
		return nil
	}

	lessonManager = content.NewLessonManager(db)
	chapterManager = content.NewChapterManager(db)
	userManager = auth.NewUserManager(auth.NewUserStore(db))
	tutorialRecordManager = records.NewTutorialRecordManager(db)

	app.BaseRouter = mux.NewRouter()

	lessonRouter := app.BaseRouter.PathPrefix("/lesson").Subrouter()
	lessonRouter.HandleFunc("/", LessonsHandler)
	lessonRouter.HandleFunc("/{id}", LessonHandler)
	lessonRouter.HandleFunc("/current/{uid}", getNextNonCompletedLesson)

	chapterRouter := app.BaseRouter.PathPrefix("/chapter").Subrouter()
	chapterRouter.HandleFunc("/{id}", ChapterHandler)
	chapterRouter.HandleFunc("/", ChaptersHandler)
	chapterRouter.HandleFunc("/current/{uid}", getNextNonCompletedChapter)

	userRouter := app.BaseRouter.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", newUserHandler).Methods("POST")
	userRouter.HandleFunc("/edit/password", editPasswordHandler).Methods("POST")
	userRouter.HandleFunc("/authenticate", authenticateHandler).Methods("POST")

	recordRouter := app.BaseRouter.PathPrefix("/records").Subrouter()

	tutorialRouter := recordRouter.PathPrefix("/tutorial").Subrouter()
	tutorialRouter.HandleFunc("/lessons/{uid}", getLessonRecordsForUserHandler)
	tutorialRouter.HandleFunc("/save/{type}", saveTutorialRecord).Methods("POST")

	statsRouter := app.BaseRouter.PathPrefix("/stats").Subrouter()
	statsRouter.HandleFunc("/tutorial/lesson/{uid}", getTutorialHollisticLessonStatsHandler)
	statsRouter.HandleFunc("/tutorial/lesson/{lessonid}/{uid}", getTutorialLessonStatsHandler)

	return &app
}

func initDB(source string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	return db, nil
}
