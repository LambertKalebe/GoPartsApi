package main

import (
	"database/sql"
	"fmt"
	"g0/internal/config"
	"g0/internal/database"
	"g0/internal/routes"
	"log"

	"github.com/labstack/echo/v5"
)

func main() {
	fmt.Println("DatabaseUrl:", config.GetDatabasePath())
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database: ", database.DB.Stats())
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
