package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lavazares/models"
	"github.com/lavazares/routes"
	_ "github.com/lib/pq"
)

const (
	connStr = "user=codephil dbname=lavazaresDB password=password port=5432 host=localhost sslmode=disable"
)

func main() {
	err := models.InitDB(connStr)
	if err != nil {
		log.Panicf("database could not be opened: %s", err)
	}

	err = models.InitRedisCache()

	router := mux.NewRouter()
	auth := router.PathPrefix("/auth").Subrouter()
	home := router.PathPrefix("/home").Subrouter()

	auth.HandleFunc("/login", routes.HandleLogin).Methods("POST")
	auth.HandleFunc("/signup", routes.HandleSignup).Methods("POST")

	home.HandleFunc("/Test", routes.Test).Methods("GET")
	home.Use(routes.AuthMiddleware)

	log.Println("listening on port 8081")
	http.ListenAndServe(":8081", router)
	log.Println("finished server")
}
