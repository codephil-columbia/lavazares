package main

import "github.com/lavazares/app"

func main() {
	a := app.App{}
	a.Init("codephil", "lavazaresDB", "password", "5432", "localhost")
	a.Run(":8081")
}
