package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"lavazares/routes"

	"github.com/rs/cors"

	"github.com/gorilla/handlers"
)

func main() {
	env := flag.Bool("local", true, "Specifies local vs prod env")
	flag.Parse()

	api := routes.Run(*env)
	loggingRouter := handlers.LoggingHandler(os.Stdout, api.BaseRouter)

	log.Println("listening on port 5000")
	log.Println(http.ListenAndServe(":5000", cors.Default().Handler(loggingRouter)))
}
