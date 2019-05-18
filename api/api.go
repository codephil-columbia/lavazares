package api

import (
	"lavazares/auth"
	"lavazares/content"
	"lavazares/records"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
)

var (
	lessonManager         *content.DefaultLessonManager
	chapterManager        *content.DefaultChapterManager
	userManager           *auth.DefaultUserManager
	tutorialRecordManager *records.TutorialRecordManager
	gameManager           *content.DefaultGameManager

	a *API
)

const (
	localPostgres = "port=5432 host=localhost sslmode=disable user=postgres dbname=postgres"
)

func Run() *API {
	return initAPI()
}

type API struct {
	BaseRouter *mux.Router
}

func initAPI() *API {
	a := API{}

	db, err := initDB(localPostgres)
	if err != nil {
		return nil
	}

	lessonManager = content.NewDefaultLessonManager(db)
	chapterManager = content.NewDefaultChapterManager(db)
	userManager = auth.NewDefaultUserManager(auth.NewUserStore(db))
	tutorialRecordManager = records.NewTutorialRecordManager(db)
	gameManager = content.NewDefaultGameContentManager(db)

	baseRouter := mux.NewRouter()
	a.BaseRouter = baseRouter

	lessonRouter := a.BaseRouter.PathPrefix("/lesson").Subrouter()
	lessonRouter.HandleFunc("/", LessonsHandler)
	lessonRouter.HandleFunc("/{id}", LessonHandler)

	chapterRouter := a.BaseRouter.PathPrefix("/chapter").Subrouter()
	chapterRouter.HandleFunc("/{id}", ChapterHandler)
	chapterRouter.HandleFunc("/", ChaptersHandler)

	userRouter := a.BaseRouter.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", newUserHandler).Methods("POST")
	userRouter.HandleFunc("/edit/password", editPasswordHandler).Methods("POST")
	userRouter.HandleFunc("/authenticate", authenticateHandler).Methods("POST")

	recordRouter := a.BaseRouter.PathPrefix("/records").Subrouter()
	recordRouter.HandleFunc("/tutorial", addLessonRecordHandler).Methods("POST")

	gameRouter := a.BaseRouter.PathPrefix("/game").Subrouter()
	gameRouter.HandleFunc("/boatrace", BoatgameHandler).Methods("GET")

	return &a
}

func initDB(source string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initRecordRouter() *mux.Router {
	recordRouter := a.BaseRouter.PathPrefix("/records").Subrouter()
	recordRouter.HandleFunc("/tutorial", addLessonRecordHandler).Methods("POST")
	return recordRouter
}
