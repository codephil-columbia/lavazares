package routes

import (
	"errors"
	"fmt"
	"lavazares/auth"
	"lavazares/content"
	"lavazares/records"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
)

var (
	lessonManager         *content.LessonManager
	chapterManager        *content.ChapterManager
	userManager           *auth.UserManager
	tutorialRecordManager *records.TutorialRecordManager

	a *API
)

const (
	localConnStr = "port=5432 host=localhost sslmode=disable user=postgres dbname=postgres"
	prodConnStr  = "user=codephil dbname=lavazaresdb password=codephil! port=5432 host=lavazares-db1.cnodp99ehkll.us-west-2.rds.amazonaws.com sslmode=disable"
)

var (
	errMissingUID = errors.New("Missing UID")
)

// Run initializes the App
func Run(isLocal bool) *API {
	fmt.Println("is local", isLocal)
	if isLocal {
		return initAPI(localConnStr)
	}
	return initAPI(prodConnStr)
}

// API holds information about the currently running app
type API struct {
	BaseRouter *mux.Router
}

func initAPI(connStr string) *API {
	a := API{}

	db, err := initDB(connStr)
	if err != nil {
		return nil
	}

	lessonManager = content.NewLessonManager(db)
	chapterManager = content.NewChapterManager(db)
	userManager = auth.NewUserManager(auth.NewUserStore(db))
	tutorialRecordManager = records.NewTutorialRecordManager(db)

	baseRouter := mux.NewRouter()
	a.BaseRouter = baseRouter

	lessonRouter := a.BaseRouter.PathPrefix("/lesson").Subrouter()
	lessonRouter.HandleFunc("/", LessonsHandler)
	lessonRouter.HandleFunc("/{id}", LessonHandler)
	lessonRouter.HandleFunc("/current/{uid}", getNextNonCompletedLesson)

	chapterRouter := a.BaseRouter.PathPrefix("/chapter").Subrouter()
	chapterRouter.HandleFunc("/{id}", ChapterHandler)
	chapterRouter.HandleFunc("/", ChaptersHandler)
	chapterRouter.HandleFunc("/current/{uid}", getNextNonCompletedChapter)

	userRouter := a.BaseRouter.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", newUserHandler).Methods("POST")
	userRouter.HandleFunc("/edit/password", editPasswordHandler).Methods("POST")
	userRouter.HandleFunc("/authenticate", authenticateHandler).Methods("POST")

	recordRouter := a.BaseRouter.PathPrefix("/records").Subrouter()

	tutorialRouter := recordRouter.PathPrefix("/tutorial").Subrouter()
	tutorialRouter.HandleFunc("/lesson", addLessonRecordHandler).Methods("POST")
	tutorialRouter.HandleFunc("/lessons/{uid}", getLessonRecordsForUserHandler)

	statsRouter := a.BaseRouter.PathPrefix("/stats").Subrouter()
	statsRouter.HandleFunc("/tutorial/lesson/{uid}", getTutorialHollisticLessonStatsHandler)
	statsRouter.HandleFunc("/tutorial/lesson/{lessonid}/{uid}", getTutorialLessonStatsHandler)

	return &a
}

func initDB(source string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	return db, nil
}
