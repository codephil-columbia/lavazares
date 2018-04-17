package main

import (
	"flag"
	"fmt"

	"lavazares/app"
	"lavazares/models"
)

const (
	connStr = "user=codephil dbname=lavazaresDB password=password port=5432 host=localhost sslmode=disable"
)

func main() {

	user := flag.String("user", "codephil", "database user")
	dbname := flag.String("dbname", "lavazaresDB", "dbname")
	password := flag.String("password", "password", "dbpassword")
	port := flag.String("port", "5432", "dbport")
	host := flag.String("host", "localhost", "dbhost")
	ssl := flag.String("ssl", "disable", "ssl mode")

	connStr := fmt.Sprintf("user=%s dbname=%s, password=%s, port-%s, host=%s, sslmode=%s", *user, *dbname, *password, *port, *host, *ssl)

	if err := models.InitDB(connStr); err != nil {
		fmt.Println(err)
	}

	if err := models.InitRedisCache(); err != nil {
		fmt.Println(err)
	}

	app := app.NewApp()

	app.Run("8081")
}
