package main

import (
	"log"

	"g0/database"
	"g0/routes"

	"github.com/labstack/echo/v5"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	e := echo.New()

	routes.RegisterRoute(e)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("Failed to start server", "error", err)
	}
}
