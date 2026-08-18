package main

import (
	"database/sql"
	"fmt"
	"g0/internal/common"
	"g0/internal/config"
	"g0/internal/database"
	"g0/internal/http"
	"g0/internal/middleware"
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
	e.HTTPErrorHandler = httpcustom.HTTPErrorHandler
	common.NewLogger()                  // new
	e.Use(middleware.LoggingMiddleware) // new

	routes.Route(e)
	common.Logger.LogInfo().Msg(e.Start(":8080").Error())

}
