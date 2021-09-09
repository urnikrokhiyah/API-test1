package main

import (
	"cleancode/config"
	"cleancode/middlewares"
	"cleancode/routes"
)

func main() {
	config.InitDb()
	e := routes.New()
	middlewares.LogMiddleware(e)
	e.Logger.Fatal(e.Start(":8000"))
}
