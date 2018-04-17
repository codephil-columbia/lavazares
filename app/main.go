package app

import (
	"fmt"
	"lavazares/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type App struct {
	Router *mux.Router
}

func NewApp() (app *App) {
	router := mux.NewRouter()
	a := App{Router: router}
	return &a
}

func (a *App) Run(port string) {

	auth := a.Router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/login", routes.HandleLogin).Methods("POST")
	auth.HandleFunc("/signup", routes.HandleSignup).Methods("POST")

	lesson := a.Router.PathPrefix("/lesson").Subrouter()
	lesson.HandleFunc("/create", routes.HandleLessonCreate).Methods("POST")

	chapter := a.Router.PathPrefix("/chapter").Subrouter()
	chapter.HandleFunc("/create", routes.HandleChapterCreate).Methods("POST")

	a.Router.HandleFunc("/bulk", routes.HandleBulkGet).Methods("POST")

	// home.Use(routes.AuthMiddleware)
	loggingRouter := handlers.LoggingHandler(os.Stdout, a.Router)

	log.Println(http.ListenAndServe(fmt.Sprintf(":%s", port), cors.Default().Handler(loggingRouter)))
}
