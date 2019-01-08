package api

import (
	"lavazares/auth"
	"lavazares/content"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
)

var lessonManager *content.DefaultLessonManager
var chapterManager *content.DefaultChapterManager
var userManager *auth.DefaultUserManager

const (
	localPostgres = "port=5432 host=localhost sslmode=disable user=postgres dbname=postgres"
)

func Run() API {
	a := API{}
	a.initAPI()
	return a
}

type API struct {
	BaseRouter *mux.Router
}

func (a *API) initAPI() {
	db, err := initDB(localPostgres)
	if err != nil {
		return
	}

	lessonManager = content.NewDefaultLessonManager(db)
	chapterManager = content.NewDefaultChapterManager(db)
	userManager = auth.NewDefaultUserManager(auth.NewUserStore(db))

	a.BaseRouter = mux.NewRouter()

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

}

func initDB(source string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	return db, nil
}
