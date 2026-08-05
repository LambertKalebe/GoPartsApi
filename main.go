package main

import (
	"database/sql"
	"g0/internal/database"
	"g0/internal/routes"
	"log"

	"github.com/labstack/echo/v5"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {

		}
	}(database.DB)

	e := echo.New()

	routes.Route(e)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("Failed to start server", "error", err)
	}
}
