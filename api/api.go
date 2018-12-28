package api

import (
	"lavazares/content"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
)

var lessonManager *content.DefaultLessonManager
var chapterManager *content.DefaultChapterManager

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

	a.BaseRouter = mux.NewRouter()

	lessonManager = content.NewDefaultLessonManager(db)
	chapterManager = content.NewDefaultChapterManager(db)

	lessonRouter := a.BaseRouter.PathPrefix("/lesson").Subrouter()
	lessonRouter.HandleFunc("/", LessonsHandler)
	lessonRouter.HandleFunc("/{id}", LessonHandler)

	chapterRouter := a.BaseRouter.PathPrefix("/chapter").Subrouter()
	chapterRouter.HandleFunc("/{id}", ChapterHandler)
	chapterRouter.HandleFunc("/", ChaptersHandler)

}

func initDB(source string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	return db, nil
}
