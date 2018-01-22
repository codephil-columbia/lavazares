package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/lavazares/app/routes"
)

//App holds application data
type App struct {
	Router         *mux.Router
	SessionManager *redis.Client
	DB             *sql.DB
}

//Run inits the app on the specified port
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

//Init does all the initialization required for the app to run
func (a *App) Init(user, dbname, password, port, host string) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s port=%s host=%s sslmode=disable", user, dbname, password, port, host)

	var err error
	if a.DB, err = sql.Open("postgres", connStr); err != nil {
		log.Fatal(err)
	}

	a.SessionManager = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err = a.SessionManager.Ping().Result()
	if err != nil {
		log.Fatalln("error connecting to redis: %s", err)
	}

	a.Router = mux.NewRouter()
	a.InitRoutes()
}

//InitRoutes initializes the apps routes
func (a *App) InitRoutes() {
	router := a.Router
	auth := router.PathPrefix("/auth").Subrouter()
	home := router.PathPrefix("/home").Subrouter()

	auth.HandleFunc("/login", routes.HandleLogin).Methods("POST")
	auth.HandleFunc("/signup", routes.HandleSignup).Methods("POST")

	home.HandleFunc("/Test", routes.Test).Methods("GET")
	home.Use(routes.RequiresAuth)
}
