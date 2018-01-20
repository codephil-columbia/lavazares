package main

import (
	"log"
	"net/http"

	"github.com/GoPhil/models"
	"github.com/GoPhil/routes"
	"github.com/gorilla/mux"
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
	router.HandleFunc("/login", routes.HandleLogin).Methods("POST")
	router.HandleFunc("/signup", routes.HandleSignup).Methods("POST")
	router.HandleFunc("/Test", routes.Test).Methods("GET")

	log.Println("listening on port 8081")
	http.ListenAndServe(":8081", router)
	log.Println("finished server")
}
