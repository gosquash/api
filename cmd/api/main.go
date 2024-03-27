package main

import (
	"fmt"
	"gosquash/api/internal/db"
	"gosquash/api/internal/routes"
	"gosquash/api/internal/validator"
	"gosquash/api/pkg/structs"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError

	he, ok := err.(*echo.HTTPError)
	if ok {
		code = he.Code
	}

	c.JSON(code, echo.Map{
		"error": echo.Map{
			"message": he.Message.(string),
			"status":  code,
		},
	})
}

func main() {

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Init()

	// sync database
	err = db.DB.AutoMigrate(
		&structs.User{},
		&structs.Game{},
		&structs.Player{},
	)

	e := echo.New()

	// Validator
	e.Validator = validator.NewValidator()

	// Custom error handler.
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Define routes
	routes.Init(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
